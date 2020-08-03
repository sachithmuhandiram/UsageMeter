package main

import (
	"fmt"
	"net/http"
	"regexp"
	"net"
	"os"
)

var userservice = os.Getenv("USERSERVICE")
func main() {

	http.HandleFunc("/requestdata", requestData)
	http.ListenAndServe("0.0.0.0:7171", nil)
}

func requestData(res http.ResponseWriter, req *http.Request) {

	validUserInput,userIP := validateUserInput(req)

	if validUserInput{
		//fmt.Fprintf(res, "Welcome to Home Page hahaha!")
		fmt.Println("User IP address : ",userIP)

		// send request to user service.

		parameters := userservice+"/checkuser?userip="+ userIP

		http.Redirect(res, req, parameters, http.StatusSeeOther)
		// response success
		// send success/failure to frontend 

	}else{
		fmt.Fprintf(res, "Wrong input!")
	}

}

func validateUserInput(req *http.Request)( bool,string){

	requestedDataQuota := req.FormValue("data_amount")

	validInt := regexp.MustCompile(`^[0-9]+`)

	if validInt.MatchString(requestedDataQuota){
		
		userIPaddress := req.RemoteAddr
		
		userIPaddress,_,_= net.SplitHostPort(userIPaddress)
		return true,userIPaddress

	}else{
		return false,""
	}
}