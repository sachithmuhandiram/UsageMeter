package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"errors"

	"github.com/go-sql-driver/mysql"
	//"reflect"
)

type userObject struct {
	UserChain    string `json: "userchain"`
	UserEmail    string `json: "useremail"`
	DefaultQuota int    `json: "defaultquota"`
	IsManager    bool   `json: "ismanager"`
}

type managerDetails struct{
	ManagerEmail string `json: "manageremail"`
}

var userservice = os.Getenv("USERSERVICE")
var gatewayDB = os.Getenv("MYSQLDBGATEWAY")
var notificationservice = os.Getenv("NOTIFICATIONSERVICE")

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
	http.HandleFunc("/insertmanagerstouser", insertManagersToUser)
	http.HandleFunc("/insertuserdevices", insertUserDevices)
	http.ListenAndServe("0.0.0.0:7171", nil)
}

func requestData(res http.ResponseWriter, req *http.Request) {

	validUserInput, userIP := validateUserInput(req)

	/*
		This could be a problem as only one user connection active per time. Neet to test and verify.
	*/
	if validUserInput {

		validIP := checkIP(userIP)

		if validIP {

			userDetails, err := getUserDetails(userIP) // returns JSON object

			if err != nil {
				log.Println("There was a problem getting user data")
				return
			}
			// get remaining data quota
			eligibleToReq,err := checkRemainingQuota(userDetails.UserChain)

			if err != nil{
				log.Println("Error getting remaing data quota")
				return
			}

			if eligibleToReq == false{
				log.Println("User has sufficient data quota",userDetails.UserChain)
				fmt.Fprintf(res, "You have sufficient data quota for this month")
				return
			}
			// check whether he/she a manager

			if userDetails.IsManager {

				if userDetails.DefaultQuota == 0 {
					fmt.Fprintf(res, "You have no data quota restrictions")
					return
				}

				addedDataQuotatoManager := addQuotaToManager(userDetails.UserChain)

				if addedDataQuotatoManager {
					log.Println("Added data quota to Manager : ", userDetails.UserChain)
					// send an email to manager yet to develop
					return
				}
				log.Println("There is a problem adding data quota to manager : ", userDetails.UserChain)
				return

			}
			// check user has pending request pendingRequest table
			hasPendingReq := checkPendingRequest(userDetails.UserChain)

			if hasPendingReq {
				fmt.Fprintf(res, "You have a pending data quota request")
				return
			}
			// user details are ok, take manager emails
			managerEmails, err := getManagerEmails(userDetails.UserChain)

			if err != nil {
				log.Println("There is a problem with getting managers emails for User : ", userDetails.UserChain)
				return
			}
			// got manager emails

			adminEmails, err := getAdminEmails()

			if err != nil {
				log.Println("There is a problem getting admin emails")
				return
			}
			log.Println("Admin emaisl : ",adminEmails)
			requestedDataQuota := req.FormValue("data_amount")
			managers := strings.Join(managerEmails, ",")
			admins := strings.Join(adminEmails, ",")
			sentQuotaReq := dataQuotaRequest(userDetails.UserChain, requestedDataQuota, managers, admins)

			if sentQuotaReq != true {
				log.Println("There was a problem sending quota request email for user : ", userDetails.UserChain)
				return
			}
			log.Printf("Data quota request for %s email sent to managers", userDetails.UserChain)

			fmt.Fprintf(res, "Data requested sucessfully")

			//insert record to pendingRequest table | userchain | 1
			insertToPendingRequest(userDetails.UserChain)
		} else {
			log.Println("Wrong IP address. No user for this IP address  : ", userIP)
		}

	} else {
		fmt.Fprintf(res, "Wrong input!")
	}

}

func dataQuotaRequest(user string, quotaReq string, managers string, admins string) bool {

	// send email to admins and managers
	sendQuotaRequest := notificationservice + "/sendquotarequestmail"

	quotaRequestRes, err := http.PostForm(sendQuotaRequest, url.Values{"user": {user}, "requestedQuota": {quotaReq}, "managers": {managers}, "admins": {admins}})

	if err != nil {
		log.Println("Couldnt post data quota request")
		defer quotaRequestRes.Body.Close()
		return false
	}

	respBytes, err := ioutil.ReadAll(quotaRequestRes.Body)
	if err != nil {
		log.Println("Couldn't read quotaRequestRes body")
		return false
	}

	log.Println("quotaRequestRes", quotaRequestRes)
	respBool, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from userDetquotaRequestResailsRes body")
		return false
	}
	defer quotaRequestRes.Body.Close()
	return respBool

}

func validateUserInput(req *http.Request) (bool, string) {

	requestedDataQuota := req.FormValue("data_amount")

	validInt := regexp.MustCompile(`^[0-9]+`)

	if validInt.MatchString(requestedDataQuota) {

		userIPaddress := req.RemoteAddr

		userIPaddress, _, _ = net.SplitHostPort(userIPaddress)
		return true, userIPaddress

	} else {
		return false, ""
	}
}

func bulkUserInsert(res http.ResponseWriter, req *http.Request) {

	bulkDataFile := req.FormValue("bulk_data")

	db := dbConn()
	mysql.RegisterLocalFile(bulkDataFile)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + bulkDataFile + "' INTO TABLE users FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("User Bulk data inserted")
	defer db.Close()
}

func insertManagersToUser(res http.ResponseWriter, req *http.Request) {

	bulkDataFile := req.FormValue("user_managers")

	db := dbConn()
	mysql.RegisterLocalFile(bulkDataFile)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + bulkDataFile + "' INTO TABLE userManagers FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("User Managers data inserted")
	defer db.Close()
}

func insertUserDevices(res http.ResponseWriter, req *http.Request) {
	bulkDataFile := req.FormValue("user_devices")

	db := dbConn()
	mysql.RegisterLocalFile(bulkDataFile)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + bulkDataFile + "' INTO TABLE userDevices FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n' (userChain,deviceIP) SET isActive=1;")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("User Device data inserted")
	defer db.Close()
}

func checkIP(userIP string) bool {

	validUser := userservice + "/validuser"
	validUserRes, err := http.PostForm(validUser, url.Values{"userip": {userIP}})

	respBytes, err := ioutil.ReadAll(validUserRes.Body)
	if err != nil {
		log.Println("Couldn't read body")
	}

	respBool, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from body")
	}
	defer validUserRes.Body.Close()

	return respBool
}

func getUserDetails(userIP string) (userObject, error) {
	getUserDetails := userservice + "/userdetails"

	userDetailsRes, err := http.PostForm(getUserDetails, url.Values{"userip": {userIP}})

	if err != nil {
		log.Println("Problem getting User details")
		return userObject{}, err
	}

	userDetaildecoder := json.NewDecoder(userDetailsRes.Body)
	var userDetail userObject

	err = userDetaildecoder.Decode(&userDetail)
	if err != nil {
		log.Println("Could not decode user details")
		return userDetail, err
	}

	if userDetail.UserChain == "" {
		log.Println("Null data received from User service")
		return userDetail, errors.New("Data not received from User service")
	}
	defer userDetailsRes.Body.Close()
	return userDetail, nil
}

func getManagerEmails(userChain string) ([]string, error) {

	getManagereEmails := userservice + "/getmanageremails"

	managerEmailsRes, err := http.PostForm(getManagereEmails, url.Values{"userchain": {userChain}})

	managerEmailDecoder := json.NewDecoder(managerEmailsRes.Body)
	var managerEmail []string

	err = managerEmailDecoder.Decode(&managerEmail)
	if err != nil {
		log.Println("Could not decode Manager details")
		return managerEmail, err
	}

	defer managerEmailsRes.Body.Close()
	log.Println("Manager emails from main module : ",managerEmail)
	return managerEmail, nil
}

func addQuotaToManager(userChain string) bool {
	// this should directly call to adddataquota script

	managerDataQuota := getManagerDataQuota() // convert to int

	if managerDataQuota == "" {
		log.Println("Problem getting data quota to manager :", userChain)
		return false
	}
	// call to adddataquota script
	return true
}

func getManagerDataQuota() string {

	getManagerDataQuota := userservice + "/getmanagerdataquota"

	managerQuota, err := http.PostForm(getManagerDataQuota, url.Values{})

	if err != nil {
		log.Println("Error posting data to user service")
		return ""
	}
	bodyBytes, err := ioutil.ReadAll(managerQuota.Body)
	if err != nil {
		log.Fatal(err)
	}
	managerDataQuota := string(bodyBytes)

	return managerDataQuota
}

func getAdminEmails() ([]string, error) {
	getAdminEmails := userservice + "/getadminemails"

	adminEmailsRes, err := http.PostForm(getAdminEmails, url.Values{})

	adminEmailDecoder := json.NewDecoder(adminEmailsRes.Body)
	var adminEmail []string

	err = adminEmailDecoder.Decode(&adminEmail)
	if err != nil {
		log.Println("Could not decode admin email addresses")
		return adminEmail, err
	}

	defer adminEmailsRes.Body.Close()
	return adminEmail, nil
}

func checkRemainingQuota(user string)(bool,error){

	var remainingQuotaChecker := os.Getenv("QUOTACHECK")
	remainingQuota := userservice + "/checkquota"
	remainingQuotaRes, err := http.PostForm(validUser, url.Values{"user": {user},"method":{remainingQuotaChecker}})

	respBytes, err := ioutil.ReadAll(remainingQuotaRes.Body)
	if err != nil {
		log.Println("Couldn't read body of RemainingDataQuota")
	}

	respBool, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from RemainingDataQuota body")
		return false,errors.New("Could not read remaining data quota response")
	}
	defer remainingQuotaRes.Body.Close()

	return respBool,nil
}


func checkPendingRequest(user string) bool {
	db := dbConn()

	var hasPendingReq bool
	row := db.QueryRow("SELECT EXISTS(SELECT id FROM pendingRequest WHERE userChain=? AND isPending=1)", user)

	err := row.Scan(&hasPendingReq)
	if err != nil {
		log.Println("Error checking pendingRequest")
	}
	log.Println("check pending request : ", hasPendingReq)
	defer db.Close()
	return hasPendingReq
}

func insertToPendingRequest(user string) {
	// insert into pendingRequest table
	db := dbConn()

	insData, err := db.Prepare("INSERT INTO pendingRequest (userChain,isPending) VALUES(?,?)")
	if err != nil {
		log.Println("There is a problem inserting to pendingRequest table")
	}
	insData.Exec(user, 1)
	log.Println("Inserted record to pendingRequest table for user : ", user)
	defer db.Close()
}
