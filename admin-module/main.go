package main

import (
	"encoding/csv"
	"fmt"
	"log"
    "os"
	// "time"
	"net/http"
)

func main() {

	//http.HandleFunc("/requestdata", requestData)
	// add users
	// bulkuserinsert
	// singleuserinsert
	// addusersmanager
	// bulkusersmanager
	// singleusermanager
	// manager data quota
	http.HandleFunc("/addmanagerstouser",addManagersToUsers)
	http.ListenAndServe("0.0.0.0:7575", nil)

}

func addManagersToUsers(res http.ResponseWriter,req *http.Request){

	file := req.FormValue("usermanagers")

	fd, err := os.Open(file)
	if err != nil {
		//log.Fatal(err)
		log.Println("Error opening file, false sent to frontend")
		fmt.Fprintf(res,"false")
		return
	}
	defer fd.Close()
	records, err := csv.NewReader(fd).ReadAll()
	if err != nil {
		//log.Fatal(err)
		fmt.Fprintf(res,"false")
		log.Println("Error parsing file, false sent to frontend")
	}
	m := make(map[string][]string)
	for _, v := range records {
		m[v[0]] = append(m[v[0]], v[1])
	}
	var userManagers [][]string

	for k, v := range m {
		userManagers = append(userManagers, append([]string{k}, v...))
		
	}

	go printRes(userManagers)
	//time.Sleep(2000 * time.Millisecond)

	fmt.Fprintf(res,"true")

}

func printRes(userManagers [][]string){

for _, v := range userManagers {
	fmt.Print("User ",v[0])
	fmt.Println(v[1:])
}
}