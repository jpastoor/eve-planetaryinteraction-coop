package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Server struct {
	http *http.Client
	r    *mux.Router
	db   *gorm.DB
}

func NewServer(r *mux.Router, db   *gorm.DB) (*Server) {

	s := &Server{
		http: &http.Client{},
		r:    r,
		db: db,
	}

	r.HandleFunc("/parse", s.parseLog).Methods(http.MethodPost)

	return s
}

func (s *Server) parseLog(w http.ResponseWriter, r *http.Request) {

	bytes, _ := ioutil.ReadAll(r.Body)
	tp := NewTransactionParser()
	ts, errs := tp.Parse(string(bytes))

	for _, t := range ts {
		s.db.Create(t)
	}

	fmt.Print(ts)
	fmt.Print(errs)

	// TODO Add response
}
