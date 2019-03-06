
# CMSC 691 Project 1

## Installation and Dependencies

#### Instructions for installing on Ubuntu 18.10

Setting up Go Environment:
```
git clone https://github.com/Mjacks3/OS_Proj_1.git
cd OS_Proj_1
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
```

Setting up MySQL database:
```
sudo apt-get install mysql-server
sudo systemctl start mysql
sudo mysql -u root < init.sql
```

## Running the REST API server
```
go run proj1.go
```

## Interacting with the REST API
```
CREATE - curl --header "Content-Type: application/json" --request POST --data '[json data]' 127.0.0.1:8000/site/
READ - curl 127.0.0.1:8000/site/[name]
DELETE - curl -X "DELETE" 127.0.0.1:8000/site/[name]
UPDATE - ???
```
