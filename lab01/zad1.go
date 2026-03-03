package main

import (
	"fmt"
	"math/rand"
)

func main() {
	N := 5
	k := 2
	liczbaRozgrywek := 100000

	if k > N-2 {
		fmt.Println("Błąd: Prowadzący nie może otworzyć tylu pudeł!")
		return
	}

	sukcesBezZmiany := 0
	sukcesPoZmianie := 0

	for i := 0; i < liczbaRozgrywek; i++ {
		pudloZnagroda := rand.Intn(N)
		pierwszyWybor := rand.Intn(N)

		if pierwszyWybor == pudloZnagroda {
			sukcesBezZmiany++
		}

		var dostepneDlaProwadzacego []int
		for j := 0; j < N; j++ {
			if j != pierwszyWybor && j != pudloZnagroda {
				dostepneDlaProwadzacego = append(dostepneDlaProwadzacego, j)
			}
		}

		var pokazanePuste []int
		for j := 0; j < k; j++ {
			pokazanePuste = append(pokazanePuste, dostepneDlaProwadzacego[j])
		}

		var alternatywnePudla []int
		for j := 0; j < N; j++ {
			if j == pierwszyWybor {
				continue
			}

			czyPokazane := false
			for _, p := range pokazanePuste {
				if j == p {
					czyPokazane = true
					break
				}
			}

			if !czyPokazane {
				alternatywnePudla = append(alternatywnePudla, j)
			}
		}

		ostatecznaDecyzja := alternatywnePudla[rand.Intn(len(alternatywnePudla))]

		if ostatecznaDecyzja == pudloZnagroda {
			sukcesPoZmianie++
		}
	}

	fmt.Printf("Przeprowadzono %d gier (Pudeł: %d, Prowadzący otwiera: %d)\n", liczbaRozgrywek, N, k)
	fmt.Println("---------------------------------------------------")

	skutecznoscBrakZmiany := float64(sukcesBezZmiany) / float64(liczbaRozgrywek) * 100
	skutecznoscZmiana := float64(sukcesPoZmianie) / float64(liczbaRozgrywek) * 100

	fmt.Printf("Wygrane (BRAK ZMIANY): %d (%.2f%%)\n", sukcesBezZmiany, skutecznoscBrakZmiany)
	fmt.Printf("Wygrane (ZMIANA):      %d (%.2f%%)\n", sukcesPoZmianie, skutecznoscZmiana)
}