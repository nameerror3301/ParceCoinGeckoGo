package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

type DataCoin struct {
	CoinName             string  `json:"name"`
	CoinSimbol           string  `json:"symbol"`
	CoinPrice            float64 `json:"current_price"`
	MarcetCapitalization int     `json:"market_cap"`
}

type Config struct {
	TypeWallet    string `yaml:"TypeWallet"`
	NumberOfCoins int64  `yaml:"NumberOfCoins"`
	PathToFolder  string `yaml:"PathToFolder"`
}

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
			log.Println(err)
		}

		for _, val := range data {
			nameList = append(nameList, val.CoinName)
			priceList = append(priceList, val.CoinPrice)
			simbolList = append(simbolList, val.CoinSimbol)
			capList = append(capList, val.MarcetCapitalization)
		}
		i++

		// Add more process ... wait..
		for idxN, valN := range nameList {
			f.SetCellValue("CoinList", fmt.Sprintf("A%d", idxN+1), valN)
		}
		for idxS, valS := range simbolList {
			f.SetCellValue("CoinList", fmt.Sprintf("B%d", idxS+1), valS)
		}
		for idxP, valP := range priceList {
			f.SetCellValue("CoinList", fmt.Sprintf("C%d", idxP+1), valP)
		}
		for idxC, valC := range capList {
			f.SetCellValue("CoinList", fmt.Sprintf("D%d", idxC+1), valC)
		}
	}
	if err := f.SaveAs("Coins.xlsx"); err != nil {
		fmt.Println(err)
	}
	// 2.325321123s or 3.132414124s or 4.015732478s
	fmt.Println(time.Since(timeStart))
}

// Parce userConfig yaml
func ParceUserConfig() (error, string, int64, string) {
	var conf Config
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		return fmt.Errorf("err read config file - %s", err), "", 0, ""
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return fmt.Errorf("err parce config file - %s", err), "", 0, ""
	}

	return nil, conf.TypeWallet, conf.NumberOfCoins, conf.PathToFolder
}
