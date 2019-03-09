package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "github.com/gorilla/mux"
)


type AccessPoint struct {
    Label string `json:"label"`
    URL string `json:"url"`
}

type Site struct {
    Name string `json:"name"`
    Role string `json:"role"`
    URI string `json:"uri"`
    APs []AccessPoint `json:"ap"`
}


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
	     query = "INSERT INTO site_aps (ap, url, name) VALUES (?, ?, ?)"
	     stmt, err = db.Prepare(query)
        if err != nil {
            fmt.Println(err)
        }

        // Execute MySQL statement
        _, err = stmt.Exec(AP.Label, AP.URL, site.Name)
        if err != nil {
            fmt.Println(err)
        }
    }
}


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

    // Prepare MySQL query for site's access points
    query = "SELECT ap, url FROM site_aps WHERE name=?"
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

    fmt.Println(label + " " + url)

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

    // TODO: Shouldn't send a response if no results in db
}


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

    //fmt.Println("Decoder: ",decoder)
    fmt.Println("request:", request["name"])

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

    // Prepare MySQL query for site's access points
    query = "SELECT ap, url FROM site_aps WHERE name=?"
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
    fmt.Println("updating record:", name, role, uri, label, url)
    fmt.Println("With new att: ", site.Name, site.Role, site.URI) // need to add the AP

    //now that we have the current record, update it to the new one
    query = "UPDATE sites SET name=? WHERE name=?"
    stmt, err = db.Prepare(query)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()

    // Execute MySQL statement
    _, err = stmt.Exec(site.Name, name)
    if err != nil {
        fmt.Println(err)
    }
/*
    // Iterate through access points
    for _, AP := range site.APs{

        // Prepare MySQL statement for creating access point
	     query = "INSERT INTO site_aps (ap, url, name) VALUES (?, ?, ?)"
	     stmt, err = db.Prepare(query)
        if err != nil {
            fmt.Println(err)
        }

        // Execute MySQL statement
        _, err = stmt.Exec(AP.Label, AP.URL, site.Name)
        if err != nil {
            fmt.Println(err)
        }
    }

    // Send response to client
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, string(response_json))
    */
}


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


func main(){
    router := mux.NewRouter()

    // Handlers for site API calls
    router.HandleFunc("/site/", CreateSite).Methods("POST")
    router.HandleFunc("/site/{name}", ReadSite).Methods("GET")
    router.HandleFunc("/site/{name}", UpdateSite).Methods("POST") // ???
    router.HandleFunc("/site/{name}", DeleteSite).Methods("Delete")

    // Listen on port 8000 for REST API calls
    fmt.Println("Listening on 127.0.0.1:8000")
    http.ListenAndServe(":8000", router)
}
