#dashDB and Go

For issues that you encounter with this service, Go to [**Get help**](https://www.ibmdw.net/bluemix/get-help/) in the Bluemix development community.

You can bind a dashDB service instance to a Go runtime in Bluemix and then work with the data in the dashDB database.

###Required components

The following components are required to connect dashDB service from a Go application. They are all described in further detail in this topic.

- Go program
- db2cli Go module to connect to dashDB
- db2gobuildpack to install DB2 driver dependencies


###Go Program

#### Obtainig the dashDB connection details from VCAP_Services

In your Go code, you will need to import [go-cfenv](https://github.com/cloudfoundry-community/go-cfenv) package, which provides convenience functions and structures that map to Cloud Foundry environment variable primitives.  Using cfenv package we can parse the VCAP_SERVICES environment variable to retrieve the database connection information and connect to the database as shown in the following example.


```

    // Fetch dashDB service details
    appEnv, _ := cfenv.Current()
    services, _ := appEnv.Services.WithLabel("dashDB")
    dashDB := services[0]

    // Build connection string from the dashDB service detail fetched from VCAP_SERVICES
    connStr := []string{"DATABASE=", dashDB.Credentials["db"], ";", "HOSTNAME=", dashDB.Credentials["hostname"], ";", 
            "PORT=", dashDB.Credentials["port"], ";", "PROTOCOL=TCPIP", ";", "UID=", dashDB.Credentials["username"], ";", "PWD=", dashDB.Credentials["password"]};
    conn := strings.Join(connStr, "")
```
For more information on the structure of the VCAP_SERVICES environment variable see [Getting Started with dashDB Service](http://www.ng.bluemix.net/docs/#services/dashDB/index.html#dashDB)

#### Using the db2cli Go package to execute queries against dashDB

Now connect to dashDB using db2cli Go package, run your SQL query and process the resultset. Example code snippet is shown below

```

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

```
Please download and examine the attached sample for the rest of the code. If you use this sample, please edit the manifest.yml and specify a unique 'host' value.

###Uploading your application to bluemix

Use the -b option of the 'cf push' command to specify the [db2gobuildpack](https://github.com/ibmdb/db2gobuildpack) 

`cf push <app-name>` -b https://github.com/ibmdb/db2gobuildpack.git

###Related Links
- [Go db2cli package](https://bitbucket.org/phiggins/db2cli)

- [IBM DB2 v10.5 Knowledge Center](https://www-01.ibm.com/support/knowledgecenter/SSEPGG_10.5.0/com.ibm.db2.luw.kc.doc/welcome.html)
- [IBM DB2 developerWorks](http://www.ibm.com/developerworks/data/products/db2/)

