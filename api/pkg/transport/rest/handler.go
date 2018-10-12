package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//Handler manages http rest methods handling
func Handler(l listing.Service) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.Headers("Content-Type", "application/json")
	s.HandleFunc("/account/{hash}", getAccount(l)).Queries("signature", "{signature}")
	return s
}

func getAccount(l listing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := []byte(vars["hash"])
		sig := []byte(vars["signature"])
		acc, err := l.GetAccount(hash, sig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(acc)
	}
}
