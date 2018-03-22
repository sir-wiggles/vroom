package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var (
	host string
	port string
	db   *sql.DB

	pgUser = "foobar"
	pgPass = "cheese"
	pgHost = "localhost"
	pgPort = "5432"
	pgName = "foobar"
)

func init() {
	var (
		ok  bool
		err error
	)

	if host, ok = os.LookupEnv("HOST"); !ok {
		host = "127.0.0.1"
	}

	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}

	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgName,
	)

	log.Println("postgres: ", pgURL)
	if db, err = sql.Open("postgres", pgURL); err != nil {
		log.Fatal(err)
	}

}

func main() {
	r := mux.NewRouter()

	r.Use(mux.MiddlewareFunc(Middleware))

	s := r.Path("/user").Subrouter().StrictSlash(true)
	s.Methods("DELETE").Path("").HandlerFunc(UserDelete)
	s.Methods("GET").Path("").HandlerFunc(UserGet).Name("foobar")
	s.Methods("HEAD").Path("").HandlerFunc(UserHead)
	s.Methods("POST").Path("").HandlerFunc(UserPost)
	s.Methods("PUT").Path("").HandlerFunc(UserPut)

	server := http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	//schema := gojsonschema.NewReferenceLoader("file:///home/jeff/test.json")
	//document := gojsonschema.NewStringLoader(`{"name": "jeff", "billing_address": "1234 fake st"}`)

	//result, err := gojsonschema.Validate(schema, document)
	//if err != nil {
	//    panic(err.Error())
	//}

	//if result.Valid() {
	//    fmt.Printf("The document is valid\n")
	//} else {
	//    fmt.Printf("The document is not valid. see errors :\n")
	//    for _, desc := range result.Errors() {
	//        fmt.Printf("- %s\n", desc)
	//    }
	//}
	//os.Exit(0)

	log.Printf("service at %s:%s", host, port)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}

}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		log.Println("=========", route.GetName())
		next.ServeHTTP(w, r)
		w.Write([]byte("after \n"))
	})
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("From DELETE\n"))
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("From GET\n"))
}

func UserHead(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("From HEAD\n"))
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("From POST\n"))
}

func UserPut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("From PUT\n"))
}
