package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Status     string
	msg        string
	codigoErro int
}

func validaEmail(w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("email")

	result := Result{Status: "invalid", codigoErro: 0}

	if email == "teste@teste.com" {
		result.Status = "valid"
		result.codigoErro = 0
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}

func main() {

	http.HandleFunc("/", validaEmail)
	http.ListenAndServe(":9093", nil)
}
