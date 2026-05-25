package components

import (
	"context"
	"sync"
	"time"

	"grid-simulator/config"
	"grid-simulator/models"
)

type WindFarm struct {
	currentPower float64
	mu           sync.RWMutex
	wg           *sync.WaitGroup
}

func NewWindFarm(wg *sync.WaitGroup) *WindFarm {
	return &WindFarm{wg: wg}
}

func (w *WindFarm) Start(ctx context.Context, weatherChan <-chan models.WeatherData) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-weatherChan:
				if !ok {
					return
				}
				power := data.WindSpeed * 2.0
				w.mu.Lock()
				w.currentPower = power
				w.mu.Unlock()
			}
		}
	}()
}

func (w *WindFarm) GetPower() float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.currentPower
}

type CoalPlant struct {
	maxPower  float64
	isOn      bool
	isWarming bool
	mu        sync.RWMutex
	wg        *sync.WaitGroup
}

func NewCoalPlant(maxPower float64, wg *sync.WaitGroup) *CoalPlant {
	return &CoalPlant{maxPower: maxPower, wg: wg}
}

func (c *CoalPlant) Start(ctx context.Context) {
	c.mu.Lock()
	if c.isOn || c.isWarming {
		c.mu.Unlock()
		return
	}
	c.isWarming = true
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		
		timer := time.NewTimer(config.GridStep * 2)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			c.mu.Lock()
			c.isWarming = false
			c.isOn = true
			c.mu.Unlock()
		}
	}()
}

func (c *CoalPlant) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isOn = false
	c.isWarming = false
}

func (c *CoalPlant) GetPower() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.isOn {
		return c.maxPower
	}
	return 0.0
}