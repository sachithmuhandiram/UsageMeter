package main

import(
	"net/http"
	"os"
	"net/smtp"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strings"

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

	//get user convert to first name
	// get managers
	//get admins
	user := req.FormValue("user")
	requestedQuota := req.FormValue("requestedQuota")
	managers := req.FormValue("managers")
	//admins := req.FormValue("admins")

	toManagers := strings.Split(managers, ",") 
	msg := "Additional " + requestedQuota + "GB data quota is requested by "+ user +"."
	log.Println("From notificateion service : ",managers)
	log.Println("To managers : ",toManagers)

	

	body := msg
	from, pass := getCredintials()

	emailMsg := "From: " + from + "\n" +
	"To: " + managers + "\n" +
	"Subject: Request : additional dataquota \n\n" +
	body

	err := smtp.SendMail("smtp.gmail.com:587",
	smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
	from, toManagers, []byte(emailMsg))

	if err != nil{
		fmt.Fprintf(res,"false") 
		log.Println("Error : ",err)
	}
	fmt.Fprintf(res,"true") 
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