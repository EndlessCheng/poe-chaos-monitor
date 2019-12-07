package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var RecordFileName = "currency_records.txt"

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func init() {
	if err := touchFile(RecordFileName); err != nil {
		panic(err)
	}
}

type Record struct {
	Minute        int
	Worth         int
	Rate          int // Ceil(rate)
	NumExaltedOrb int
	NumChaosOrb   int
}

func (r *Record) ShownString(prevWorth int) string {
	s := fmt.Sprintf("%05d %d", r.Minute, r.Worth)
	if prevWorth != -1 && r.Worth != prevWorth {
		s += fmt.Sprintf(" (%+d)", r.Worth-prevWorth)
	}
	return s
}

func (r *Record) FileString() string {
	return fmt.Sprintf("%d %d %d %d %d", r.Minute, r.Worth, r.Rate, r.NumExaltedOrb, r.NumChaosOrb)
}

func parseRecord(line string) (r *Record, err error) {
	splits := strings.Split(strings.TrimSpace(line), " ")
	nonEmpty := []string{}
	for _, s := range splits {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	if len(nonEmpty) != 5 {
		return nil, fmt.Errorf("invalid record %v", nonEmpty)
	}

	minute, _ := strconv.Atoi(nonEmpty[0])
	worth, _ := strconv.Atoi(nonEmpty[1])
	rate, _ := strconv.Atoi(nonEmpty[2])
	numExaltedOrb, _ := strconv.Atoi(nonEmpty[3])
	numChaosOrb, _ := strconv.Atoi(nonEmpty[4])
	return &Record{minute, worth, rate, numExaltedOrb, numChaosOrb}, nil
}

func GetLastRecord() (r *Record, err error) {
	data, err := ioutil.ReadFile(RecordFileName)
	if err != nil {
		return
	}
	if len(data) == 0 {
		return
	}

	lines := strings.Split(string(data), "\n")
	lastLine := lines[len(lines)-1]
	if lastLine == "" {
		lastLine = lines[len(lines)-2]
	}
	return parseRecord(lastLine)
}
