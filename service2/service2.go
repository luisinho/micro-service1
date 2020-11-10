package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
	Msg    string
	Codigo int
}

func main() {

	http.HandleFunc("/", processCoupon)
	http.ListenAndServe(":9091", nil)
}

func processCoupon(w http.ResponseWriter, r *http.Request) {

	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	email := r.PostFormValue("email")

	resultCoupon := makeHTTPCall("http://localhost:9092/", coupon, email)

	result := Result{Status: "declined", Codigo: resultCoupon.Codigo, Msg: ""}

	if resultCoupon.Codigo == 500 {
		result.Status = resultCoupon.Status
		result.Codigo = resultCoupon.Codigo
		result.Msg = ""
	}

	if ccNumber == "1" && resultCoupon.Codigo != 500 {
		result.Status = "approved"
		result.Codigo = resultCoupon.Codigo
		result.Msg = ""
	}

	if resultCoupon.Status == "invalid" && resultCoupon.Codigo != 500 {
		result.Status = "invalid coupon"
		result.Codigo = resultCoupon.Codigo
		result.Msg = ""
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHTTPCall(urlMicroService string, coupon string, email string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("email", email)

	res, err := http.PostForm(urlMicroService, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar", Codigo: http.StatusInternalServerError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}
	json.Unmarshal(data, &result)

	return result
}
