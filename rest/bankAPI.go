package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/florhusq/digibank/account"
	"github.com/florhusq/digibank/event"
	"github.com/gorilla/mux"
)

// bankHandler holds the manager and all the handlers
type bankHandler struct {
	Manager *account.Manager
}

// newAccountHandler handles requests of new account
func (h *bankHandler) newAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	openReq := struct {
		Customer string `json:"customer"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&openReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	account, err := h.Manager.Process(&account.OpenAccountCommand{Customer: openReq.Customer})

	resp := &struct {
		Account string `json:"account"`
	}{
		Account: account,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	//w.WriteHeader(http.StatusCreated)
}

// viewBalanceHandler handles request of the current balance for an account
func (h *bankHandler) viewBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{error: no account id found}")
		return
	}

	balance, err := h.Manager.ViewBalance(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	resp := &struct {
		Account string  `json:"account"`
		Balance float64 `json:"balance"`
	}{
		Account: account,
		Balance: balance,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

// newTransferHandler handles requests of new transder from a account to another
func (h *bankHandler) newTransferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	transacReq := struct {
		AccountFrom string  `json:"from"`
		AccountTo   string  `json:"to"`
		Amount      float64 `json:"amount"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&transacReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if _, err := h.Manager.Process(&account.TransferCommand{
		AccountFrom: transacReq.AccountFrom,
		AccountTo:   transacReq.AccountTo,
		Amount:      transacReq.Amount,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// newWithdrawHandler handles requests of withdrawal from an account
func (h *bankHandler) newWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	transacReq := struct {
		AccountFrom string  `json:"from"`
		Amount      float64 `json:"amount"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&transacReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if _, err := h.Manager.Process(&account.WithdrawCommand{
		AccountFrom: transacReq.AccountFrom,
		Amount:      transacReq.Amount,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// newDepositHandler handles requests of new transaction deposit to an account
func (h *bankHandler) newDepositHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	transacReq := struct {
		AccountTo string  `json:"to"`
		Amount    float64 `json:"amount"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&transacReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if _, err := h.Manager.Process(&account.DepositCommand{
		AccountTo: transacReq.AccountTo,
		Amount:    transacReq.Amount,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// viewTransactionHandler handles requests of viewing the whole history of transactions for an account
func (h *bankHandler) viewTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{error: no account id found}")
		return
	}

	transactions, err := h.Manager.ViewTransactions(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	resp := make([]struct {
		AccountFrom string  `json:"from"`   // The account which sends the money
		AccountTo   string  `json:"to"`     // The account which receives the money
		Amount      float64 `json:"amount"` // The amount
	}, len(transactions))

	// Mapping
	for i, t := range transactions {
		resp[i].AccountFrom = t.AccountFrom
		resp[i].AccountTo = t.AccountTo
		resp[i].Amount = t.Amount
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func newBankHandler(db *event.Storage) *bankHandler {
	manager, err := account.NewManager(db)
	if err != nil {
		panic(err)
	}

	return &bankHandler{Manager: manager}
}

// ServeAPI serves the API of the bank.
func ServeAPI(endpoint, metricsEndpoint string, db *event.Storage) error {
	r := mux.NewRouter()
	accountRouter := r.PathPrefix("/account").Subrouter()
	transferRouter := r.PathPrefix("/transfer").Subrouter()

	handler := newBankHandler(db)

	accountRouter.Methods("POST").Path("/").HandlerFunc(handler.newAccountHandler)
	accountRouter.Methods("GET").Path("/{account}/").HandlerFunc(handler.viewBalanceHandler)
	transferRouter.Methods("POST").Path("/transfer/").HandlerFunc(handler.newTransferHandler)
	transferRouter.Methods("POST").Path("/deposit/").HandlerFunc(handler.newDepositHandler)
	transferRouter.Methods("POST").Path("/withdraw/").HandlerFunc(handler.newWithdrawHandler)
	transferRouter.Methods("GET").Path("/{account}/").HandlerFunc(handler.viewTransactionHandler)

	return http.ListenAndServe(endpoint, r)
}
