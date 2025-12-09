package examples

import (
	"io"
	"log"

	"github.com/yousef-muc/httpx"
)

func SimpleUsage() {
	// new instance
	client := httpx.New()

	res, err := client.Get("https://google.com", nil)
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
