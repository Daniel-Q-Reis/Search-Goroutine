package main

import (
	"buscador/internal/fetcher"
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	priceChannel := make(chan float64)
	done := make(chan struct{}) // Canal de sinalização

	var wg sync.WaitGroup
	wg.Add(3)

	// Consumidor
	go func() {
		var totalPrice float64
		countPrice := 0.0
		for price := range priceChannel {
			totalPrice += price
			countPrice++
			fmt.Printf("R$ %.2f \n", price)
			fmt.Printf("Preço médio: R$ %.2f \n", (totalPrice / countPrice))
		}
		// ✅ Sinaliza que terminou
		done <- struct{}{}
	}()

	// Produtores
	go func() {
		defer wg.Done()
		priceChannel <- fetcher.FetchPriceFromSite1()
	}()
	go func() {
		defer wg.Done()
		priceChannel <- fetcher.FetchPriceFromSite2()
	}()
	go func() {
		defer wg.Done()
		priceChannel <- fetcher.FetchPriceFromSite3()
	}()

	// Fecha o canal após todos produtores terminarem

	wg.Wait()
	close(priceChannel)

	// Espera o consumidor terminar
	<-done

	fmt.Printf("\nTempo total: %s\n", time.Since(start))
}
