package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type AssetDataResponse struct {
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type QuoteResponse struct {
	Result []AssetDataResponse `json:"result"`
}

type YahooResp struct {
	QuoteResponse QuoteResponse `json:"quoteResponse"`
}

const YahooBaseUrl = "https://query1.finance.yahoo.com/v7/finance/quote?=&symbols="

func main() {
	log.Println("main started")
	initDB()
	var assets []Asset
	err := Database.Model(&assets).Select()
	if err != nil {
		panic(err)
	}
	if len(assets) < 1 {
		log.Println("No assets found in database, exiting.")
		os.Exit(0)
	}
	c := make(chan map[string]float64)
	for _, asset := range assets {
		go getAssetPrice(asset.Ticker, c) // TODO what happens with many requests?
	}
	// keeps track of goroutine count
	count := 0
	for msg := range c {
		updateAssetPrice(msg)
		// increment after update
		count++
		if count == len(assets) {
			close(c)
		}
	}
	// calc net worth

}
func updateAssetPrice(asset map[string]float64) {
	var key string
	var value float64
	for k, v := range asset { // TODO has to be a better way..?
		key = k
		value = v
	}
	dt := time.Now()
	fmt.Println(dt)
	model := &Asset{
		Price:     value,
		UpdatedAt: dt,
	}
	var id string
	_, err := Database.Model(model).Column("price").Column("updated_at").Where("asset.ticker = ?", key).Returning("id").Update(&id)
	if err != nil {
		panic(err)
	}
}
func getAssetPrice(ticker string, c chan map[string]float64) {
	url := YahooBaseUrl + ticker
	resp, err := http.Get(url)

	data := YahooResp{}
	if err != nil {
		log.Println("Failed to get url: " + url)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}
	json.Unmarshal(body, &data)
	m := make(map[string]float64)
	m[ticker] = data.QuoteResponse.Result[0].RegularMarketPrice
	// send map back through channel
	c <- m
}
