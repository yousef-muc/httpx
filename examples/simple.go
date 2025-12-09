package examples

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/yousef-muc/httpx"
)

func simpleGetRequest() {
	// new instance
	client := httpx.New()

	defaultHeaders := make(http.Header)
	defaultHeaders.Set("Authorization", "Bearer ABC-123")

	client.SetHeaders(defaultHeaders)

	requestHeaders := make(http.Header)
	requestHeaders.Set("Content-Type", "application/json")

	res, err := client.Get("https://dummyjson.com/carts", requestHeaders)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	log.Println(string(bytes))
}

func simplePostRequest() {

	type User struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	user := &User{
		Firstname: "Yousef",
		Lastname:  "Hejazi",
	}

	requestBody, _ := json.Marshal(&user)

	// new instance
	client := httpx.New()

	defaultHeaders := make(http.Header)
	defaultHeaders.Set("Authorization", "Bearer ABC-123")

	client.SetHeaders(defaultHeaders)

	requestHeaders := make(http.Header)
	requestHeaders.Set("Content-Type", "application/json")

	res, err := client.Post("https://dummyjson.com/carts/add", requestHeaders, requestBody)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	log.Println(string(bytes))
}
