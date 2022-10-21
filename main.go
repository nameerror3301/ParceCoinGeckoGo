package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

type DataCoin struct {
	CoinName             string  `json:"name"`
	CoinSimbol           string  `json:"symbol"`
	CoinPrice            float64 `json:"current_price"`
	MarcetCapitalization float32 `json:"market_cap"`
}

type Config struct {
	TypeWallet     string `yaml:"TypeWallet"`
	NumbersOfCoins int    `yaml:"NumbersOfCoins"`
	PathToFolder   string `yaml:"PathToFolder"`
}

var (
	NameList   []string
	SimbolList []string
	PriceList  []float64
	CapList    []float32
)

func main() {
	timeStart := time.Now()
	f := excelize.NewFile()
	var data []DataCoin

	err, typeWallet, numCoins, pathToSave := ParceUserConfig()
	if err != nil {
		log.Fatal(err)
	}

	if typeWallet == "eur" || typeWallet == "usd" {
		fmt.Printf("The currency of your choice - %s\n", typeWallet)
	} else {
		log.Fatal("The currency you selected was not recognized")
	}

	if numCoins >= 1 || numCoins <= 130 {
		fmt.Printf("Number of pages from which the information will be collected - %d\n", numCoins)
	} else {
		log.Fatal("The number of pages can not be more than 130 or less than/equal to zero.")
	}

	var wg sync.WaitGroup
	for i := 1; i < numCoins+1; {
		wg.Add(5)
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&order=market_cap_desc&per_page=100&page=%d&sparkline=false", typeWallet, i)
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
			log.Println("The maximum number of requests per minute was exceeded. Performing blocking bypass...")
			time.Sleep(90 * time.Second)
		}
		fmt.Printf("Visited page - %d of %d\n", i, numCoins)

		go func() {
			for _, val := range data {
				NameList = append(NameList, val.CoinName)
				PriceList = append(PriceList, val.CoinPrice)
				SimbolList = append(SimbolList, val.CoinSimbol)
				CapList = append(CapList, val.MarcetCapitalization)
			}
			wg.Done()
		}()

		go func() {
			for idxN, valN := range NameList {
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", idxN+1), valN)
			}
			fmt.Println("1")
			wg.Done()
		}()

		go func() {
			for idxS, valS := range SimbolList {
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", idxS+1), valS)
			}
			fmt.Println("2")
			wg.Done()
		}()

		go func() {
			for idxP, valP := range PriceList {
				f.SetCellValue("Sheet1", fmt.Sprintf("C%d", idxP+1), valP)
			}
			fmt.Println("3")
			wg.Done()
		}()

		go func() {
			for idxC, valC := range CapList {
				f.SetCellValue("Sheet1", fmt.Sprintf("D%d", idxC+1), valC)
			}
			fmt.Println("4")
			wg.Done()
		}()
		i++
	}
	wg.Wait()
	if pathToSave == "" {
		if err := f.SaveAs("Coins.xlsx"); err != nil {
			log.Fatalf("Excel file creation error - %s", err)
		}
	} else {
		if err := f.SaveAs(filepath.Join(pathToSave, "Coins.xlsx")); err != nil {
			log.Fatalf("Excel file creation error - %s", err)
		}
	}
	fmt.Println("Total execution time - ", time.Since(timeStart))
}

// Parce user сonfig...
func ParceUserConfig() (error, string, int, string) {
	var conf Config
	file, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		return fmt.Errorf("err read config file - %s", err), "", 0, ""
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return fmt.Errorf("err parce config file - %s", err), "", 0, ""
	}

	return nil, conf.TypeWallet, conf.NumbersOfCoins, conf.PathToFolder
}
