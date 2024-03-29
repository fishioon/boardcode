package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type server struct {
	mux  *http.ServeMux
	tess *Tess
}

func newServer() (*server, error) {
	t, err := NewTess()
	if err != nil {
		return nil, err
	}
	s := &server{mux: http.NewServeMux(), tess: t}

	s.init()
	return s, nil
}

type response struct {
	Errors string      `json:"errors"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (s *server) init() {
	s.mux.HandleFunc("/compile/image", s.handleImage)
	s.mux.HandleFunc("/compile/code", s.handleCode)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	s.mux.Handle("/static/", staticHandler)
}

func (s *server) handleImage(w http.ResponseWriter, r *http.Request) {
	/*
		req := struct {
			Image []byte `json:"image"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("error decoding request: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	*/
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	code, err := s.tess.Image2text(b)
	if err != nil {
		log.Printf("error image orc request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res, err := play(code)
	if err != nil {
		log.Printf("error decoding request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	resp := &response{
		Data: &struct {
			Compile string `json:"compile"`
		}{res},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy(w, &buf): %v", err)
		return
	}
}

func (s *server) handleCode(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Code string `json:"code"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res, err := play(req.Code)
	if err != nil {
		log.Printf("error decoding request: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resp := &response{
		Data: &struct {
			Compile string `json:"compile"`
		}{res},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy(w, &buf): %v", err)
		return
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Forwarded-Proto") == "http" {
		r.URL.Scheme = "https"
		r.URL.Host = r.Host
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
	}
	s.mux.ServeHTTP(w, r)
}
