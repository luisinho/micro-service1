package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type Result struct {
	Status     string
	msg        string
	codigoErro int
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {

	result := makeHttpcall("http://localhost:9091/", r.FormValue("coupon"), r.FormValue("cc-number"), r.FormValue("email"))

	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, result)
}

func makeHttpcall(urlMicroService string, coupon string, ccNumber string, email string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)
	values.Add("email", email)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroService, values)

	if err != nil {
		result := Result{Status: "Servidor fora do ar", codigoErro: 0, msg: ""}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	log.Println(data)

	result := Result{}
	json.Unmarshal(data, &result)

	// log.Println(result)

	return result
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}
