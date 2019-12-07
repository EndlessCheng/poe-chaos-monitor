package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	accountName := os.Getenv("ACCOUNT_NAME")
	poeSessionID := os.Getenv("POE_SESSION_ID")
	leagueName, err := GetCurrentLeague()
	if err != nil {
		panic(err)
	}
	//leagueName := LeagueNameStandard
	tabIndex := 0

	fmt.Println("Current League:", leagueName)

	lastRecord, err := GetLastRecord()
	if err != nil {
		panic(err)
	}
	lastMinute := 0
	prevWorth := -1
	rate := 0.0
	if lastRecord != nil {
		lastMinute = lastRecord.Minute
		prevWorth = lastRecord.Worth
		rate = float64(lastRecord.Rate)
	}

	f, err := os.OpenFile(RecordFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const timeBase = time.Minute
	const gap = 5 * timeBase
	for minute := lastMinute; ; minute += int(gap / timeBase) {
		t0 := time.Now().UnixNano()

		r, err := GetExRate(leagueName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r = rate
		} else {
			rate = r
		}
		intRate := int(math.Ceil(r))

		var currency *MainCurrency
		for {
			currency, err = GetMainCurrency(poeSessionID, leagueName, accountName, tabIndex)
			if err == nil {
				break
			}
			fmt.Fprintln(os.Stderr, err)
		}

		record := &Record{
			Minute:        minute,
			Worth:         intRate*currency.NumExaltedOrb + currency.NumChaosOrb,
			Rate:          intRate,
			NumExaltedOrb: currency.NumExaltedOrb,
			NumChaosOrb:   currency.NumChaosOrb,
		}
		fmt.Println(record.ShownString(prevWorth))

		if _, err := f.WriteString(record.FileString() + "\n"); err != nil {
			panic(err)
		}

		prevWorth = record.Worth
		time.Sleep(gap - time.Duration(time.Now().UnixNano()-t0))
	}
}
