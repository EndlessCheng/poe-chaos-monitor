package main

import (
	"os"
	"testing"
)

func TestGetCurrentLeague(t *testing.T) {
	leagueName, err := GetCurrentLeague()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(leagueName)
}

func TestGetMainCurrency(t *testing.T) {
	leagueName, err := GetCurrentLeague()
	if err != nil {
		t.Fatal(err)
	}
	leagueName = LeagueNameStandard
	currency, err := GetMainCurrency(os.Getenv("POE_SESSION_ID"), leagueName, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(currency)
}

func TestGetExRate(t *testing.T) {
	leagueName, err := GetCurrentLeague()
	if err != nil {
		t.Fatal(err)
	}
	r, err := GetExRate(leagueName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}
