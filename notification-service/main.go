package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type emailDetails struct {
	From  string `json:"from"`
	Parse string `json:"parse"`
}

func main() {

	http.HandleFunc("/sendquotarequestmail", routeSample) //sendQuotaRequestEmail
	http.ListenAndServe("0.0.0.0:7474", nil)
}

func routeSample(res http.ResponseWriter, req *http.Request) {

	boolChan := make(chan bool)

	go blockingFunc(boolChan)

	if <-boolChan {
		fmt.Fprintf(res, "true")
		return
	}

	fmt.Fprintf(res, "false")
	return
}

func blockingFunc(ch chan bool) <-chan bool {
	time.Sleep(5000 * time.Millisecond)

	ch <- true

	return ch
}

func sendQuotaRequestEmail(res http.ResponseWriter, req *http.Request) {
	sendEmailChan := make(chan bool)

	go sendRequestEmail(req, sendEmailChan)

	sentEmail := <-sendEmailChan

	log.Print("sent email boolean channel value : ", sentEmail)
	if sentEmail {
		fmt.Fprintf(res, "true")
		close(sendEmailChan)
		return
	}

	fmt.Fprintf(res, "false")
	close(sendEmailChan)
	return
	//get user convert to first name
	// get managers
	//get admins

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

// This is used to create as a separate go route.
func sendRequestEmail(req *http.Request, ch chan bool) <-chan bool {
	user := req.FormValue("user")
	requestedQuota := req.FormValue("requestedQuota")
	managers := req.FormValue("managers")
	admins := req.FormValue("admins")

	receivers := managers + "," + admins
	toReceivers := strings.Split(receivers, ",")
	msg := "Additional " + requestedQuota + "GB data quota is requested by " + user + "."

	body := msg
	from, pass := getCredintials()

	emailMsg := "From: " + from + "\n" +
		"To: " + managers + "\n" +
		"Cc: " + admins + "\n" +
		"Subject: Request : additional dataquota \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, toReceivers, []byte(emailMsg))

	if err != nil {
		log.Println("Error : ", err)
		ch <- false
		return ch
	}
	ch <- true
	return ch
}
