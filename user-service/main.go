package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"errors"
	"encoding/csv"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var gatewayDB = os.Getenv("MYSQLDBGATEWAY")

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", gatewayDB)

	if err != nil {
		log.Println("Cant open database connection")
	}
	return db
}

type userObject struct {
	UserChain    string
	UserEmail    string
	DefaultQuota int
	IsManager    bool
}

type managerDetails struct{
	ManagerEmail string
}
func main() {
	http.HandleFunc("/validuser", validUser)
	http.HandleFunc("/userdetails", getUserDetails)
	http.HandleFunc("/getmanageremails", getManagerEmails)
//	http.HandleFunc("/getmanagerdataquota", getManagerDataQuota)
	http.HandleFunc("/getadminemails", getAdminEmails)
	http.HandleFunc("/checkquota",checkQuota)
	http.ListenAndServe(":7272", nil)
}

// This returns true/false to main service
func validUser(res http.ResponseWriter, req *http.Request) {

	//userIP := req.FormValue("userip")
	fmt.Fprintf(res, "true")

}

// get user name /email
func getUserDetails(res http.ResponseWriter, req *http.Request) {

	userIP := req.FormValue("userip")

	// get user details for the IP
	db := dbConn()

	var userDetail userObject
	row := db.QueryRow("CALL GetUserDetails(?)", userIP)

	//scan sequence should be equal to user details return from SP
	err := row.Scan(&userDetail.UserChain, &userDetail.UserEmail, &userDetail.IsManager, &userDetail.DefaultQuota)
	if err != nil {
		db.Close()
		log.Println("Error getting user details", err)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(userDetail)
	}

	defer db.Close()
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(userDetail)
	
}

func getManagerEmails(res http.ResponseWriter, req *http.Request) {
//GetManagersEmail
	userChain := req.FormValue("userchain")
	db := dbConn()

	rows,err := db.Query("CALL GetManagersEmail(?)", userChain)

	if err != nil {
		fmt.Println("Failed to run query", err)
		db.Close()
        return
    }
	managerEmails:=[]string{}
    for rows.Next() {
        var managerMail string
        rows.Scan(&managerMail)
        managerEmails = append(managerEmails, managerMail)
    }

    defer db.Close()
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(managerEmails)
	
}

func getAdminEmails(res http.ResponseWriter, req *http.Request) {
	
	db := dbConn()
	rows,err := db.Query("SELECT email FROM users where isAdmin=?", 1)

	if err != nil {
        fmt.Println("Failed to run query", err)
        return
    }
	var adminEmails  []string

    for rows.Next() {
        var adminEmail string
        rows.Scan(&adminEmail)
        adminEmails = append(adminEmails, adminEmail)
    }

    defer db.Close()
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(adminEmails)
}

// send emails to managers/admins

// send response back to frontened

func checkQuota(res http.ResponseWriter, req *http.Request){

	user := req.FormValue("user")
	isManager := req.FormValue("usertype") // true /false is returned as string
	method := req.FormValue("method")

	if (method == "db"){	// possible limitation, this should come from db table
		log.Println("Call to database service to check user remaining data quota",user)
		return
	}
	// method is file
//	var remainingQuotaChecker = os.Getenv("QUOTACHECK")
	currentTime := time.Now()
  
  	month := currentTime.Format("2006-01")
	// read data from file
	availableQuota,err := availableDataQuota(user,month)

	if err != nil{
		log.Println("Couldnt find record for user : ",user)
	}
	// get user's min data quota
	if isManager == "true"{
		//get managers min quota
		minimumQuotaLimit := getManagerDataQuotaLimit()
		if (availableQuota > minimumQuotaLimit){
			fmt.Fprintf(res, "false")
			return
		}
		
		fmt.Fprintf(res, "true")
		return
	}
	// for normal user
	minimumQuotaLimit := getUserDataQuotaLimit()
	log.Println("Available and minimum quotas : ",availableQuota,minimumQuotaLimit)
	if (availableQuota > minimumQuotaLimit){
		fmt.Fprintf(res, "false")
		return
	}
	
	fmt.Fprintf(res, "true")
	return
}

// readDataFile
func availableDataQuota(user string,month string) (int,error){

	log.Println("User to read from file : ",user)
	log.Println("file to read ",month)

	//quotaFileLocation := os.Getenv("QUOTAFILE")

	monthQuotaFile := "/files/"+month+".txt" //quotaFileLocation+"/"+

	userQuota,currentUsage := getQuotaAndUsage(user,monthQuotaFile)

	if ((userQuota == 0) && (currentUsage == 0)){
		return 0, errors.New("No user record found for this user")
	}

	currentUsage = (currentUsage/(1024*1024))
	userQuota = (userQuota/(1024*1024))
log.Println("User's current usage ",currentUsage)
log.Println("User's quota ",userQuota)
	return (userQuota-currentUsage),nil // this is in bytes
	// 
}

// return user's current usage and userQuota
func getQuotaAndUsage(user string, file string)(int,int){

	usageFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer usageFile.Close()

	csvLines, err := csv.NewReader(usageFile).ReadAll()
    if err != nil {
        fmt.Println(err)
    }    
    for _, line := range csvLines {
        if (line[0]==user){
			quota,_ := strconv.Atoi(line[2])
			usage,_ := strconv.Atoi(line[3])
			return quota,usage
		}
	}
	
	return 0,0
}


func getManagerDataQuotaLimit() int {

	return 2000
	// This is only for testing, should get from DB
}

func getUserDataQuotaLimit()int{
	return 1000
}