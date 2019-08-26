package main

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/JacobSMoller/attendance/pkg/attendance"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg api.Config
	err := envconfig.Process("attendance", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	dbConnectString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbName, cfg.DbPw)
	//conect to db
	db, err := gorm.Open(
		"postgres",
		dbConnectString,
	)
	db.LogMode(true)
	fmt.Println("Connected")
	if err != nil {
		panic(err.Error())
	}
	db.LogMode(true)
	defer db.Close()
	cfg.DB = db
	service := api.NewAPIService(cfg)
	router := mux.NewRouter()
	router.HandleFunc("/receive", service.ReceiveSms).Methods("POST")
	router.HandleFunc("/match/{id:[0-9]+}", service.GetMatch).Methods("GET")
	router.HandleFunc("/match/create", service.CreateMatch).Methods("POST")
	router.HandleFunc("/match/update", service.UpdateMatch).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}
