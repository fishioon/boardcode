package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type event struct {
	Message string
	Kind    string
	Delay   int
}

type compile struct {
	Errors      string
	Events      []event
	Status      int
	IsTest      bool
	TestsFailed int
	VetOK       bool
}

func play(code string) (string, error) {
	resp, err := http.Post("https://play.golang.org/compile", "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewBufferString(code))
	if err != nil {
		return "", err
	}
	res := new(compile)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Events[0].Message, nil
}
