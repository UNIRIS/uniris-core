package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/uniris/uniris-core/api/pkg/crypto"
	"github.com/uniris/uniris-core/api/pkg/listing"
	"github.com/uniris/uniris-core/api/pkg/mock"
	"github.com/uniris/uniris-core/api/pkg/transport/rest"
)

func main() {

	port := flag.Int("port", 8080, "API port")
	flag.Parse()

	r := mux.NewRouter()

	loadSwagger(r)
	loadAPI(r)

	log.Printf("Server running on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

func loadSwagger(r *mux.Router) {
	staticDir, _ := filepath.Abs("../../web/static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootPage, _ := filepath.Abs("../../web/index.html")
		http.ServeFile(w, r, rootPage)
	})

	r.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		swaggerFile, _ := filepath.Abs("../../api/swagger-spec/swagger.yaml")
		http.ServeFile(w, r, swaggerFile)
	})
}

func loadAPI(r *mux.Router) {
	sharedBioPrivKey := []byte("")
	client := mock.NewClient()
	validator := new(crypto.RequestValidator)
	lister := listing.NewService(sharedBioPrivKey, client, validator)

	apiR := rest.Handler(lister)
	r.Handle("/api", apiR)
}
