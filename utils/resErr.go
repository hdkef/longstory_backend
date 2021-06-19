package utils

import (
	"net/http"
)

//Respond with error
func ResErr(res *http.ResponseWriter, code int, err error) {
	(*res).WriteHeader(code)
	(*res).Write([]byte(err.Error()))
}
