package components

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"grid-simulator/config"
	"grid-simulator/models"
)

type Broadcaster struct {
	subscribers []chan<- models.WeatherData
	mu          sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan<- models.WeatherData, 0),
	}
}

func (b *Broadcaster) Subscribe(ch chan<- models.WeatherData) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, ch)
}

func (b *Broadcaster) Broadcast(data models.WeatherData) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- data:
		default:
		}
	}
}

type WeatherStation struct {
	broadcaster *Broadcaster
	wg          *sync.WaitGroup
}

func NewWeatherStation(b *Broadcaster, wg *sync.WaitGroup) *WeatherStation {
	return &WeatherStation{
		broadcaster: b,
		wg:          wg,
	}
}

func (ws *WeatherStation) Start(ctx context.Context) {
	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()
		
		ticker := time.NewTicker(config.WeatherStep)
		defer ticker.Stop()

		currentWindSpeed := 15.0 

		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				change := (rnd.Float64() * 2) - 1.0
				currentWindSpeed += change

				if currentWindSpeed < 0 {
					currentWindSpeed = 0
				}

				data := models.WeatherData{
					WindSpeed: currentWindSpeed,
				}

				ws.broadcaster.Broadcast(data)
			}
		}
	}()
}