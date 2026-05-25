package components

import (
	"context"
	"fmt"
	"sync"
	"time"

	"grid-simulator/config"
	"grid-simulator/models"
)

type BaseConsumer struct {
	ID       string
	Priority int
	BaseReq  float64
	logger   *CSVLogger
	wg       *sync.WaitGroup
	hourTick int
}

func (c *BaseConsumer) GetPriority() int {
	return c.Priority
}

func (c *BaseConsumer) Start(ctx context.Context, demandChan chan<- models.DemandReport, calcDemand func(int, float64) float64) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(config.GridStep)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.hourTick = (c.hourTick + 1) % 24
				currentDemand := calcDemand(c.hourTick, c.BaseReq)

				replyChan := make(chan models.SupplyStatus, 1)

				report := models.DemandReport{
					ID:       c.ID,
					PDemand:  currentDemand,
					Priority: c.Priority,
					ReplyTo:  replyChan,
				}

				demandChan <- report

				select {
				case <-ctx.Done():
					return
				case status := <-replyChan:
					if status.AllocatedMW < currentDemand {
						c.logger.LogEvent(fmt.Sprintf("Konsument %s: Partial/LoadShed (Żądano: %.2f MW, Otrzymano: %.2f MW. Powód: %s)",
							c.ID, currentDemand, status.AllocatedMW, status.Reason))
					}
				}
			}
		}
	}()
}

func NewResidentialConsumer(id string, logger *CSVLogger, wg *sync.WaitGroup) *BaseConsumer {
	c := &BaseConsumer{ID: id, Priority: config.PriorityResidential, BaseReq: 5.0, logger: logger, wg: wg}
	return c
}
func (c *BaseConsumer) StartResidential(ctx context.Context, demandChan chan<- models.DemandReport) {
	c.Start(ctx, demandChan, func(hour int, base float64) float64 {

		if (hour >= 7 && hour <= 9) || (hour >= 18 && hour <= 22) {
			return base * 2.5
		}
		return base
	})
}

func NewIndustrialConsumer(id string, logger *CSVLogger, wg *sync.WaitGroup) *BaseConsumer {
	c := &BaseConsumer{ID: id, Priority: config.PriorityIndustrial, BaseReq: 20.0, logger: logger, wg: wg}
	return c
}
func (c *BaseConsumer) StartIndustrial(ctx context.Context, demandChan chan<- models.DemandReport) {
	c.Start(ctx, demandChan, func(hour int, base float64) float64 {
		if hour >= 6 && hour <= 18 {
			if hour == 8 || hour == 14 {
				return base * 1.8
			}
			return base
		}
		return base * 0.2 
	})
}

func NewCriticalConsumer(id string, logger *CSVLogger, wg *sync.WaitGroup) *BaseConsumer {
	c := &BaseConsumer{ID: id, Priority: config.PriorityCritical, BaseReq: 10.0, logger: logger, wg: wg}
	return c
}
func (c *BaseConsumer) StartCritical(ctx context.Context, demandChan chan<- models.DemandReport) {
	c.Start(ctx, demandChan, func(hour int, base float64) float64 {
		return base
	})
}