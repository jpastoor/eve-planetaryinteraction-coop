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

	r.HandleFunc("/", s.GetRoot).Methods(http.MethodGet)
	r.HandleFunc("/parse", s.parseLog).Methods(http.MethodPost)
	r.HandleFunc("/transactions", s.GetTransactions).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{transactionId}/markForCorp", s.MarkTransactionForCorp).Methods(http.MethodPost)
	r.HandleFunc("/inventory", s.GetInventory).Methods(http.MethodGet)
	r.HandleFunc("/ledger", s.GetLedger).Methods(http.MethodGet)
	r.HandleFunc("/ledger/reset", s.ResetLedger).Methods(http.MethodPost)
	r.HandleFunc("/players/{playerName}", s.UpdatePlayer).Methods(http.MethodPut)

	return s
}

func (s *Server) GetRoot(w http.ResponseWriter, r *http.Request) {

	t := Transaction{}
	s.db.Order("creation_date DESC").First(&t)

	w.Header().Set("Content-Type", "text/plain")
	output := fmt.Sprintf("Server is working. Last log is from: %s (%s - %s %dx %s)", t.CreationDate, t.PlayerName, t.Action, t.Quantity, t.TypeName)
	w.Write([]byte(output))
}

func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var ts []Transaction

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

func (s *Server) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetInventory(w http.ResponseWriter, r *http.Request) {

	var ts []Transaction
	s.db.Find(&ts)

	handler := &Handler{
		inv:           NewInventory(),
		ledger:        NewLedger(&EveMarketerAPI{client: &http.Client{}}),
		typeFetcher:   NewDbTypeFetcher(s.db),
		playerFetcher: NewDbPlayerFetcher(s.db),
	}

	handler.Process(ts)

	var rsp []GetInventoryRspItem

	for typeId, contents := range handler.inv.contents {
		ty, _ := handler.typeFetcher.getTypeById(typeId)

		rsp = append(rsp, GetInventoryRspItem{
			TypeId:   typeId,
			TypeName: ty.TypeName,
			Stacks:   contents,
		})
	}

	body, _ := json.Marshal(rsp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type GetInventoryRspItem struct {
	TypeId   int
	TypeName string
	Stacks   []InventoryStack
}

func (s *Server) GetLedger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) ResetLedger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) parseLog(w http.ResponseWriter, r *http.Request) {

	bytes, _ := ioutil.ReadAll(r.Body)
	tp := NewTransactionParser(NewDbTypeFetcher(s.db))
	ts, errs := tp.Parse(string(bytes))

	for _, t := range ts {
		// TODO this call fails, i think when it tries to update/insert its associations (type most likely)
		s.db.Create(t)
	}

	fmt.Print(ts)
	fmt.Print(errs)

	// TODO Add response
}
