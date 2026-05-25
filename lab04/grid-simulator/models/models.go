package models

//żądanie wysyłane przez konsumenta do GridHub
type DemandReport struct {
	ID       string
	PDemand  float64 
	Priority int
	ReplyTo  chan<- SupplyStatus
}

//odpowiedź z GridHub do konsumenta
type SupplyStatus struct {
	AllocatedMW float64
	Reason      string 
}

//paczka danych rozsyłana przez Broadcastera
type WeatherData struct {
	WindSpeed float64
	SunIntensity float64
}

//prognoza wysyłana z Predictora do GridHub
type ForecastReport struct {
	Trend     float64
	Message   string
}