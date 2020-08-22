package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type userObject struct {
	UserChain    string
	UserEmail    string
	DefaultQuota int
	IsManager    bool
}

func main() {
	http.HandleFunc("/validuser", validUser)
	http.HandleFunc("/userdetails", userDetails)
	http.HandleFunc("/getmanageremails", getManagerEmails)
	http.HandleFunc("/getmanagerdataquota", getManagerDataQuota)
	http.HandleFunc("/getadminemails", getAdminEmails)
	http.ListenAndServe(":7272", nil)
}

// This returns true/false to main service
func validUser(res http.ResponseWriter, req *http.Request) {

	//userIP := req.FormValue("userip")
	fmt.Fprintf(res, "true")

}

// get user name /email
func userDetails(res http.ResponseWriter, req *http.Request) {

	userIP := req.FormValue("userip")
	log.Println("User IP address : ", userIP)

	// get user details for the IP

	userDetail := userObject{
		UserChain:    "nalakachain",
		UserEmail:    "nalaka@vx.com",
		DefaultQuota: 0,
		IsManager:    true}

	log.Println("User details from  user service", userDetail)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(userDetail)
}

func getManagerEmails(res http.ResponseWriter, req *http.Request) {

	userChain := req.FormValue("userchain")
	log.Println("User chain", userChain)
	//DB call to get manager emails
	managerEmails := []string{"sachith@vizuamatix.com", "sachithnalaka@gmail.com"}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(managerEmails)

}

func getAdminEmails(res http.ResponseWriter, req *http.Request) {
	adminEmails := []string{"msachithnalaka@yahoo.com", "janith@vx.com"}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(adminEmails)
}

// send emails to managers/admins

// send response back to frontened

func getManagerDataQuota(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "2000")
}
