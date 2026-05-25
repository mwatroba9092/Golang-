package interfaces

import (
	"context"
	"grid-simulator/models"
)

// podstawowe zachowania każdego źródła energii
type EnergySource interface {
	Start(ctx context.Context)
	Stop()
	GetPower() float64
}

// analizuje dane historyczne w celu przewidywania przyszłej produkcji
type Predictor interface {
	Start(ctx context.Context, weatherChan <-chan models.WeatherData, forecastChan chan<- models.ForecastReport)
}

// definiuje podstawowe zachowania odbiorcy końcowego
type Consumer interface {
	Start(ctx context.Context, demandChan chan<- models.DemandReport)
	GetPriority() int
}

// definiuje operacje na magazynach energii (ESS)
type EnergyStorage interface {
	Charge(power float64) float64
	Discharge(power float64) float64
	GetSoC() float64
}

// definiuje źródło danych pogodowych
type WeatherProvider interface {
	Start(ctx context.Context, broadcastChan chan<- models.WeatherData)
}

// pozwala na trwały zapis stanu systemu (CSV/JSON)
type DataLogger interface {
	Start(ctx context.Context)
	LogEvent(event interface{}) 
	Flush() error
}