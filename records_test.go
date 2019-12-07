package main

import "testing"

func TestGetLastRecord(t *testing.T) {
	r, err := GetLastRecord()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}
