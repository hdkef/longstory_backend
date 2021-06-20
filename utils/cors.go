package utils

import (
	"net/http"
)

func handleCors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*res).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Accept-Encoding, X-CSRF-Token, token")
}

func handleOptions(res *http.ResponseWriter, req *http.Request) bool {
	if req.Method == http.MethodOptions {
		(*res).WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func Cors(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		handleCors(&res)
		if handleOptions(&res, req) {
			return
		}
		next.ServeHTTP(res, req)
	}
}
