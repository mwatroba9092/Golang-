package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"grid-simulator/components"
	"grid-simulator/models"
)

func main() {
	fmt.Println("Uruchamianie symulatora sieci energetycznej...")

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nZłapano sygnał przerwania. Inicjacja Graceful Shutdown...")
		cancel()
	}()

	// 1. Inicjalizacja Loggera
	logger, err := components.NewCSVLogger(&wg)
	if err != nil {
		fmt.Printf("Błąd inicjalizacji loggera: %v\n", err)
		return
	}
	logger.Start(ctx)

	// 2. Kanały komunikacyjne
	weatherChanForOZE := make(chan models.WeatherData, 10)
	weatherChanForPred := make(chan models.WeatherData, 10)
	forecastChan := make(chan models.ForecastReport, 1)
	demandChan := make(chan models.DemandReport, 100)

	// 3. Stacja pogodowa i Pub/Sub
	broadcaster := components.NewBroadcaster()
	broadcaster.Subscribe(weatherChanForOZE)
	broadcaster.Subscribe(weatherChanForPred)
	
	weatherStation := components.NewWeatherStation(broadcaster, &wg)
	weatherStation.Start(ctx)

	// 4. Inteligencja (Predictor)
	predictor := components.NewPredictor(&wg)
	predictor.Start(ctx, weatherChanForPred, forecastChan)

	// 5. Źródła Energii i Baterie
	oze := components.NewWindFarm(&wg)
	oze.Start(ctx, weatherChanForOZE)

	coal := components.NewCoalPlant(150.0, &wg)

	battery := components.NewBatteryStorage(50.0, 25.0)

	// 6. Gorutyna Centralna (GridHub)
	gridHub := components.NewGridHub(oze, coal, battery, logger, demandChan, forecastChan, &wg)
	gridHub.Start(ctx)

	// 7. Konsumenci energii
	residential := components.NewResidentialConsumer("Kowalski_Dom", logger, &wg)
	residential.StartResidential(ctx, demandChan)

	industrial := components.NewIndustrialConsumer("Fabryka_Stali", logger, &wg)
	industrial.StartIndustrial(ctx, demandChan)

	critical := components.NewCriticalConsumer("Szpital_Miejski", logger, &wg)
	critical.StartCritical(ctx, demandChan)

	logger.LogEvent("Symulacja uruchomiona pomyślnie. Utworzono wszystkie gorutyny.")
	fmt.Println("Symulacja działa. Wciśnij Ctrl+C aby zakończyć z Graceful Shutdown.")

	<-ctx.Done()

	fmt.Println("Oczekiwanie na bezpieczne zamknięcie gorutyn i zrzut logów do dysku...")
	wg.Wait()
	fmt.Println("Symulator zamknięty pomyślnie. Sprawdź folder logs/grid_history.csv.")
}