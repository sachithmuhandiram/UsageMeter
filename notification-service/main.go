package main

import(
	"net/http"
)

func main(){

	http.ListenAndServe("0.0.0.0:7474", nil)
}