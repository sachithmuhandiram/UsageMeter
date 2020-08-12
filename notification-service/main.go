package main

import(
	"net/http"
	"os"
	//"net/smtp"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type emailDetails struct {
	From    string `json:"from"`
	Parse   string `json:"parse"`
}

func main(){

	http.HandleFunc("/sendquotarequestmail",sendQuotaRequestEmail)
	http.ListenAndServe("0.0.0.0:7474", nil)
}

func sendQuotaRequestEmail(res http.ResponseWriter, req *http.Request) {

	//get user
	// get managers
	//get admins
	user := req.FormValue("user")
	requestedQuota := req.FormValue("requestedQuota")
	managers := req.FormValue("managers")
	admins := req.FormValue("admins")

	log.Println("From notificateion service : ",user,requestedQuota,managers,admins)

	fmt.Fprintf(res,"true") 

	// body := msg + "\n" +  + "=" +  user.token + "\n This link valid only for 30 minutes"
	// from, pass := getCredintials()

	// emailMsg := "From: " + from + "\n" +
	// 	"To: " + user.email + "\n" +
	// 	"Subject: Register to the system\n\n" +
	// 	body

	// err := smtp.SendMail("smtp.gmail.com:587",
	// 	smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
	// 	from, []string{user.email}, []byte(emailMsg))

}

func getCredintials() (string, string) {
	jsonFile, err := os.Open("emailData.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var email emailDetails
	json.Unmarshal(byteValue, &email)
	//log.Println("Received email : " + email.From)

	return email.From, email.Parse

}