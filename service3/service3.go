package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

type Result struct {
	Status string
	Msg    string
	Codigo int
}

var coupons Coupons

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}

	return "invalid"
}

func checkCoupon(w http.ResponseWriter, r *http.Request) {

	coupon := r.PostFormValue("coupon")
	email := r.PostFormValue("email")

	resultEmail := makeHTTPCall("http://localhost:9093/", email)

	valid := coupons.Check(coupon)
	if resultEmail.Status == "invalid" {
		valid = resultEmail.Status
	}

	result := Result{Status: valid, Msg: "", Codigo: 0}
	if resultEmail.Codigo == 500 {
		result = Result{Status: resultEmail.Status, Msg: "", Codigo: resultEmail.Codigo}
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	log.Println(string(jsonResult))

	fmt.Fprintf(w, string(jsonResult))
}

func makeHTTPCall(urlMicroService string, email string) Result {

	values := url.Values{}
	values.Add("email", email)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroService, values)

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

func main() {

	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", checkCoupon)
	http.ListenAndServe(":9092", nil)
}
