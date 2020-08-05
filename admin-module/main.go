package main

import(
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
	http.ListenAndServe("0.0.0.0:7575", nil)
}