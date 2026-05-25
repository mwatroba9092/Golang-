package components

import (
	"sync"
)

// implementuje interfejs interfaces.EnergyStorage
type BatteryStorage struct {
	capacity      float64
	currentCharge float64 
	mu            sync.RWMutex
}

func NewBatteryStorage(capacityMWh float64, initialChargeMWh float64) *BatteryStorage {
	return &BatteryStorage{
		capacity:      capacityMWh,
		currentCharge: initialChargeMWh,
	}
}

// przyjmuje energię i zwraca ile faktycznie udało się przyjąć
func (b *BatteryStorage) Charge(power float64) float64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	spaceLeft := b.capacity - b.currentCharge
	if spaceLeft >= power {
		b.currentCharge += power
		return power 
	} else {
		b.currentCharge = b.capacity
		return spaceLeft
	}
}

// oddaje energięi zwraca ile faktycznie udało się oddać
func (b *BatteryStorage) Discharge(power float64) float64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.currentCharge >= power {
		b.currentCharge -= power
		return power
	} else {
		available := b.currentCharge
		b.currentCharge = 0.0
		return available
	}
}

// zwraca poziom naładowania w zakresie od 0.0 do 1.0
func (b *BatteryStorage) GetSoC() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.currentCharge / b.capacity
}