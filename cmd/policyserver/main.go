package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/souza-bruno/policy-module/pkg/policydb"
	"github.com/souza-bruno/policy-module/pkg/resourceloader"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// hardcoded address and credentials for now
	db, err := sql.Open("mysql", "root:mypass@tcp(127.0.0.1:3306)/policy")
	if err != nil {
		log.Fatalf("failed to open db connection: %s", err)
	}
	defer db.Close()

	storage := policydb.NewStorage(db)

	err = resourceloader.ImportResourcesIntoDb(storage)
	if err != nil {
		log.Fatalf("failed to import resoruces into db: %s", err)
	}

	pServer := &policyServer{storage: &storage}

	r := mux.NewRouter()
	pServer.registerRoutes(r)

	http.Handle("/", r)
	log.Printf("serving API")
	log.Fatal(http.ListenAndServe(":8180", nil))
}

type policyServer struct {
	storage *policydb.Storage
}

func (s *policyServer) registerRoutes(router *mux.Router) {
	router.Path("/policy").Methods(http.MethodPost).HandlerFunc(s.postPolicyHandler)
	return
}

func (s *policyServer) postPolicyHandler(w http.ResponseWriter, r *http.Request) {
	policyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %s", err)
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	var policy policydb.Policy
	err = json.Unmarshal(policyJson, &policy)
	if err != nil {
		log.Printf("failed to unmarshall policy post body: %s", err)
		http.Error(w, "failed to unmarshall policy post body", http.StatusBadRequest)
		return
	}

	// TODO: should handle "already exists" case
	err = s.storage.CreatePolicy(policy)
	if err != nil {
		log.Printf("failed to create new policy: %s", err)
		http.Error(w, "no policies for you", http.StatusBadRequest)
		return
	}

	log.Printf("successfully created new policy %q", policy.Id)
	w.WriteHeader(http.StatusOK)
	return
}
