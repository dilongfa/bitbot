// docker build -t bitbot-img . && docker run --rm bitbot-img
package main

import (
	"exchanger/bitfinex"
	"exchanger/hitbtc"
	"exchanger/kraken"
	"fmt"
)

func main() {
	krakenBook, err := kraken.FetchOrderBook("XXBTXLTC")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", krakenBook.Asks[0].Price)
	}

	hitbtcBook, err := hitbtc.FetchOrderBook("LTCBTC")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", hitbtcBook.Asks[0].Price)
	}

	bitfinexBook, err := bitfinex.FetchOrderBook("LTCBTC")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", bitfinexBook.Asks[0].Price)
	}
}
