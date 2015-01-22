package main

import (
    "fmt"
    "net/http"
    "os"
    //"time"
    "database/sql"
    "strings"
    _ "bitbucket.org/phiggins/db2cli"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

func main() {

    http.HandleFunc("/", hello)
    http.HandleFunc("/connect", connect)
    fmt.Println("listening on port "+os.Getenv("PORT")+" ...")

    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func hello(res http.ResponseWriter, req *http.Request) {
	
    	fmt.Fprintln(res, "Hello World!!")
}

func connect(res http.ResponseWriter, req *http.Request) {
	
	fmt.Fprintln(res, "<html><head><title>Sample go dashDB application</title></head><body>")
	
	// Fetch dashDB service details
	appEnv, _ := cfenv.Current()
	services, _ := appEnv.Services.WithLabel("dashDB")
	
	dashDB := services[0]
	
	connStr := []string{"DATABASE=", dashDB.Credentials["db"], ";", "HOSTNAME=", dashDB.Credentials["hostname"], ";", 
			"PORT=", dashDB.Credentials["port"], ";", "PROTOCOL=TCPIP", ";", "UID=", dashDB.Credentials["username"], ";", "PWD=", dashDB.Credentials["password"]};
	conn := strings.Join(connStr, "")
	
	db, err := sql.Open("db2-cli", conn)
	if err != nil {

		fmt.Fprintln(res, "<h3>" )
		fmt.Fprintln(res, err )
		fmt.Fprintln(res, "</h3>" )
		return
	}
	
	defer db.Close()
	
	var (
		first_name string
		last_name string
	)
	
	stmt, err := db.Prepare("SELECT FIRST_NAME, LAST_NAME from GOSALESHR.employee FETCH FIRST 10 ROWS ONLY")
	if err != nil {
		fmt.Fprintln(res, "<h3>" )
		fmt.Fprintln(res, err )
		fmt.Fprintln(res, "</h3>" )
		return
	}
	defer stmt.Close()
	
	rows, err := stmt.Query()
    if err != nil {
        fmt.Fprintln(res, "<h3>" )
		fmt.Fprintln(res, err )
		fmt.Fprintln(res, "</h3>" )
		return
    }
    defer rows.Close()
	fmt.Fprintln(res, "<h3>Query: <br>SELECT FIRST_NAME, LAST_NAME from GOSALESHR.employee FETCH FIRST 10 ROWS ONLY</h3>" )
	fmt.Fprintln(res, "<h3>Result:<br></h3><table border=\"1\"><tr><th>First Name</th><th>Last Name</th></tr>" )
    for rows.Next() {
		err := rows.Scan(&first_name, &last_name)
		if err != nil {
			fmt.Fprintln(res, "<h3>" )
			fmt.Fprintln(res, err )
			fmt.Fprintln(res, "</h3>" )
			return;
			
		}
		fmt.Fprintln(res, "<tr><td>", first_name, "</td><td>",last_name, "</td></tr>")
	}
	fmt.Fprintln(res, "</table></body></html>")
}