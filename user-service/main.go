package main

import (

	"net/http"

)
func main() {

	http.HandleFunc("/checkuser", checkUser)
	http.ListenAndServe("0.0.0.0:7272", nil)
}

func checkUser(res http.ResponseWriter,req *http.Request){

	// get user name /email

	// get managers / emails

	// get admins/emails

	// send emails to managers/admins

	// send response back to frontened

}