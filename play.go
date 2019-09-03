package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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
	log.Printf("code: %s", code)
	req := url.Values{"body": {code}, "withVet": {"true"}, "version": {"2"}}
	resp, err := http.PostForm("https://play.golang.org/compile", req)
	if err != nil {
		return "", err
	}
	res := new(compile)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	log.Printf("compile: %+v", res)
	return res.Events[0].Message, nil
}
