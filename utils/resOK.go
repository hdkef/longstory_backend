package utils

import "net/http"

func ResOK(res *http.ResponseWriter, message string) {
	(*res).WriteHeader(http.StatusOK)
	(*res).Write([]byte(message))
}
