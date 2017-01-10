package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"bitbot/exchanger"
)

type transaction struct {
	orig   string
	dest   string
	amount float64
}

func rebalance(withdrawers map[string]Withdrawer, pair exchanger.Pair) {
	wg := sync.WaitGroup{}

	f := func(cur string) {
		wg.Add(1)
		defer wg.Done()
		execRebalanceTransactions(withdrawers, cur)
	}

	// execRebalanceTransactions triggers several API requests. With latency issues, the exchanger
	// could receive requests in a different order than what we sent. This involves that the nounce
	// will be invalid and a "Kraken errors: [EAPI:Invalid nonce]" can occur. To fix this quickly
	// we just wait 10 seconds here...
	go f(pair.Base)
	time.Sleep(10 * time.Second)
	go f(pair.Quote)
	wg.Wait()
}

func execRebalanceTransactions(withdrawers map[string]Withdrawer, cur string) {
	masterBal, err := getBalances(withdrawers)
	if err != nil {
		log.Printf("execRebalanceTransactions: call to getBalances() failed - %s (%s)", err, cur)
		return
	}

	curBal := map[string]float64{}
	for ex, bal := range masterBal {
		curBal[ex] = bal[cur]
	}

	var wg sync.WaitGroup
	total := map[string]float64{}

	for _, t := range findRebalanceTransactions(curBal) {
		wg.Add(1)

		go func(t *transaction) {
			defer wg.Done()
			err := execTransaction(withdrawers[t.orig], withdrawers[t.dest], cur, t.amount)
			if err != nil {
				log.Printf("execRebalanceTransactions: call to execTransaction() failed - %s (%s)", err, cur)
			} else {
				total[t.dest] += t.amount
			}
		}(t)
	}

	wg.Wait()

	for ex, amount := range total {
		// we only take 90% to remove the transaction fee
		target := 0.9 * (curBal[ex] + amount)
		err = withdrawers[ex].WaitBalance(cur, target)
		if err != nil {
			log.Printf("execRebalanceTransactions: call to waitBalanceChange() failed - %s (%s)", err, cur)
		}
	}
}

func findRebalanceTransactions(balances map[string]float64) []*transaction {
	var total float64
	for _, balance := range balances {
		total += balance
	}

	const threshold = 0.05
	targetBal := total / float64(len(balances))
	positives := map[string]float64{}
	negatives := map[string]float64{}

	for exchanger, balance := range balances {
		alloc := balance / total
		delta := balance - targetBal

		if alloc < threshold {
			negatives[exchanger] = -delta
		} else if delta > 0 {
			positives[exchanger] = delta
		}
	}

	var amount float64
	transactions := []*transaction{}

	for dest, negDelta := range negatives {
		for orig, posDelta := range positives {
			if posDelta <= 0 || negDelta == 0 {
				continue
			} else if posDelta > negDelta {
				amount = negDelta
				positives[orig] -= amount
			} else {
				amount = posDelta
				negDelta -= posDelta
				delete(positives, orig)
			}

			t := &transaction{orig, dest, amount}
			transactions = append(transactions, t)
		}
	}

	return transactions
}

func execTransaction(org, dest Withdrawer, cur string, vol float64) error {
	log.Printf("Starting transfert of %f %s from %s to %s\n", vol, cur, org.Exchanger(), dest.Exchanger())

	var address string
	var err error

	if org.Exchanger() == "Kraken" {
		// Kraken requires to input the withdrawal addresses in the UI and to
		// give them unique name. The convention is ExchangerName + "-" + cur.
		// Example: Poloniex-ZEC
		address = fmt.Sprintf("%s-%s", dest.Exchanger(), cur)
	} else {
		address, err = dest.PaymentAddress(cur)
		if err != nil {
			return err
		}
	}

	ack, err := org.Withdraw(vol, cur, address)
	if err != nil {
		return fmt.Errorf("Cannot withdraw `%s` from %s: %s\n", cur, err, org.Exchanger())
	} else {
		log.Printf("Transfer registered: %s\n", ack)
	}

	return nil
}

func getBalances(withdrawers map[string]Withdrawer) (map[string]map[string]float64, error) {
	out := map[string]map[string]float64{}

	for _, w := range withdrawers {
		b, err := w.TradingBalances()
		if err != nil {
			return nil, err
		}
		out[w.Exchanger()] = b
	}

	return out, nil
}

func printBalances(balances map[string]map[string]float64, pair exchanger.Pair) {
	var totalBase float64
	var totalQuote float64

	for ex, bal := range balances {
		totalBase += bal[pair.Base]
		totalQuote += bal[pair.Quote]
		log.Printf("Balance: %s %s: %f, %s %f\n", ex, pair.Base, bal[pair.Base], pair.Quote, bal[pair.Quote])
	}

	log.Printf("Balance: Total %s: %f, %s %f\n", pair.Base, totalBase, pair.Quote, totalQuote)
}
