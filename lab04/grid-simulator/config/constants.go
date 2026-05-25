package config

import "time"

const (
	// Skale czasowe
	WeatherStep = 5 * time.Millisecond   // 1 krok = ~5 minut czasu
	GridStep    = 100 * time.Millisecond // 1 krok = 1 godzina czasu
	
	// Zależności czasowe
	WeatherPerGrid      = int(GridStep / WeatherStep) // 20 kroków pogodowych w jednym sieciowym
	ForecastHorizon     = 5                           // Prognoza na 5 kroków w przód
	PredictorBufferSize = WeatherPerGrid              // Bufor równy jednej godzinie

	// Priorytety konsumentów
	PriorityCritical    = 1
	PriorityIndustrial  = 2
	PriorityResidential = 3
)