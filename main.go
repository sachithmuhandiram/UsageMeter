package main

import (
	"fmt"
	"net/http"
	"regexp"
	"net"
	"os"
	"log"
	"net/url"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

type userObject struct{
	userChain string `json: "userchain"`
	userEmail string `json: "useremail"`
	defaultQuota int `json: "defaultquota"`

}


var userservice = os.Getenv("USERSERVICE")
var gatewayDB	= os.Getenv("MYSQLDBGATEWAY")

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", gatewayDB)

	if err != nil {
		log.Println("Cant open database connection")
	}
	return db
}


func main() {

	http.HandleFunc("/requestdata", requestData)
	http.HandleFunc("/bulkuserinsert", bulkUserInsert)
	http.ListenAndServe("0.0.0.0:7171", nil)
}

func requestData(res http.ResponseWriter, req *http.Request) {

	validUserInput,userIP := validateUserInput(req)

	/*
		This could be a problem as only one user connection active per time. Neet to test and verify.
	*/
	if (validUserInput){

		validIP := checkIP(userIP)
	
		if (validIP){

			userDetails := getUserDetails(userIP) // returns JSON object
			log.Println("User details : ",userDetails)

		}else{
			log.Println("Wrong IP address. No user for this IP address")
		}
		
	
	

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

func bulkUserInsert(res http.ResponseWriter,req *http.Request){

	bulkDataFile := req.FormValue("bulk_data")

	db := dbConn()
	mysql.RegisterLocalFile(bulkDataFile)
	_,err := db.Exec("LOAD DATA LOCAL INFILE '"+bulkDataFile+"' INTO TABLE users FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'")
        if err != nil {
            log.Println(err.Error())
        }
	log.Println("User Bulk data inserted")
	defer db.Close()
}

func checkIP(userIP string) bool{

	validUser := userservice +"/validuser"
	validUserRes, err := http.PostForm(validUser, url.Values{"userip": {userIP}})

	defer validUserRes.Body.Close()

	respBytes, err := ioutil.ReadAll(validUserRes.Body)
	if err != nil {
		log.Println("Couldn't read body")
	}

	respBool, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from body")
	}

	return respBool
}

func getUserDetails(userIP string) userObject{
	getUserDetails := userservice +"/userdetails"

	userDetailsRes, err := http.PostForm(getUserDetails, url.Values{"userip": {userIP}})

	log.Println("User details response from user details :",userDetailsRes)

	userDetaildecoder := json.NewDecoder(userDetailsRes.Body)
    var userDetail userObject
    err = userDetaildecoder.Decode(&userDetail)
    if err != nil {
		log.Println("Could not decode user details")
	}

	log.Println("User details from main : ",userDetail)
	return userDetail
}