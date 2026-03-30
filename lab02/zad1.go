package main

import (
	"fmt"
	"sort"
)

type Uczestnik struct {
	ID           int
	ImieNazwisko string
	Repertuar    []string
	Oceny        map[string][]float64
}


// 1. Przypisuje repertuar uczestnikowi
func PrzypiszRepertuar(u Uczestnik, utwory []string) Uczestnik {
	nowyUczestnik := u
	nowyUczestnik.Repertuar = append([]string{}, utwory...)
	nowyUczestnik.Oceny = make(map[string][]float64)

	for utwor, oceny := range u.Oceny {
		nowyUczestnik.Oceny[utwor] = append([]float64{}, oceny...)
	}

	return nowyUczestnik
}

// 2. Przypisuje pojedynczą ocenę za dany utwór
func DodajOcene(u Uczestnik, utwor string, ocena float64) Uczestnik {
	if ocena < 0 {
		ocena = 0
	} else if ocena > 25 {
		ocena = 25
	}

	nowyUczestnik := u
	nowyUczestnik.Oceny = make(map[string][]float64)

	for k, v := range u.Oceny {
		nowyUczestnik.Oceny[k] = append([]float64{}, v...)
	}

	nowyUczestnik.Oceny[utwor] = append(nowyUczestnik.Oceny[utwor], ocena)

	return nowyUczestnik
}

// 3. Wylicza średnią ocen z mechanizmem korekcyjnym (korekta skrajnych ocen)
func WyliczSredniaZKorekta(oceny []float64) float64 {
	if len(oceny) == 0 {
		return 0.0
	}

	suma := 0.0
	for _, ocena := range oceny {
		suma += ocena
	}
	srednia := suma / float64(len(oceny))

	dopuszczalneOdchylenie := 2.0
	skorygowanaSuma := 0.0

	for _, ocena := range oceny {
		if ocena > srednia+dopuszczalneOdchylenie {
			skorygowanaSuma += srednia + dopuszczalneOdchylenie
		} else if ocena < srednia-dopuszczalneOdchylenie {
			skorygowanaSuma += srednia - dopuszczalneOdchylenie
		} else {
			skorygowanaSuma += ocena
		}
	}

	return skorygowanaSuma / float64(len(oceny))
}

// Funkcja pomocnicza: liczy łączne punkty uczestnika za wszystkie utwory
func LacznePunkty(u Uczestnik) float64 {
	suma := 0.0
	for _, utwor := range u.Repertuar {
		suma += WyliczSredniaZKorekta(u.Oceny[utwor])
	}
	return suma
}

// 4. Sortuje uczestników po zdobytych punktach (malejąco)
func PosortujUczestnikow(uczestnicy []Uczestnik) []Uczestnik {
	posortowani := append([]Uczestnik{}, uczestnicy...)

	sort.Slice(posortowani, func(i, j int) bool {
		return LacznePunkty(posortowani[i]) > LacznePunkty(posortowani[j])
	})

	return posortowani
}

// 5. Znajduje uczestnika z najlepszym wynikiem za podany utwór
func NajlepszyWUtworze(uczestnicy []Uczestnik, szukanyUtwor string) (Uczestnik, float64) {
	var najlepszy Uczestnik
	najlepszyWynik := -1.0

	for _, u := range uczestnicy {
		if oceny, maUtwor := u.Oceny[szukanyUtwor]; maUtwor {
			wynik := WyliczSredniaZKorekta(oceny)
			if wynik > najlepszyWynik {
				najlepszyWynik = wynik
				najlepszy = u
			}
		}
	}

	return najlepszy, najlepszyWynik
}

func main() {
	utwor1 := "Etiuda c-moll"
	utwor2 := "Ballada g-moll"
	utwor3 := "Polonez As-dur"
	repertuar := []string{utwor1, utwor2, utwor3}

	uczestnicy := []Uczestnik{
		{ID: 1, ImieNazwisko: "Jan Kowalski"},
		{ID: 2, ImieNazwisko: "Anna Nowak"},
		{ID: 3, ImieNazwisko: "Piotr Wiśniewski"},
	}

	for i := range uczestnicy {
		uczestnicy[i] = PrzypiszRepertuar(uczestnicy[i], repertuar)
	}

	// Oceny Jana
	u := uczestnicy[0]
	u = DodajOcene(u, utwor1, 20); u = DodajOcene(u, utwor1, 21); u = DodajOcene(u, utwor1, 10); u = DodajOcene(u, utwor1, 20); u = DodajOcene(u, utwor1, 22)
	u = DodajOcene(u, utwor2, 24); u = DodajOcene(u, utwor2, 23); u = DodajOcene(u, utwor2, 25); u = DodajOcene(u, utwor2, 24); u = DodajOcene(u, utwor2, 24)
	u = DodajOcene(u, utwor3, 18); u = DodajOcene(u, utwor3, 19); u = DodajOcene(u, utwor3, 18); u = DodajOcene(u, utwor3, 20); u = DodajOcene(u, utwor3, 19)
	uczestnicy[0] = u

	// Oceny Anny
	u = uczestnicy[1]
	u = DodajOcene(u, utwor1, 22); u = DodajOcene(u, utwor1, 23); u = DodajOcene(u, utwor1, 22); u = DodajOcene(u, utwor1, 24); u = DodajOcene(u, utwor1, 23)
	u = DodajOcene(u, utwor2, 15); u = DodajOcene(u, utwor2, 16); u = DodajOcene(u, utwor2, 14); u = DodajOcene(u, utwor2, 15); u = DodajOcene(u, utwor2, 25)
	u = DodajOcene(u, utwor3, 21); u = DodajOcene(u, utwor3, 22); u = DodajOcene(u, utwor3, 21); u = DodajOcene(u, utwor3, 20); u = DodajOcene(u, utwor3, 21)
	uczestnicy[1] = u

	// Oceny Piotra
	u = uczestnicy[2]
	u = DodajOcene(u, utwor1, 19); u = DodajOcene(u, utwor1, 19); u = DodajOcene(u, utwor1, 20); u = DodajOcene(u, utwor1, 18); u = DodajOcene(u, utwor1, 19)
	u = DodajOcene(u, utwor2, 20); u = DodajOcene(u, utwor2, 21); u = DodajOcene(u, utwor2, 20); u = DodajOcene(u, utwor2, 22); u = DodajOcene(u, utwor2, 21)
	u = DodajOcene(u, utwor3, 24); u = DodajOcene(u, utwor3, 24); u = DodajOcene(u, utwor3, 25); u = DodajOcene(u, utwor3, 23); u = DodajOcene(u, utwor3, 24)
	uczestnicy[2] = u

	// Wyniki 

	fmt.Println("--- RANKING UCZESTNIKÓW ---")
	ranking := PosortujUczestnikow(uczestnicy)
	for i, ucz := range ranking {
		fmt.Printf("%d. %s - %.2f punktów\n", i+1, ucz.ImieNazwisko, LacznePunkty(ucz))
	}

	fmt.Println("\n--- NAJLEPSZY WYKONAWCA KONKRETNEGO UTWORU ---")
	szukanyUtwor := "Polonez As-dur"
	zwyciezca, wynik := NajlepszyWUtworze(uczestnicy, szukanyUtwor)
	fmt.Printf("Utwór: %s\nNajlepszy: %s (%.2f pkt)\n", szukanyUtwor, zwyciezca.ImieNazwisko, wynik)
}