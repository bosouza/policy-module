package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/souza-bruno/policy-module/pkg/evaluator"
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

	eval, err := evaluator.NewPolicyEvaluator(storage)
	if err != nil {
		log.Fatalf("failed to create evaluator: %s", err)
	}

	pServer := &policyServer{storage: storage, eval: eval}

	r := mux.NewRouter()
	pServer.registerRoutes(r)

	http.Handle("/", r)
	log.Printf("serving API")
	log.Fatal(http.ListenAndServe(":8180", nil))
}

type policyServer struct {
	storage *policydb.Storage
	eval    *evaluator.PolicyEvaluator
}

func (s *policyServer) registerRoutes(router *mux.Router) {
	router.Path("/policy").Methods(http.MethodPost).HandlerFunc(s.postPolicyHandler)
	router.Path("/assign/{userId}/{policyId}").Methods(http.MethodPut).HandlerFunc(s.assignPolicy)
	router.Path("/evaluate").Methods(http.MethodPost).HandlerFunc(s.evaluate)
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

func (s *policyServer) assignPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	policyId := vars["policyId"]
	userId := vars["userId"]

	err := s.storage.AssignPolicyToUser(policyId, userId)
	if err != nil {
		log.Printf("failed to assign policy %q to user %q: %s", policyId, userId, err)
		http.Error(w, "failed assignment", http.StatusBadRequest)
		return
	}

	err = s.eval.RefreshData()
	if err != nil {
		log.Fatalf("failed to refresh data: %s", err)
	}

	log.Printf("successfully assigned policy %q to user %q", policyId, userId)
	w.WriteHeader(http.StatusOK)
	return
}

func (s *policyServer) evaluate(w http.ResponseWriter, r *http.Request) {
	evaluationJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %s", err)
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	var input EvaluateInput
	err = json.Unmarshal(evaluationJson, &input)
	if err != nil {
		log.Printf("failed to unmarshall evaluate post body: %s", err)
		http.Error(w, "failed to unmarshall evaluate post body", http.StatusBadRequest)
		return
	}

	result, err := s.eval.Evaluate(r.Context(), input.PolicyCheckId, input.Input)
	if err != nil {
		log.Printf("failed to evaluate: %s", err)
		http.Error(w, "failed to evaluate", http.StatusBadRequest)
		return
	}

	mResult, err := json.Marshal(result)
	if err != nil {
		log.Printf("failed to marshal evaluation result: %s", err)
		http.Error(w, "failed to marshal evaluation result", http.StatusInternalServerError)
		return
	}

	log.Printf("successfully evaluated policycheck %q", input.PolicyCheckId)
	w.Write(mResult)
}

type EvaluateInput struct {
	PolicyCheckId string      `json:"policyCheckId"`
	Input         interface{} `json:"input"`
}
