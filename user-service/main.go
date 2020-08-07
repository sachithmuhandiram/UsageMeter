package main

import (

	"net/http"
	"log"
	"fmt"

)

func main() {
	http.HandleFunc("/validuser", validUser)
	http.ListenAndServe(":7272", nil)
}

// This returns true/false to main service
func validUser(res http.ResponseWriter,req *http.Request){

	userIP := req.FormValue("userip")
	log.Println("IP address : ",userIP)

	fmt.Fprintf(res,"huta") 

}

	// get user name /email

	// get managers / emails

	// get admins/emails

	// send emails to managers/admins

	// send response back to frontened