package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	CurrencyRatesAPI = "http://api.nbp.pl/api/exchangerates"
	GoldRatesAPI = "http://api.nbp.pl/api/cenyzlota/"
)

type Client struct {
	Token string
	hc http.Client
	RemainingTime int32
}

func CreateNewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type CurrSearchResult struct {
	Page int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	TotalResults int32 `json:"total_results`
	NextPage string `json:"next_page"`
	Currencies []Currency `json:"currencies"`
}

type Currency struct {
	TopCount float32 `json:"top_count"`
	Code int32 `json:"code"`
	Bid float32 `json:"bid"`
	Ask float32 `json:"ask"`
	Mid float32 `json:"mid"`
	Table string `json:"table"`
	CurrencyName string `json:"currency_name"`
	Country string `json:"country"`
	Symbol string `json:"symbol"`
	CurrentDate string `json:"current_date"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

type GoldRate struct {
	Code int32 `json:"code"`
	TopCount float32 `json:"top_count`
	Date string `json:"date"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}

func (c *Client) SearchCurrencies(query string, perPage, page int) (*CurrSearchResult, error) {
	url := fmt.Sprintf(CurrencyRatesAPI+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	res, err := c.requestDoWithAuth("GET", url)
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result CurrSearchResult
	err = json.Unmarshall(data, &result)
	return &result, err
}

func (c *Client) GetCurrencyTable(table string) () {
	url := fmt.Sprintf(CurrencyRatesAPI+"/tables/%s/", table)
	res, err := c.RequestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result CurrSearchResult
	err = json.Unmarshall(data, &result)
	return &result, err
}

func (c *Client) GetCurrencyByTimeInterval(table, currencyCode, startDate , endDate string) () {
	url := fmt.Sprintf(CurrencyRatesAPI+"rates/%s/%s/%s/%s", table, currencyCode, startDate, endDate)
	res, err := c.RequestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result Currency
	return &result, err
}

func (c *Client) GetGoldRatesTable(tableCode string) (*[]GoldRate, error) {
	url := fmt.Sprintf(GoldRatesAPI)
	res, err := c.RequestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result []GoldRate
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetCurrency(tableId int32, code int32) (*Currency, error) {
	url := fmt.Sprintf(CurrencyRatesAPI+"/rates/%d/%d", tableId, code)
	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	var result Currency
	err = json.Unmarshal(data, &result)
	return &result, err
}

func main() {
	os.Setenv("", "")
	var token = os.Getenv("NBPToken")

	var c = NewClient(token)

	result, err := c.GetCurrencyTable()

	if err != nil {
		fmt.Errorf("Search error:%v", err)
	}

	fmt.Println(result)
}
