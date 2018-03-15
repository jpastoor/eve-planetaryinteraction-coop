package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	config := readConfig()
	dbConnStr := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Db.User, config.Db.Password, config.Db.Host, config.Db.Name)

	db := connectDb(dbConnStr)
	db.LogMode(false)

	defer db.Close()

	db.AutoMigrate(&Transaction{}, &Type{}, &Player{}, &LedgerState{}, &InventoryState{}, &Commit{})

	r := mux.NewRouter()

	s := NewServer(r, db)

	// Add CORS headers and start http listener
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	http.ListenAndServe(":1234", handlers.CORS(originsOk, headersOk, methodsOk)(s.r))
}

func connectDb(dbConnStr string) *gorm.DB {
	db, err := gorm.Open("mysql", dbConnStr)
	if err != nil {
		panic(err)
	}
	return db
}

func readConfig() Config {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var config Config
	json.Unmarshal(file, &config)
	return config
}

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Host     string
	Name     string
	User     string
	Password string
}
