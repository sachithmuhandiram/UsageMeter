package main

import (

	"net/http"
	"log"
	"fmt"
	"encoding/json"

)


type userObject struct{
	userChain string
	userEmail string
	defaultQuota int

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
		userChain : "sachithchain",
		userEmail : "sachith@vx.com",
		defaultQuota : 7000	}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(userDetail) 
}
	// get managers / emails

	// get admins/emails

	// send emails to managers/admins

	// send response back to frontened