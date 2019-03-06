package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    //TODO: will need "encoding/json"
    "github.com/gorilla/mux"
)


type Site struct {
    Name string `json:"name,omitempty"`
    Role string `json:"role,omitempty"`
    URI string `json:"uri,omitempty"`
    AP *AccessPoint `json:"ap"` // TODO: Allow 0+ access points
}

type AccessPoint struct {
    Label string `json:"label,omitempty"`
    URL string `json:"url,omitempty"`
}


func CreateSite(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called CreateSite!")
    // TODO
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

    // Query MySQL db
    query := "SELECT name, role, uri, ap FROM sites NATURAL LEFT JOIN " +
             "site_aps WHERE name=?"
    rows, err := db.Query(query, request["name"])
    if err != nil {
        fmt.Println(err)
    }

    // TODO: Parse query results into JSON
    var(
        name string
        role string
        uri string
        ap string
    )

    // Just printing the results for now
    for rows.Next(){
        err := rows.Scan(&name, &role, &uri, &ap)
        // TODO: An error is thrown here if NULL is inserted into ap
        if err != nil{
            fmt.Println(err)
	}
	fmt.Println("result:", name, role, uri, ap)
    }

    // TODO: Send response to client
}

func UpdateSite(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called UpdateSite!")
    // TODO
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


func CreateAP(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called CreateAP!")
    // TODO
}


func ReadAP(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called ReadAP!")
    // TODO
}


func UpdateAP(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called UpdateAP!")
    // TODO
}


func DeleteAP(w http.ResponseWriter, r *http.Request){
    fmt.Println("Called DeleteAP!")
    // TODO
}


func main(){
    router := mux.NewRouter()

    // Handlers for site API calls
    router.HandleFunc("/site/{name}", CreateSite).Methods("POST")
    router.HandleFunc("/site/{name}", ReadSite).Methods("GET")
    //router.HandleFunc("/site/{name}", UpdateSite).Methods("POST") ???
    router.HandleFunc("/site/{name}", DeleteSite).Methods("Delete")

    // Handlers for access point API calls
    router.HandleFunc("/ap/{label}", CreateSite).Methods("POST")
    router.HandleFunc("/ap/{label}", ReadSite).Methods("GET")
    //router.HandleFunc("/ap/{label}", UpdateSite).Methods("POST") ???
    router.HandleFunc("/ap/{label}", DeleteSite).Methods("Delete")

    // Listen on port 8000 for REST API calls
    fmt.Println("Listening on 127.0.0.1:8000")
    http.ListenAndServe(":8000", router)
}
