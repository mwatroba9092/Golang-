package main

import (
	"errors"
	"fmt"
	"time"
)

// 1. STRUKTURY DANYCH

type Samolot struct {
	Model        string
	LiczbaMiejsc int
}

type Pasazer struct {
	ID    string // mozna zmienic na pesel lub inne unikalne ID
	Imie     string
	Nazwisko string
}

type Rezerwacja struct {
	IDRezerwacji string
	Pasazer      *Pasazer
	Lot          *Lot
}

type Lot struct {
	NumerLotu  string
	Skad       string
	Dokad      string
	Samolot    Samolot
	Rezerwacje map[string]Rezerwacja
}

func NowyLot(numer, skad, dokad string, samolot Samolot) *Lot {
	return &Lot{
		NumerLotu:  numer,
		Skad:       skad,
		Dokad:      dokad,
		Samolot:    samolot,
		Rezerwacje: make(map[string]Rezerwacja),
	}
}


// 2. LOGIKA BIZNESOWA I METODY

func (l *Lot) String() string {
	return fmt.Sprintf("Lot %s [%s -> %s] | Samolot: %s | Wolne miejsca: %d/%d",
		l.NumerLotu, l.Skad, l.Dokad, l.Samolot.Model, l.WolneMiejsca(), l.Samolot.LiczbaMiejsc)
}

func (l *Lot) WolneMiejsca() int {
	return l.Samolot.LiczbaMiejsc - len(l.Rezerwacje)
}

func (l *Lot) Zarezerwuj(p *Pasazer) error {
	if l.WolneMiejsca() <= 0 {
		return errors.New("brak wolnych miejsc na ten lot")
	}

	if _, istnieje := l.Rezerwacje[p.ID]; istnieje {
		return fmt.Errorf("pasażer %s %s posiada już rezerwację na ten lot", p.Imie, p.Nazwisko)
	}

	idRezerwacji := fmt.Sprintf("RES-%s-%d", l.NumerLotu, time.Now().UnixNano()%1000)
	nowaRezerwacja := Rezerwacja{
		IDRezerwacji: idRezerwacji,
		Pasazer:      p,
		Lot:          l,
	}

	l.Rezerwacje[p.ID] = nowaRezerwacja
	fmt.Printf("Pomyślnie zarezerwowano lot %s dla %s %s (Rezerwacja: %s)\n", l.NumerLotu, p.Imie, p.Nazwisko, idRezerwacji)
	return nil
}

func (l *Lot) Odwolaj(pasazerID string) error {
	if _, istnieje := l.Rezerwacje[pasazerID]; !istnieje {
		return errors.New("nie znaleziono rezerwacji dla podanego pasażera na ten lot")
	}

	delete(l.Rezerwacje, pasazerID)
	fmt.Printf("Odwołano rezerwację pasażera o numerze pesel %s na lot %s\n", pasazerID, l.NumerLotu)
	return nil
}

// 3. WYSZUKIWANIE Z UŻYCIEM INTERFEJSÓW

type SystemRezerwacji struct {
	Loty []*Lot
}

func (s *SystemRezerwacji) ZnajdzRezerwacjePasazera(pasazerID string) []Rezerwacja {
	var znalezione []Rezerwacja
	for _, lot := range s.Loty {
		if rez, istnieje := lot.Rezerwacje[pasazerID]; istnieje {
			znalezione = append(znalezione, rez)
		}
	}
	return znalezione
}

type KryteriumLotu interface {
	SpelniaKryterium(lot *Lot) bool
}

type FiltrSkad struct {
	PortLotniczy string
}

func (f FiltrSkad) SpelniaKryterium(lot *Lot) bool {
	return lot.Skad == f.PortLotniczy
}

type FiltrDokad struct {
	PortLotniczy string
}

func (f FiltrDokad) SpelniaKryterium(lot *Lot) bool {
	return lot.Dokad == f.PortLotniczy
}

func (s *SystemRezerwacji) ZnajdzLoty(k KryteriumLotu) []*Lot {
	var wynik []*Lot
	for _, lot := range s.Loty {
		if k.SpelniaKryterium(lot) {
			wynik = append(wynik, lot)
		}
	}
	return wynik
}

// 4. DEMONSTRACJA DZIAŁANIA (main)

func main() {
	boeing := Samolot{Model: "Boeing 737", LiczbaMiejsc: 2}
	airbus := Samolot{Model: "Airbus A320", LiczbaMiejsc: 150}

	lot1 := NowyLot("LO3905", "Warszawa", "Kraków", boeing)
	lot2 := NowyLot("RY112", "Gdańsk", "Londyn", airbus)
	lot3 := NowyLot("W6340", "Warszawa", "Paryż", airbus)

	system := SystemRezerwacji{
		Loty: []*Lot{lot1, lot2, lot3},
	}

	jan := &Pasazer{ID: "99010212345", Imie: "Jan", Nazwisko: "Kowalski"}
	anna := &Pasazer{ID: "95030454321", Imie: "Anna", Nazwisko: "Nowak"}
	piotr := &Pasazer{ID: "90050611111", Imie: "Piotr", Nazwisko: "Zieliński"}

	fmt.Println("--- STAN POCZĄTKOWY (Interfejs Stringer) ---")
	for _, l := range system.Loty {
		fmt.Println(l)
	}

	fmt.Println("\n--- 1 & 2. REZERWACJE I BLOKADA PODWÓJNEJ REZERWACJI ---")
	_ = lot1.Zarezerwuj(jan)
	_ = lot1.Zarezerwuj(anna)

	err := lot1.Zarezerwuj(jan)
	if err != nil {
		fmt.Printf("Zablokowano: %s\n", err)
	}

	err = lot1.Zarezerwuj(piotr)
	if err != nil {
		fmt.Printf("Zablokowano: %s\n", err)
	}

	_ = lot2.Zarezerwuj(jan)

	fmt.Println("\n--- 3. WYSZUKIWANIE REZERWACJI PASAŻERA ---")
	rezerwacjeJana := system.ZnajdzRezerwacjePasazera(jan.ID)
	fmt.Printf("Rezerwacje dla pasażera %s %s:\n", jan.Imie, jan.Nazwisko)
	for _, r := range rezerwacjeJana {
		fmt.Printf("- %s (ID: %s)\n", r.Lot.NumerLotu, r.IDRezerwacji)
	}

	fmt.Println("\n--- ODWOŁANIE REZERWACJI ---")
	_ = lot1.Odwolaj(jan.ID)
	fmt.Printf("Liczba wolnych miejsc w locie %s po odwołaniu: %d\n", lot1.NumerLotu, lot1.WolneMiejsca())

	fmt.Println("\n--- 4. WYSZUKIWANIE LOTÓW Z UŻYCIEM INTERFEJSÓW ---")

	kryteriumSkad := FiltrSkad{PortLotniczy: "Warszawa"}
	znalezioneSkad := system.ZnajdzLoty(kryteriumSkad)
	
	fmt.Println("Loty rozpoczynające się z Warszawy:")
	for _, l := range znalezioneSkad {
		fmt.Printf("- %s do %s\n", l.NumerLotu, l.Dokad)
	}

	kryteriumDokad := FiltrDokad{PortLotniczy: "Londyn"}
	znalezioneDokad := system.ZnajdzLoty(kryteriumDokad)
	
	fmt.Println("Loty z celem: Londyn:")
	for _, l := range znalezioneDokad {
		fmt.Printf("- %s z %s\n", l.NumerLotu, l.Skad)
	}
}