package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_"github.com/godror/godror"

	"zero-to-prod.norrbom.org/models"
)

func validate(r *http.Request, params []string) error {
	for _, s := range params {
		if len(r.URL.Query()[s]) < 1 || r.URL.Query()[s][0] == "" {
			return errors.New(s + " parmater expected!")
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "My Service")
	w.Header().Set("Content-Type", "application/json")

	params := []string{"brand", "locale", "jurisdiction"}
	var response string

	err := validate(r, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
	} else {
		// DO SOMETHING!
		response = fmt.Sprintf(`{
    "message": "%s"
}
`, "hello!")
	}
	fmt.Fprintf(w, response)
}

// health endpoint
func hanlderHealth(w http.ResponseWriter, r *http.Request) {
    err := models.Ping()
    if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }
	fmt.Fprintf(w, "OK")
}

// k8s readiness and liveness probes
func handlerProbe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	
	var err error

	log.SetOutput(os.Stdout)
	log.Println("Server is starting...")
	var connectStr = fmt.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		os.Getenv("MY_DB_ORACLE_USERNAME"),
		os.Getenv("MY_DB_ORACLE_PASSWORD"),
		os.Getenv("MY_DB_ORACLE_CONNECTION_URL"),
	)

	models.Db, err = sql.Open("godror", connectStr)
	if err != nil {
		log.Fatal("error connecting to Oracle database: ", err)
	}
	defer models.Db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", hanlderHealth)
	http.HandleFunc("/probe", handlerProbe)
	http.Handle("/prometheus", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
	log.Fatal(http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
