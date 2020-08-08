package main

import (

	"net/http"
	"log"
	"fmt"
	"encoding/json"

)


type userObject struct{
	UserChain string
	UserEmail string
	DefaultQuota int

}

func main() {
	http.HandleFunc("/validuser", validUser)
	http.HandleFunc("/userdetails",userDetails)
	http.ListenAndServe(":7272", nil)
}

// This returns true/false to main service
func validUser(res http.ResponseWriter,req *http.Request){

	userIP := req.FormValue("userip")
	log.Println("IP address : ",userIP)

	fmt.Fprintf(res,"true") 

}
	// get user name /email
func userDetails(res http.ResponseWriter,req *http.Request){

	userIP := req.FormValue("userip")
	log.Println("User IP address : ",userIP)

	userDetail := userObject{
		UserChain : "sachithchain",
		UserEmail : "sachith@vx.com",
		DefaultQuota : 7000	}

		log.Println("User details from  user service",userDetail)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(userDetail) 
}
	// get managers / emails

	// get admins/emails

	// send emails to managers/admins

	// send response back to frontened