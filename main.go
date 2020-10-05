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
	"os/exec"
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

type managerDetails struct {
	ManagerEmail string `json: "manageremail"`
}

var userservice = os.Getenv("USERSERVICE")
var gatewayDB = os.Getenv("MYSQLDBGATEWAY")
var notificationservice = os.Getenv("NOTIFICATIONSERVICE")
var oldDB = os.Getenv("MYSQLOLDSYSTEM")

func dbConn(database string) (db *sql.DB) {
	db, err := sql.Open("mysql", database)

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

	log.Println("User request : ",req)

	/*
		This could be a problem as only one user connection active per time. Neet to test and verify.
	*/
	if validUserInput {

		validIP := validUserIP(userIP)

		log.Println("User IP address : ",userIP)

		if validIP {

			userDetails, err := getUserDetails(userIP) // returns JSON object

			if err != nil {
				log.Println("There was a problem getting user data",err)
				fmt.Fprintf(res,"There was a problem in our backend")
				return
			}
			// check whether he/she a manager

			if (userDetails.DefaultQuota == 0 ){

				fmt.Fprintf(res, "You have no data quota restrictions")
				return

			} // super manager without quota limitations

			eligibleToRequest,err := eligibleToRequstQuota(userDetails.UserChain,userDetails.IsManager)

			if err != nil {
				log.Println("Error getting user's remaing data quota",userDetails.UserChain)
				return
			}

			if eligibleToRequest{

				if (userDetails.IsManager){
					addDataQuotatoManager := addQuotaToManager(userDetails.UserChain)

				if addDataQuotatoManager {
					log.Println("Added data quota to Manager : ", userDetails.UserChain)
					fmt.Fprintf(res, "Data Quota added")
					// send an email to manager yet to develop
					return
				}
				log.Println("There is a problem adding data quota to manager : ", userDetails.UserChain)
				fmt.Fprintf(res, "Sorry, something not right in our side")
				return
				} // manager quota request

				hasPendingReq := checkPendingRequest(userDetails.UserChain)

				if hasPendingReq {
					fmt.Fprintf(res, "You have pending data quota requests")
					return
				}

				requestedDataQuota := req.FormValue("data_amount")
				requestAmount,_ := strconv.Atoi(requestedDataQuota)
				if requestAmount > 5 {
					fmt.Fprintf(res, "Maximum data quota request is 5GB")
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
				return
			// normal user request

			} // eligble to request quota if

			fmt.Fprintf(res, "You have sufficient data quota for this month")
			return
			
		} else {
			// Valide IP address else loop
			log.Println("Wrong IP address. No user for this IP address  : ", userIP)
			fmt.Fprintf(res, "Sorry, your device is not listed in our system")
		}

	} else {
		// Validate user input else
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

	//	userIPaddress := req.RemoteAddr
		userIPaddress := req.Header.Get("X-Forwarded-For")

		userIPaddress, _, _ = net.SplitHostPort(userIPaddress)
		return true, userIPaddress

	} else {
		return false, ""
	}
}

func bulkUserInsert(res http.ResponseWriter, req *http.Request) {

	bulkDataFile := req.FormValue("bulk_data")

	db := dbConn(gatewayDB)
	mysql.RegisterLocalFile(bulkDataFile)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + bulkDataFile + "' INTO TABLE users FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'")
	if err != nil {
		log.Println(err.Error())
		defer db.Close()
		return
	}
	log.Println("User Bulk data inserted")
	defer db.Close()
}

func insertManagersToUser(res http.ResponseWriter, req *http.Request) {

	bulkDataFile := req.FormValue("user_managers")

	db := dbConn(gatewayDB)
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

	db := dbConn(gatewayDB)
	mysql.RegisterLocalFile(bulkDataFile)
	_, err := db.Exec("LOAD DATA LOCAL INFILE '" + bulkDataFile + "' INTO TABLE userDevices FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n' (userChain,deviceIP) SET isActive=1;")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("User Device data inserted")
	defer db.Close()
}

func validUserIP(userIP string) bool {

	validUser := userservice + "/validuser"
	validUserRes, err := http.PostForm(validUser, url.Values{"userip": {userIP}})

	respBytes, err := ioutil.ReadAll(validUserRes.Body)
	defer validUserRes.Body.Close()

	if err != nil {
		log.Println("Couldn't read body")
		return false
	}

	respBool, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from body")
		return false
	}
	

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
	log.Println("Manager emails from main module : ", managerEmail)
	return managerEmail, nil
}

func addQuotaToManager(userChain string) bool {
	// this should directly call to adddataquota script

	managerDataQuota := getManagerDataQuota() // convert to int

	oldQuota := getCurrentQuota(userChain)

	if managerDataQuota == 0 {
		log.Println("Problem getting data quota to manager :", userChain)
		return false
	}
	// get default quota from both databases (sme-backend and usagemeter)
	// call to adddataquota script
	newQuota := managerDataQuota + oldQuota

	quotaAddedToManager := callAddQuotaScript(userChain,oldQuota,newQuota )

	if quotaAddedToManager {
		return true
	}

	return false
}

func getManagerDataQuota() int {

	getManagerDataQuota := userservice + "/getmanagerdataquota"

	managerQuota, err := http.PostForm(getManagerDataQuota, url.Values{})

	if err != nil {
		log.Println("Error posting data to user service")
		return 0
	}
	bodyBytes, err := ioutil.ReadAll(managerQuota.Body)
	if err != nil {
		log.Fatal(err)
	}
	managerDataQuota := string(bodyBytes)

	quota, _ := strconv.Atoi(managerDataQuota)
	return quota
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

// need to update this to accept normal users and managers
func eligibleToRequstQuota(user string,manager bool) (bool, error) {

	validToReq := userservice+ "/checkquota"
	remainingQuotaChecker := os.Getenv("QUOTACHECK")
	isManager := strconv.FormatBool(manager)
	// usertype 0 for normal users and 1 for managers
	remainingQuotaRes, err := http.PostForm(validToReq, url.Values{"user": {user}, "usertype":{isManager},"method": {remainingQuotaChecker}})

	respBytes, err := ioutil.ReadAll(remainingQuotaRes.Body)
	defer remainingQuotaRes.Body.Close()

	if err != nil {
		log.Println("Couldn't read body of RemainingDataQuota")
	}

	eligibleToReq, err := strconv.ParseBool(string(respBytes))
	if err != nil {
		log.Println("Couldn't parse bool from RemainingDataQuota body",err)
		return false, errors.New("Could not read remaining data quota response")
	}

	return eligibleToReq, nil
}

func checkPendingRequest(user string) bool {
	db := dbConn(gatewayDB)

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
	db := dbConn(gatewayDB)

	insData, err := db.Prepare("INSERT INTO pendingRequest (userChain,isPending) VALUES(?,?)")
	if err != nil {
		log.Println("There is a problem inserting to pendingRequest table")
	}
	insData.Exec(user, 1)
	log.Println("Inserted record to pendingRequest table for user : ", user)
	defer db.Close()
}

func callAddQuotaScript(userChain string, oldQuota int,newQuota int) bool {

	//updateDataQuota.sh chain oldQuota newQuota

	cmd, err := exec.Command("/bin/sh", "test.sh", userChain, strconv.Itoa(oldQuota),strconv.Itoa(newQuota)).Output()

	if err != nil {
		log.Println("There is a problem calling to Add quota script", err)
		return false
	}

	scriptRes := string(cmd)
	log.Println("Script returned", scriptRes)
	return true
}

func getCurrentQuota(userChain string)int{
	// read from sme-backend db
	db := dbConn(oldDB)

	var currentQutoa int
	row := db.QueryRow("SELECT data_quota FROM users WHERE ip_chain=? AND status=1", userChain)
	defer db.Close()
	err := row.Scan(&currentQutoa)
	if err != nil {
		log.Println("Error checking user's current data quota",userChain)
		log.Println("Error : ",err)
		return 0
	}

	return currentQutoa
}