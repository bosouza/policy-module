package policyserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	registerRoutes(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerRoutes(router *mux.Router) {
	router.Path("/test").Methods(http.MethodGet).HandlerFunc(testHandler)
	return
}

func testHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("ok"))
}
