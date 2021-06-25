package main

import (
	"database/sql"
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

	r := mux.NewRouter()
	registerRoutes(r)

	http.Handle("/", r)
	log.Printf("serving API")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerRoutes(router *mux.Router) {
	router.Path("/test").Methods(http.MethodGet).HandlerFunc(testHandler)
	//server := policyServer{}
	//router.Path("/policy").Methods(http.MethodPost).HandlerFunc(server.postPolicyHandler)
	return
}

func testHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("ok"))
}

type Policy struct {
	Id      string `json:id`
	Content string `json:content`
}

type policyServer struct {
}

func (s *policyServer) postPolicyHandler(w http.ResponseWriter, r *http.Request) {

}
