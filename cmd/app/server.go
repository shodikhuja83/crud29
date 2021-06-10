package app

import (
	"encoding/json"
	"errors"
	"github.com/shodikhuja83/crud/pkg/customers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
}

//NewServer construct
func NewServer(mux *mux.Router, customerSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	log.Println("Init method")
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerById).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
}

func (s *Server) handleGetCustomerById(writer http.ResponseWriter, request *http.Request) {
	idParam := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ByID(request.Context(), id)
	log.Println(item)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, item)
}

func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.All(r.Context())

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, items)
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.AllActive(r.Context())

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, items)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.Delete(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	customer, err := s.customerSvc.Save(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, customer)
}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

func jsonResponse(writer http.ResponseWriter, data interface{}) {
	item, err := json.Marshal(data)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(item)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}
