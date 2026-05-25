package components

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"grid-simulator/config"
	"grid-simulator/models"
)

type GridHub struct {
	ozePlant   *WindFarm
	coalPlant  *CoalPlant
	battery    *BatteryStorage
	logger     *CSVLogger
	demandChan <-chan models.DemandReport
	fcstChan   <-chan models.ForecastReport
	wg         *sync.WaitGroup
}

func NewGridHub(oze *WindFarm, coal *CoalPlant, bat *BatteryStorage, logger *CSVLogger, dc <-chan models.DemandReport, fc <-chan models.ForecastReport, wg *sync.WaitGroup) *GridHub {
	return &GridHub{
		ozePlant:   oze,
		coalPlant:  coal,
		battery:    bat,
		logger:     logger,
		demandChan: dc,
		fcstChan:   fc,
		wg:         wg,
	}
}

func (g *GridHub) Start(ctx context.Context) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		
		ticker := time.NewTicker(config.GridStep)
		defer ticker.Stop()

		currentDemands := make(map[string]models.DemandReport)
		stepCount := 0

		for {
			select {
			case <-ctx.Done():
				return

			case req := <-g.demandChan:
				currentDemands[req.ID] = req

			case forecast := <-g.fcstChan:
				g.logger.LogEvent(forecast.Message)
				if forecast.Trend < -2.0 {
					g.logger.LogEvent("Ostrzeżenie: Prognozowany spadek OZE. Uruchamianie procedury rozgrzewania elektrowni węglowej.")
					g.coalPlant.Start(ctx)
				}

			case <-ticker.C:
				stepCount++
				
				totalDemand := 0.0
				var demandsList []models.DemandReport
				
				for _, req := range currentDemands {
					totalDemand += req.PDemand
					demandsList = append(demandsList, req)
				}

				ozePower := g.ozePlant.GetPower()
				coalPower := g.coalPlant.GetPower()
				totalProduction := ozePower + coalPower
				
				balance := totalProduction - totalDemand
				
				g.manageGrid(balance, &demandsList)
				
				if stepCount%5 == 0 {
					state := "STABLE"
					if balance < 0 && g.battery.GetSoC() == 0 {
						state = "CRITICAL (Load Shedding)"
					}
					fmt.Printf("\n--- RAPORT SIECI [Krok %d] ---\n", stepCount)
					fmt.Printf("[Produkcja] OZE: %.2f MW | Konwencjonalna: %.2f MW | Baterie: %.0f%%\n", ozePower, coalPower, g.battery.GetSoC()*100)
					fmt.Printf("[Sieć] Popyt: %.2f MW | Bilans (przed baterią): %.2f MW | Stan: %s\n", totalDemand, balance, state)
					fmt.Println("------------------------------")
				}

				currentDemands = make(map[string]models.DemandReport)
			}
		}
	}()
}

func (g *GridHub) manageGrid(balance float64, demandsList *[]models.DemandReport) {
	if balance >= 0 {
		soc := g.battery.GetSoC()
		if soc < 1.0 {
			charged := g.battery.Charge(balance)
			g.logger.LogEvent(fmt.Sprintf("Nadwyżka w sieci. Ładowanie ESS: +%.2f MW", charged))
		} else {
			g.logger.LogEvent("Nadwyżka w sieci, ESS pełny (100%). Ograniczanie mocy OZE (Curtailment).")
		}

		for _, req := range *demandsList {
			req.ReplyTo <- models.SupplyStatus{AllocatedMW: req.PDemand, Reason: "OK"}
		}
		return
	}

	deficit := -balance
	soc := g.battery.GetSoC()

	if soc > 0 {
		discharged := g.battery.Discharge(deficit)
		deficit -= discharged
		g.logger.LogEvent(fmt.Sprintf("Niedobór w sieci. Rozładowanie ESS: -%.2f MW", discharged))
	}

	if deficit > 0.01 {
		g.logger.LogEvent(fmt.Sprintf("CRITICAL: Brak mocy w sieci i bateriach (Deficyt: %.2f MW). Uruchamianie Load Shedding.", deficit))

		sort.Slice(*demandsList, func(i, j int) bool {
			return (*demandsList)[i].Priority > (*demandsList)[j].Priority
		})

		for i, req := range *demandsList {
			if deficit > 0 {
				(*demandsList)[i].ReplyTo <- models.SupplyStatus{AllocatedMW: 0, Reason: "LoadShed"}
				deficit -= req.PDemand
			} else {
				(*demandsList)[i].ReplyTo <- models.SupplyStatus{AllocatedMW: req.PDemand, Reason: "OK"}
			}
		}
	} else {
		for _, req := range *demandsList {
			req.ReplyTo <- models.SupplyStatus{AllocatedMW: req.PDemand, Reason: "OK"}
		}
	}
}