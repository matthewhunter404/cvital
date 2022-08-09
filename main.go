package main

import (
	"cvital/router"
	"net/http"
)

func main() {

	http.ListenAndServe(":3000", router.NewRouter())
}
