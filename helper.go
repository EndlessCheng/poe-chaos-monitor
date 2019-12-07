package main

import (
	"fmt"
	"github.com/levigross/grequests"
	"net/http"
	"strings"
	"time"
)

const LeagueNameStandard = "Standard"

func GetCurrentLeague() (leagueName string, err error) {
	resp, err := grequests.Get("http://api.pathofexile.com/leagues", nil)
	if err != nil {
		return
	}
	if !resp.Ok {
		return "", fmt.Errorf("resp return %d", resp.StatusCode)
	}

	d := []struct {
		ID    string  `json:"id"`
		EndAt *string `json:"endAt"`
	}{}
	if err = resp.JSON(&d); err != nil {
		return
	}
	for _, league := range d {
		if league.EndAt != nil && !strings.Contains(league.ID, " ") {
			return league.ID, nil
		}
	}
	return "", fmt.Errorf("can't find league name")
}

const (
	CurrencyExaltedOrb = "Exalted Orb"
	CurrencyChaosOrb   = "Chaos Orb"
)

type MainCurrency struct {
	NumExaltedOrb int
	NumChaosOrb   int
}

func GetMainCurrency(poeSessionID, leagueName, accountName string, tabIndex int) (currency *MainCurrency, err error) {
	api := fmt.Sprintf("https://www.pathofexile.com/character-window/get-stash-items?league=%s&accountName=%s&tabIndex=%d", leagueName, accountName, tabIndex)
	resp, err := grequests.Get(api, &grequests.RequestOptions{Cookies: []*http.Cookie{{Name: "POESESSID", Value: poeSessionID}}})
	if err != nil {
		return
	}
	if !resp.Ok {
		return nil, fmt.Errorf("resp return %d", resp.StatusCode)
	}

	d := struct {
		Items []struct {
			TypeLine  string `json:"typeLine"`
			StackSize int    `json:"stackSize"`
		} `json:"items"`
	}{}
	if err = resp.JSON(&d); err != nil {
		return
	}

	currency = &MainCurrency{}
	for _, item := range d.Items {
		switch item.TypeLine {
		case CurrencyExaltedOrb:
			currency.NumExaltedOrb = item.StackSize
		case CurrencyChaosOrb:
			currency.NumChaosOrb = item.StackSize
		}
	}
	return
}

// Buy one Ex need ? C
func GetExRate(leagueName string) (rate float64, err error) {
	api := fmt.Sprintf("https://poe.ninja/api/data/currencyoverview?league=%s&type=Currency", leagueName)
	resp, err := grequests.Get(api, &grequests.RequestOptions{
		//DialTimeout: 15 * time.Second,
	})
	if err != nil {
		return
	}
	if !resp.Ok {
		return 0, fmt.Errorf("resp return %d", resp.StatusCode)
	}

	d := struct {
		Lines []struct {
			CurrencyTypeName string `json:"currencyTypeName"`
			Receive          struct {
				SampleTimeUTC string  `json:"sample_time_utc"`
				Value         float64 `json:"value"`
			} `json:"receive"`
		} `json:"lines"`
	}{}
	if err = resp.JSON(&d); err != nil {
		return
	}

	for _, l := range d.Lines {
		if l.CurrencyTypeName == CurrencyExaltedOrb {
			t, er := time.Parse(time.RFC3339, l.Receive.SampleTimeUTC)
			if er != nil {
				return 0, er
			}
			rate = l.Receive.Value
			fmt.Printf("EX:C = %.1f (%d minutes ago)\n", rate, time.Since(t)/time.Minute)
			return
		}
	}
	return 0, nil
}
