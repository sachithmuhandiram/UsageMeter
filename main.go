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
)

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
		validUser := userservice +"/validuser"

    	validUserRes, err := http.PostForm(validUser, url.Values{"userip": {userIP}})

	if err != nil {
		log.Println("Error from user service")
	}
		// response success
		// send success/failure to frontend 
	log.Println("Response from User service",validUserRes)

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