package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JacobSMoller/attendance/sms"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func handleSms(w http.ResponseWriter, r *http.Request) {
	//conect to db
	db, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=attendance password=docker sslmode=disable"
	)
	defer db.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal json body to Sms.
	var sms sms.Sms
	err = json.Unmarshal(b, &sms)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	guess, err := sms.GuessFromSms()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = guess.GuessExists(db)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
	result := db.Table("guess").Create(&guess)
	if result.Error != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	guess.RespondToGuess()
	output, err := json.Marshal(guess)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func main() {
	http.HandleFunc("/receive", handleSms)
	//Connect to database
	log.Fatal(http.ListenAndServe(":8080", nil))
}
