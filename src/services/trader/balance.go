package main

import (
	"fmt"
	"log"

	"bitbot/exchanger"
)

var minBalance = map[string]float64{
	"BTC": 0.0005,
	"ZEC": 0.001,
}

func rebalance(clients map[string]Client, arb *arbitrage, pair exchanger.Pair) error {
	balances, err := getBalances(clients)
	if err != nil {
		return err
	}

	availableSellVol := balances[arb.sellEx.Exchanger][pair.Base]
	if availableSellVol <= minBalance[pair.Base] {
		printBalances(balances, pair)

		err := transfert(
			clients["Hitbtc"],
			clients["Poloniex"],
			pair.Base,
			balances[arb.buyEx.Exchanger][pair.Base],
		)

		if err != nil {
			return err
		}
	}

	availableBuyVol := balances[arb.buyEx.Exchanger][pair.Quote] / arb.buyEx.Asks[0].Price
	if availableBuyVol <= minBalance[pair.Quote] {
		printBalances(balances, pair)

		err := transfert(
			clients["Hitbtc"],
			clients["Poloniex"],
			pair.Quote,
			balances[arb.sellEx.Exchanger][pair.Quote],
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func transfert(org, dest Client, cur string, vol float64) error {
	log.Printf("Starting transfert of %f %s from %s to %s\n", vol, cur, org.Exchanger(), dest.Exchanger())

	address, err := dest.PaymentAddress(cur)
	if err != nil {
		return err
	}

	ack, err := org.Withdraw(vol, cur, address)
	if err != nil {
		return fmt.Errorf("Cannot withdraw `%s` from %s: %s\n", cur, err, org.Exchanger())
	} else {
		log.Printf("Transfer registered: %s\n", ack)
	}

	// Wait until we see the amout on Hitbtc main account
	return dest.WaitBalance(cur)
}

func getBalances(clients map[string]Client) (map[string]map[string]float64, error) {
	out := map[string]map[string]float64{}

	for _, c := range clients {
		b, err := c.TradingBalances()
		if err != nil {
			return nil, err
		}
		out[c.Exchanger()] = b
	}

	return out, nil
}

func printBalances(balances map[string]map[string]float64, pair exchanger.Pair) {
	log.Printf("Balance: Hitbtc %s: %f, %s %f\n", pair.Base, balances["Hitbtc"][pair.Base], pair.Quote, balances["Hitbtc"][pair.Quote])
	log.Printf("Balance: Poloniex %s: %f, %s %f\n", pair.Base, balances["Poloniex"][pair.Base], pair.Quote, balances["Poloniex"][pair.Quote])

	totalBase := balances["Hitbtc"][pair.Base] + balances["Poloniex"][pair.Base]
	totalQuote := balances["Hitbtc"][pair.Quote] + balances["Poloniex"][pair.Quote]
	log.Printf("Balance: Total %s: %f, %s %f\n", pair.Base, totalBase, pair.Quote, totalQuote)
}
