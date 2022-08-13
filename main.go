package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
)

type DataCoin struct {
	CoinName             string  `json:"name"`
	CoinSimbol           string  `json:"symbol"`
	CoinPrice            float64 `json:"current_price"`
	MarcetCapitalization int     `json:"market_cap"`
}

/*Add the ability to configure the ability to get information about the coins.*/

func main() {
	timeStart := time.Now()
	f := excelize.NewFile()
	var nameList []string
	var simbolList []string
	var priceList []float64
	var capList []int
	var data []DataCoin

	// Fix this Crutch
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
			nameList = append(nameList, val.CoinName)
			priceList = append(priceList, val.CoinPrice)
			simbolList = append(simbolList, val.CoinSimbol)
			capList = append(capList, val.MarcetCapitalization)
		}
		i++

		// Add more process
		for idxN, valN := range nameList {
			f.SetCellValue("List", fmt.Sprintf("A%d", idxN+1), valN)
		}
		for idxS, valS := range simbolList {
			f.SetCellValue("List", fmt.Sprintf("B%d", idxS+1), valS)
		}
		for idxP, valP := range priceList {
			f.SetCellValue("List", fmt.Sprintf("C%d", idxP+1), valP)
		}
		for idxC, valC := range capList {
			f.SetCellValue("List", fmt.Sprintf("D%d", idxC+1), valC)
		}
	}
	if err := f.SaveAs("Coins.xlsx"); err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Since(timeStart))
}
