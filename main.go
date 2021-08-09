package main

import (
	"encoding/json"
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
	for {
		log.Println("updating asset prices]")
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
			if len(msg) > 0 { // if an empty map is returned that means there was an error getting the data
				updateAssetPrice(msg)
			} else {
				log.Println("Failed to update asset")
			}
			// increment after update
			count++
			if count == len(assets) {
				close(c)
			}
		}
		log.Println("finished updating assets")
		time.Sleep(2 * time.Minute)

	}

}
func updateAssetPrice(asset map[string]float64) {
	var key string
	var value float64
	for k, v := range asset { // TODO has to be a better way..?
		key = k
		value = v
	}
	dt := time.Now()
	model := &Asset{
		Price:     value,
		UpdatedAt: dt,
	}
	var id string
	var price string
	var qty string
	_, err := Database.Model(model).Column("price").Column("updated_at").Where("asset.ticker = ?", key).Returning("id, price, qty").Update(&id, &price, &qty)
	if err != nil {
		panic(err)
	}

	tableName := "asset_" + id
	createTimeseriesAssetTable(tableName)
	insertIntoTimeSeriesTable(tableName, price, qty)
}
func createTimeseriesAssetTable(tableName string) {
	_, err := Database.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (price NUMERIC, qty NUMERIC, updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		panic(err)
	}
}
func insertIntoTimeSeriesTable(tableName string, price string, qty string) {
	_, err := Database.Exec("INSERT INTO " + tableName + " (price, qty) VALUES (" + price + ", " + qty + ");")
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
	if len(data.QuoteResponse.Result) == 0 { // if yahoo sends back an empty list return an empty map
		c <- m
		return
	}
	m[ticker] = data.QuoteResponse.Result[0].RegularMarketPrice
	// send map back through channel
	c <- m
}
