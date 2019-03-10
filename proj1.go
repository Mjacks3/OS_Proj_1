package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "github.com/gorilla/mux"
)

// Defines the accesspoint structure
type AccessPoint struct {
    Label string `json:"label"`
    URL string `json:"url"`
}

// Defines the site structure
type Site struct {
    Name string `json:"name"`
    Role string `json:"role"`
    URI string `json:"uri"`
    APs []AccessPoint `json:"label"`
}


// Called via a POST request to /site/
// HTTP request should contain JSON of site to be created, ie:
// {
//   "name":"[name]",
//   "role":"[role]",
//   "uri":"[uri]",
//   "label":
//   [
//     {
//       "ap":"[ap]",
//       "url":"[url]"
//     }
//   ]
// }
// This will fail if the site already exists as a primary key in the database
// Empty fields in the passed in json data will be left as empty strings
// The list of access points will be null unless specified
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
        // APs will automatically be linked to the the proper site due to the
	// foreign key constraint on name
        _, err = stmt.Exec(AP.Label, AP.URL, site.Name)
        if err != nil {
            fmt.Println(err)
        }
    }
}

// Called via a GET request to /site/{name}
// Request will be processed to determine which site to be queried
// Generates and prints out JSON object representing queried site to user
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

    // If no results, print message to client and exit
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


// Called via PUT request to /site/{name}
// Will only update attributes of site:[site_name]
// Name is primary key and cannot be updated
// User must delete and re-create a site to change name
// Function accepts JSON data that holds the attribute(s) to be updated
// Updates the given values in the database
// Access points must be updated via their own function call
func UpdateSite(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called UpdateSite!")

    // Parse PUT request
    request := mux.Vars(r)

    // Decode request JSON
    var site Site
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&site)
    if err != nil{
        fmt.Println(err)
    }

    // Generates the base of the query
    // Rest of query is filled in once properly formatted
    query := "UPDATE sites SET "

    // Construct MySQL query from attributes given in HTTP PUT request
    if site.Role != "" {
        query += "role='"+site.Role+"'"
    }
    if site.URI != "" {
        if site.Role != "" { query+=", "}
        query += "uri='"+site.URI+"'"
    }

    // Last piece needed to complete MySQL query
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


// Called via DELETE request to /site/{name}
// If no site by the requested name exists, function does nothing
// MySQL cascade on delete ensures all associated APs are deleted too
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


// Called via POST request to /site/{name}/ap/
// Request should be sent with JSON data to add an access point, ie:
// {
//   "ap":"[ap]",
//   "url":"[url]"
// }
// AP label must not already exist as a primary key in the database
// Automatically associates AP and corresponding site via foreign key constraint
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


// Called via GET request to /ap/{label}
// Request will be processed to determine which site to be queried
// Generates and prints out JSON object representing queried AP to user
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

    // If label is not in database, notify user and return
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
}


// Called via a PUT request to /ap/{label}
// Will only update attributes of ap:[ap_name]
// Label is primary key and cannot be updated
// To update a label, must delete and create a new entry in database
// Accepts JSON data that holds the attribute(s) to be updated
// Changes the given values in the database
func UpdateSiteAP(w http.ResponseWriter, r *http.Request){

    fmt.Println("Called UpdateSiteAP!")

    // Parse PUT request
    request := mux.Vars(r)

    // Decode request JSON
    var accessPoint AccessPoint
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&accessPoint)
    if err != nil{
        fmt.Println(err)
    }

    // Generate MySQL Query
    // Only url can be updated, so no additional checks needed
    query := "UPDATE site_aps SET url='"+accessPoint.URL+"' WHERE label=?"

    // Open MySQL Database
    db, err := sql.Open("mysql", "proj1user:password@/proj1")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    fmt.Println(query)

    // Now that we have the current record, update it to the new one
    // query = "UPDATE sites SET name=? WHERE name=?"
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


// Called via a DELETE request to /ap/{label}/
// If no AP by the requested label exists, function does nothing
// Deletes AP from database table
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

// Main function of program
// Listens for HTTP requests to REST API on 127.0.0.1:9000
// Debug info is printed to the terminal
func main(){

    // Creates a new router for handling curl requests
    router := mux.NewRouter()

    // Handlers for site API calls
    router.HandleFunc("/site/", CreateSite).Methods("POST")
    router.HandleFunc("/site/{name}/", ReadSite).Methods("GET")
    router.HandleFunc("/site/{name}/", UpdateSite).Methods("PUT")
    router.HandleFunc("/site/{name}/", DeleteSite).Methods("DELETE")

    // Handlers for access point API calls
    router.HandleFunc("/site/{name}/ap/", CreateSiteAP).Methods("POST")
    router.HandleFunc("/ap/{label}/", ReadSiteAP).Methods("GET")
    router.HandleFunc("/ap/{label}/", UpdateSiteAP).Methods("PUT")
    router.HandleFunc("/ap/{label}/", DeleteSiteAP).Methods("DELETE")

    // Listen on port 8000 for REST API calls
    fmt.Println("Listening on 127.0.0.1:8000")
    http.ListenAndServe(":8000", router)
}
