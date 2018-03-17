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
	http    *http.Client
	r       *mux.Router
	db      *gorm.DB
	handler *Handler
}

func NewServer(r *mux.Router, db *gorm.DB) (*Server) {

	handler := &Handler{
		inv:           NewInventory(),
		ledger:        NewLedger(&EveMarketerAPI{client: &http.Client{}}),
		typeFetcher:   NewDbTypeFetcher(db),
		playerFetcher: NewDbPlayerFetcher(db),
	}

	s := &Server{
		http:    &http.Client{},
		r:       r,
		db:      db,
		handler: handler,
	}

	r.HandleFunc("/", s.GetRoot).Methods(http.MethodGet)
	r.HandleFunc("/parse", s.parseLog).Methods(http.MethodPost)
	r.HandleFunc("/transactions", s.GetTransactions).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{transactionId}/markForCorp", s.MarkTransactionForCorp).Methods(http.MethodPost)
	r.HandleFunc("/inventory", s.GetInventory).Methods(http.MethodGet)
	r.HandleFunc("/ledger", s.GetLedger).Methods(http.MethodGet)
	r.HandleFunc("/commits", s.Commit).Methods(http.MethodPost)
	r.HandleFunc("/commits/{commitId}/rollback", s.RollbackCommit).Methods(http.MethodPost)
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

/**
Fetches a list of uncommitted transactions
 */
func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var ts []Transaction
	s.db.Order("creation_date DESC").Find(&ts)
	// TODO Filter on uncomitted

	bytes, err := json.Marshal(ts)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Could not convert transactions to JSON: %s", err)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

/**
Marking transactions for corp changes the playerName of the transaction at runtime to ADHC
 */
func (s *Server) MarkTransactionForCorp(w http.ResponseWriter, r *http.Request) {

	// TODO We need a way to unmark transactions as well
	// TODO We need to make sure only uncommitted transactions can be changed

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
	s.db.Order("creation_date ASC").Find(&ts)
	// TODO Filter on uncomitted

	creditMuts, debitMuts, err := s.handler.Process(ts)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
	}

	var inv []GetInventoryRspItem

	for typeId, contents := range s.handler.inv.contents {
		ty, _ := s.handler.typeFetcher.getTypeById(typeId)

		inv = append(inv, GetInventoryRspItem{
			TypeId:   typeId,
			TypeName: ty.TypeName,
			Stacks:   contents,
		})
	}

	body, _ := json.Marshal(&GetInventoryRsp{
		Inventory: inv,
		Mutations: append(creditMuts, debitMuts...),
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type GetInventoryRsp struct {
	Inventory []GetInventoryRspItem
	Mutations []InventoryMutation
}

type GetInventoryRspItem struct {
	TypeId   int
	TypeName string
	Stacks   []InventoryStack
}

func (s *Server) GetLedger(w http.ResponseWriter, r *http.Request) {
	var ts []Transaction
	s.db.Order("creation_date ASC").Find(&ts)
	// TODO Filter on uncomitted

	creditMuts, debitMuts, err := s.handler.Process(ts)

	mutations, err := s.handler.ledger.HandleMutations(debitMuts, creditMuts)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
	}

	ledgerSummary := s.handler.ledger.CalculateLedgerSummary(mutations)

	body, _ := json.Marshal(&GetLedgerRsp{
		Ledger:    ledgerSummary,
		Mutations: mutations,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type GetLedgerRspItem struct {
	PlayerName string
	Amount     float64
}

type GetLedgerRsp struct {
	Ledger    []GetLedgerRspItem
	Mutations []LedgerMutation
}

func (s *Server) RollbackCommit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) Commit(w http.ResponseWriter, r *http.Request) {

	// TODO There is quite some duplication between this, GetLedger and GetInventory
	/**

So when the payout happens we want to
- Fetch all transactions that have not been processed
- Calculate and store new inventory state
- Calculate per person a bill
- Store that bill for future reference
- Mark transactions as processed

	So in short we want to do a GetLedger, then store the results under a Commit with a CommitId so we can rollback if needed.

	 */
	var ts []Transaction
	s.db.Where("processed_in_commit IS null").Find(&ts)

	commit := Commit{}
	s.db.Create(&commit)

	w.WriteHeader(http.StatusNotImplemented)
}

type Commit struct {
	gorm.Model
}

func (s *Server) parseLog(w http.ResponseWriter, r *http.Request) {

	bytes, _ := ioutil.ReadAll(r.Body)
	tp := NewTransactionParser(NewDbTypeFetcher(s.db))
	ts, errs := tp.Parse(string(bytes))

	for _, t := range ts {
		s.db.Create(t)
	}

	w.WriteHeader(http.StatusNoContent)
	body, _ := json.Marshal(&ParseLogRsp{
		Errors:       errs,
		Transactions: ts,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type ParseLogRsp struct {
	Transactions []Transaction
	Errors       []error
}
