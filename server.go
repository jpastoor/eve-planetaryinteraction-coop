package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"github.com/jinzhu/gorm"
	"encoding/json"
)

type Server struct {
	http *http.Client
	r    *mux.Router
	db   *gorm.DB
}

func NewServer(r *mux.Router, db *gorm.DB) (*Server) {

	s := &Server{
		http: &http.Client{},
		r:    r,
		db:   db,
	}

	r.HandleFunc("/parse", s.parseLog).Methods(http.MethodPost)
	r.HandleFunc("/transactions", s.GetTransactions).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{transactionId}/markForCorp", s.MarkTransactionForCorp).Methods(http.MethodPost)
	r.HandleFunc("/ledgers", s.GetLedgers).Methods(http.MethodGet)
	r.HandleFunc("/ledgers/reset", s.ResetLedgers).Methods(http.MethodPost)

	return s
}

func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var ts []Transaction
	s.db.Find(&ts)

	bytes, err := json.Marshal(ts)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Could not convert transactions to JSON: %s", err)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (s *Server) MarkTransactionForCorp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	transactionId := vars["transactionId"]

	t := Transaction{}
	s.db.Where("id = ?", transactionId).First(&t)

	if &t == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not find transaction with supplied ID"))
	}

	s.db.Model(&t).Update("MarkedForCorp", true)

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetLedgers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) ResetLedgers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
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
