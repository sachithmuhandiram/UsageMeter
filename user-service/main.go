package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	http.HandleFunc("/getmanagerdataquota", getManagerDataQuota)
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

	userIP := "192.168.10.16" //req.FormValue("userip")

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

func getManagerDataQuota(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "2000")
}

func checkQuota(res http.ResponseWriter, req *http.Request){

	user := req.FormValue("user")
	method := req.FormValue("method")

	if (method == "db"){
		log.Println("Call to database service to check user remaining data quota",user)
		return
	}
	// read data from file

}