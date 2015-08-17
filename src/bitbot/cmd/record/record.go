/* Record exchangers orderbooks on a provided frequency (ex: every minute) */
package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"
	"time"

	"bitbot/exchanger/bitfinex"
	"bitbot/exchanger/btce"
	"bitbot/exchanger/cex"
	"bitbot/exchanger/hitbtc"
	"bitbot/exchanger/kraken"
	"bitbot/exchanger/orderbook"
	"bitbot/httpreq"
)

const pair = "BTC_USD"

type bookFunc func(string) (*orderbook.OrderBook, error)

var exchangers = map[string]bookFunc{
	"bitfinex": bitfinex.OrderBook,
	"btce":     btce.OrderBook,
	"hitbtc":   hitbtc.OrderBook,
	"cex":      cex.OrderBook,
	"kraken":   kraken.OrderBook,
}

var (
	elasticseatchHost = flag.String("e", "http://exchanger-db:9200", "Elasticsearch host address.")
	periodicity       = flag.Int64("p", 5, "Periodicity expressed in seconds.")
)

func main() {
	flag.Parse()

	// wait until ES starts
	time.Sleep(10 * time.Second)

	for {
		log.Println("Fetching orderbooks...")
		for exchanger, f := range exchangers {
			go fetchRecord(exchanger, pair, f)
		}
		time.Sleep(time.Duration(*periodicity) * time.Second)
	}
}

func fetchRecord(exchanger, pair string, f bookFunc) {
	defer logPanic()

	start := time.Now().UnixNano()
	book, err := f(pair)
	end := time.Now().UnixNano()

	if err != nil {
		log.Println(err)
		return
	}

	record := map[string]interface{}{
		"StartTime": start,
		"EndTime":   end,
		"Bids":      book.Bids[:25],
		"Asks":      book.Asks[:25],
	}

	body, err := json.Marshal(record)
	if err != nil {
		log.Println(err)
		return
	}

	var result struct {
		Id      string `json:"_id"`
		Created bool   `json:"created"`
	}

	err = httpreq.Post(*elasticseatchHost+"/exchanger/orderbook/", nil, string(body), &result)
	if err != nil {
		log.Println(err)
		return
	}

	if !result.Created {
		log.Printf("record: can't save %s orderbook in ElasticSearch.\n", exchanger)
		return
	}

	log.Printf("record: %s orderbook saved in ElasticSearch (_id: %s)\n", exchanger, result.Id)
}

// logPanic logs a formatted stack trace of the panicing goroutine. The stack trace is truncated
// at 4096 bytes (https://groups.google.com/d/topic/golang-nuts/JGraQ_Cp2Es/discussion)
func logPanic() {
	if err := recover(); err != nil {
		const size = 4096
		buf := make([]byte, size)
		stack := buf[:runtime.Stack(buf, false)]
		log.Printf("Error: %v\n%s", err, stack)
	}
}
