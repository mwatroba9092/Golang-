package components

import (
	"context"
	"fmt"
	"sync"

	"grid-simulator/config"
	"grid-simulator/models"
)

type TrendPredictor struct {
	buffer []float64
	wg     *sync.WaitGroup
}

func NewPredictor(wg *sync.WaitGroup) *TrendPredictor {
	return &TrendPredictor{
		buffer: make([]float64, 0, config.PredictorBufferSize),
		wg:     wg,
	}
}

func (p *TrendPredictor) Start(ctx context.Context, weatherChan <-chan models.WeatherData, forecastChan chan<- models.ForecastReport) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		
		readingsCount := 0

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-weatherChan:
				if !ok {
					return
				}
				
				if len(p.buffer) >= config.PredictorBufferSize {
					p.buffer = p.buffer[1:]
				}
				p.buffer = append(p.buffer, data.WindSpeed)
				readingsCount++

				if readingsCount >= config.WeatherPerGrid && len(p.buffer) > 1 {
					first := p.buffer[0]
					last := p.buffer[len(p.buffer)-1]
					trend := last - first

					forecastMsg := fmt.Sprintf("Prognoza: Prędkość wiatru zmieni się o %.2f jednostek w ciągu 1 GridStep", trend)
					report := models.ForecastReport{
						Trend:   trend,
						Message: forecastMsg,
					}

					select {
					case forecastChan <- report:
					default:
					}

					readingsCount = 0
				}
			}
		}
	}()
}