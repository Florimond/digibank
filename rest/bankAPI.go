package rest

import (
	"fmt"
	"net/http"

	"github.com/florhusq/digibank/account"
	"github.com/florhusq/digibank/event"
	"github.com/gorilla/mux"
)

type bankHandler struct {
	Manager *account.Manager
}

func (h *bankHandler) newAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
}

func (h *bankHandler) viewBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")

	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{error: no account id found}")
		return
	}

	h.Manager.ViewBalance(account)
}

func (h *bankHandler) newTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
}

func (h *bankHandler) viewTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf8")
}

func newBankHandler(db *event.Storage) *bankHandler {
	manager, err := account.NewManager(db)
	if err != nil {
		panic(err)
	}

	return &bankHandler{Manager: manager}
}

func ServeAPI(endpoint, metricsEndpoint string, db *event.Storage) error {
	r := mux.NewRouter()
	accountRouter := r.PathPrefix("/account").Subrouter()
	transferRouter := r.PathPrefix("/transfer").Subrouter()

	handler := newBankHandler(db)

	accountRouter.Methods("POST").Path("").HandlerFunc(handler.newAccountHandler)
	accountRouter.Methods("GET").Path("/{account}").HandlerFunc(handler.viewBalanceHandler)
	transferRouter.Methods("POST").Path("").HandlerFunc(handler.newTransactionHandler)
	transferRouter.Methods("GET").Path("/{account}").HandlerFunc(handler.viewTransactionHandler)

	/*eventsRouter.Methods("GET").Path("//{search}").HandlerFunc(handler.findEventHandler)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrChan := make(chan error)
	httpsErrChan := make(chan error)

	server := handlers.CORS()(r)
	go func() {
		fmt.Println("Starting the HTTPS server...")
		httpsErrChan <- http.ListenAndServeTLS(tlsendpoint, certPath+"cert.pem", certPath+"key.pem", server)
	}()

	go func() {
		fmt.Println("Starting the HTTP server...")
		httpErrChan <- http.ListenAndServe(endpoint, server)
	}()

	// Metrics are not as important as the main server. So no stopping it in case of error here. TODO: debate.
	go func() {
		fmt.Println("Starting prometheus...")
		h := mux.NewRouter()
		h.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(metricsEndpoint, h)
		fmt.Println("Error serving prometheus.", err)
	}()

	return httpErrChan, httpsErrChan*/

	return http.ListenAndServe(endpoint, r)
}
