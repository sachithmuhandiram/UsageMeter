package main

import (

	"net/http"
	"log"

)

func main() {
	http.HandleFunc("/checkuser", checkUser)
	http.ListenAndServe(":7272", nil)
}

func checkUser(res http.ResponseWriter,req *http.Request){

	userIP := req.FormValue("userip")
	log.Println("IP address : ",userIP)

	// get user name /email

	// get managers / emails

	// get admins/emails

	// send emails to managers/admins

	// send response back to frontened

}