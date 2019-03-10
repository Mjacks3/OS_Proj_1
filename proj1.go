package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "github.com/gorilla/mux"
)

// defines the accesspoint structure
type AccessPoint struct {
    Label string `json:"label"`
    URL string `json:"url"`
}

// definse the site structure
type Site struct {
    Name string `json:"name"`
    Role string `json:"role"`
    URI string `json:"uri"`
    APs []AccessPoint `json:"label"`
}


// called when a corresponding http request is passed in
// request should be sent with json formatted data to be place in the table
// ie: '{"name":"name1", "role":"role1","uri":"uri1","label":[{"ap":"ap1","url":"url1"}]}''
// generates and inserts entry that matches call provided the site doesn't already exist as a primary ket
// empty fields in the passed in json data will be left as empty strings, and the list of aps will be null unless specified
func CreateSite(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called CreateSite!")

    // Decode request JSON
    var site Site
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&site)
    if err != nil{
        fmt.Println(err)
    }

    // Open MySQL database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Prepare MySQL statement for creating site
    query := "INSERT INTO sites (name, role, uri) VALUES (?, ?, ?)"
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Execute MySQL statement
    _, err = stmt.Exec(site.Name, site.Role, site.URI)
    if err != nil {
        fmt.Println(err)
    }

    // Iterate through access points
    for _, AP := range site.APs{

        // Prepare MySQL statement for creating access point
    query = "INSERT INTO site_aps (label, url, name) VALUES (?, ?, ?)"
    stmt, err = db.Prepare(query)
        if err != nil {
            fmt.Println(err)
        }

        // Execute MySQL statement for inserting aps
        // APs will automatically be linked the the proper site because of the foreign key [name]
        _, err = stmt.Exec(AP.Label, AP.URL, site.Name)
        if err != nil {
            fmt.Println(err)
        }
    }
}

// called when a corresponding http request is passed in
// request will be processed to determine which site to be queried
// generates and prints out json object representing queried site to user
func ReadSite(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called ReadSite!")

    // Parse GET request
    request := mux.Vars(r)

    // Open MySQL database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Prepare MySQL query for sites
    query := "SELECT name, role, uri FROM sites WHERE name=?"
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Retrieve query results
    var name string
    var role string
    var uri string
    err = stmt.QueryRow(request["name"]).Scan(&name, &role, &uri)
    if err != nil{
        fmt.Println(err)
    }
    fmt.Println("result:", name, role, uri)

    // if name is empty string, no site name existed and no output should be produced
    if name == "" {
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "No such site found\n")
        return
    }

    // Prepare MySQL query for site's access points
    query = "SELECT label, url FROM site_aps WHERE name=?"
    stmt, err = db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Retrieve query results
    var APs []AccessPoint
    var label string
    var url string
    rows, err := stmt.Query(request["name"])
    if err != nil {
        fmt.Println(err)
    }
    for rows.Next(){
        err := rows.Scan(&label, &url)
        if err != nil {
            fmt.Println(err)
        }
        var AP AccessPoint
    AP.Label = label
    AP.URL = url
    APs = append(APs, AP)
    }

    // Construct JSON response
    var response Site
    response.Name = name
    response.Role = role
    response.URI = uri
    response.APs = APs
    response_json, err := json.Marshal(response)
    if err != nil {
        fmt.Println(err)
    }

    // Send response to client
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, string(response_json)+"\n")

}


// called when a corresponding http request is passed in
// will only update attributes of site:[site_name]
// name is primary key and cannot be updated, if change is desired,
// then user is responsible for deleting entry and creating updated entry to replace it
// is given json data that holds the attribute(s) to be updated and changes their value in the db
// access points must be updated via their own function call
func UpdateSite(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called UpdateSite!")

    var site Site

    // Parse PUT request
    request := mux.Vars(r)

    // Decode request JSON
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&site)
    if err != nil{
        fmt.Println(err)
    }

    // generates the base of the query and then fills it in once properly formatted
    query := "UPDATE sites SET "
    // if the role is not an empty string, user requested that it be updated, and it is added to the query
    if site.Role != "" {
        query += "role='"+site.Role+"'"
    }
    // if the uri is not an empty string, user requested that it be updated, and it is added to the query
    if site.URI != "" {
        // checks to see if role is also being updated, so a comma is needed for query format
        if site.Role != "" { query+=", "}
        query += "uri='"+site.URI+"'"
    }
    // adds on the last piece needed to generate a valid query
    query += " WHERE name=?"

    // Open MySQL Database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
	fmt.Println(query)

    // Prepares the query for execution
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Execute MySQL statement
    _, err = stmt.Exec(request["name"])
    if err != nil {
        fmt.Println(err)
    }

}


// called when a corresponding http request is passed in
// if no site by the requested name exists, function does nothing
// MySQL cascade on delete is used to ensure all underlying APs are deleted with the site
func DeleteSite(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called DeleteSite!")

    // Parse GET request
    request := mux.Vars(r)

    // Open MySQL database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Prepare statement
    stmt, err := db.Prepare("DELETE FROM sites WHERE name=?")
    if err != nil {
        fmt.Println(err)
    }

    // Delete row from sites table
    _, err = stmt.Exec(request["name"])
    if err != nil {
        fmt.Println(err)
    }
}


// called when a corresponding http request is passed in
// request should be sent with json formatted data to be place in the table
// ie: '{"ap":"ap1","url":"url1"}''
// generates and inserts entry that matches call provided the label doesn't already exist as a primary key
// empty fields in the passed in json data will be left as empty strings
// automatically associates AP with the user provided site by using site.name as foreign key
func CreateSiteAP(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called CreateSiteAP!")

    // Decode request JSON
    var siteAP AccessPoint
    request := mux.Vars(r)
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&siteAP)

    if err != nil{
        fmt.Println(err)
    }

    // Open MySQL database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()


    // Prepare MySQL statement for creating siteAP
    query := "INSERT INTO site_aps (label, url, name) VALUES (?, ?, ?)"
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Execute MySQL statement
    _, err = stmt.Exec(siteAP.Label, siteAP.URL, request["name"])
    if err != nil {
        fmt.Println(err)
    }

}

// called when a corresponding http request is passed in
// request will be processed to determine which site to be queried
// generates and prints out json object representing queried AP to user
func ReadSiteAP(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called ReadSiteAP!")

    // Parse GET request
    request := mux.Vars(r)

    // Open MySQL database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Prepare MySQL query for sites
    query := "SELECT label, url FROM site_aps WHERE label=?"
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Retrieve query results
    var label string
    var url string
    err = stmt.QueryRow(request["label"]).Scan(&label, &url)
    if err != nil{
        fmt.Println(err)
    }
    fmt.Println("result:", label, url)

    // if label is empty string, requested label does not exist and user should be notified
    if label == "" {
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "No such access point found\n")
        return
    }

    // Construct JSON response
    var response AccessPoint
    response.Label = label
    response.URL = url
    response_json, err := json.Marshal(response)
    if err != nil {
        fmt.Println(err)
    }

    // Send response to client
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, string(response_json) +"\n")

    // TODO: Shouldn't send a response if no results in db
}

// called when a corresponding http request is passed in
// will only update attributes of ap:[ap_name] (can only edit url of existing ap)
// label is primary key and cannot be updated, if change is desired,
// then user is responsible for deleting entry and creating updated entry to replace it
// is given json data that holds the attribute(s) to be updated and changes their value in the db
func UpdateSiteAP(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called UpdateSiteAP!")

    var accessPoint AccessPoint

    // Parse PUT request
    request := mux.Vars(r)

    // Decode request JSON
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&accessPoint)
    if err != nil{
        fmt.Println(err)
    }

    // generates query given users input
    // only url can be updated, so no additional checks needed
    query := "UPDATE site_aps SET url='"+accessPoint.URL+"' WHERE label=?"

    // Open MySQL Database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    fmt.Println(query)

    //now that we have the current record, update it to the new one
    //query = "UPDATE sites SET name=? WHERE name=?"
    stmt, err := db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Execute MySQL statement
    _, err = stmt.Exec(request["label"])
    if err != nil {
        fmt.Println(err)
    }

}

// called when a corresponding http request is passed in
// if no AP by the requested label exists, function does nothing
// deletes AP from table
func DeleteSiteAP(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called DeleteSiteAP!")

    // Parse GET request
    request := mux.Vars(r)

    // Open MySQL Database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()

    // Prepares Query
    stmt, err := db.Prepare("DELETE FROM site_aps WHERE label=?")
    if err != nil {
        fmt.Println(err)
    }

    // Executes Query
    _, err = stmt.Exec(request["label"])
    if err != nil {
        fmt.Println(err)
    }
}

// main function of program
// will cause terminal to hang while executing as it is listening for requests
// user can use curl to send requests to local host, responses will be sent to terminal user enters curl command from
// terminal that is running go program will print out debug information
func main(){

    // Creates a new router for handling curl requests
    router := mux.NewRouter()

    // Handlers for site API calls
    router.HandleFunc("/site/", CreateSite).Methods("POST")
    router.HandleFunc("/site/{name}", ReadSite).Methods("GET")
    router.HandleFunc("/site/{name}", UpdateSite).Methods("PUT") // ???
    router.HandleFunc("/site/{name}", DeleteSite).Methods("Delete")

    // Handlers for access point API calls
    router.HandleFunc("/site/{name}/ap/", CreateSiteAP).Methods("POST")
    router.HandleFunc("/ap/{label}", ReadSiteAP).Methods("GET")
    router.HandleFunc("/ap/{label}", UpdateSiteAP).Methods("PUT")
    router.HandleFunc("/ap/{label}", DeleteSiteAP).Methods("Delete")

    // Listen on port 8000 for REST API calls
    fmt.Println("Listening on 127.0.0.1:8000")
    http.ListenAndServe(":8000", router)
}