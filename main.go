package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xuri/excelize/v2"
)

type DataCoin struct {
	CoinName             string  `json:"name"`
	CoinSimbol           string  `json:"symbol"`
	CoinPrice            float64 `json:"current_price"`
	MarcetCapitalization int     `json:"market_cap"`
}

func main() {
	f := excelize.NewFile()
	var nameList []string
	var simbolList []string
	var priceList []float64
	var capList []int
	var data []DataCoin
	for i := 1; i < 11; {
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=%d&sparkline=false", i)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}

		for _, val := range data {
			//fmt.Printf("Name - [%s] Simbol - [%s] Price - [%f] Capitalize - [%f]\n", val.CoinName, val.CoinSimbol, val.PriceCurrent, val.MarcetCapitalization)
			nameList = append(nameList, val.CoinName)
			priceList = append(priceList, val.CoinPrice)
			simbolList = append(simbolList, val.CoinSimbol)
			capList = append(capList, val.MarcetCapitalization)
		}
		i++
		for idxN, valN := range nameList {
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", idxN+1), valN)
		}
		for idxS, valS := range simbolList {
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", idxS+1), valS)
		}
		for idxP, valP := range priceList {
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", idxP+1), valP)
		}
		for idxC, valC := range capList {
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", idxC+1), valC)
		}
	}
	if err := f.SaveAs("Coins.xlsx"); err != nil {
		fmt.Println(err)
	}
}
