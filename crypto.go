package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiUrl = "https://api.coingecko.com/api/v3/coins/markets"

type CoinData struct {
	Id           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}

func getCoinData(currency string) ([]CoinData, error) {
	url := fmt.Sprintf("%s?vs_currency=%s&order=market_cap_desc&per_page=250&page=1", apiUrl, currency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var coinData []CoinData
	if err := json.Unmarshal(bodyBytes, &coinData); err != nil {
		return nil, err
	}

	return coinData, nil
}

func getSpecificCoinData(currency string, symbol string) (*CoinData, error) {
	coinData, err := getCoinData(currency)
	if err != nil {
		return nil, err
	}

	for _, coin := range coinData {
		if coin.Symbol == symbol {
			return &coin, nil
		}
	}

	return nil, fmt.Errorf("Криптовалюта с символом %s не найдена", symbol)
}

func main() {
	currency := "usd" // валюта по умолчанию
	interval := 10 * time.Minute

	for {
		// Запрос пользователя для ввода символа криптовалюты
		fmt.Print("Введите символ криптовалюты (например, btc): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		symbolToFind := strings.TrimSpace(scanner.Text())

		if symbolToFind == "" {
			fmt.Println("Символ не может быть пустым. Попробуйте еще раз.")
			continue
		}

		coin, err := getSpecificCoinData(currency, symbolToFind)
		if err != nil {
			fmt.Printf("Ошибка при получении данных: %s\n", err)
		} else {
			fmt.Printf("%s (%s): $%.2f\n", coin.Name, coin.Symbol, coin.CurrentPrice)
		}

		time.Sleep(interval)
	}
}
