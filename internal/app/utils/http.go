package utils

import (
	"log"
	"net/http"
)

func SafeCloseBody(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Println("error closing response body: ", err)
	}
}
