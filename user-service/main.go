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

func main() {
	http.HandleFunc("/validuser", validUser)
	http.HandleFunc("/userdetails", getUserDetails)
	http.HandleFunc("/getmanageremails", getManagerEmails)
	http.HandleFunc("/getmanagerdataquota", getManagerDataQuota)
	http.HandleFunc("/getadminemails", getAdminEmails)
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
		log.Println("Error getting user details", err)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(userDetail)
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(userDetail)
	defer db.Close()
}

func getManagerEmails(res http.ResponseWriter, req *http.Request) {

	userChain := req.FormValue("userchain")
	log.Println("User chain", userChain)
	//DB call to get manager emails
	managerEmails := []string{"sachith@vizuamatix.com", "sachithnalaka@gmail.com"}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(managerEmails)

}

func getAdminEmails(res http.ResponseWriter, req *http.Request) {
	adminEmails := []string{"msachithnalaka@yahoo.com", "janith@vx.com"}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(adminEmails)
}

// send emails to managers/admins

// send response back to frontened

func getManagerDataQuota(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "2000")
}
