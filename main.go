package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/JacobSMoller/attendance/guess"
	"github.com/JacobSMoller/attendance/match"
	"github.com/JacobSMoller/attendance/sms"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Env contains server setup.
type Service struct {
	DB    *gorm.DB
	GwKey string
}

func (s *Service) handleSms(w http.ResponseWriter, r *http.Request) {
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

	newGuess, err := sms.GuessFromSms()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get todays match.
	match, err := match.TodaysMatch(s.DB)
	if err != nil {
		fmt.Println(err.Error())
		guess.SendMtsms("No match today", s.GwKey, newGuess.UserMsisdn)
		return
	}
	newGuess.MatchID = match.ID

	err = newGuess.GuessExists(s.DB, s.GwKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result := s.DB.Table("guess").Create(&newGuess)
	if result.Error != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	newGuess.RespondToGuess(s.GwKey, match.HomeTeam, match.AwayTeam)
	output, err := json.Marshal(newGuess)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func main() {
	//conect to db
	db, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=attendance password=docker sslmode=disable",
	)
	db.LogMode(true)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	gwKey := os.Getenv("GWKEY")
	if gwKey == "" {
		panic("GWKEY env variable not found.")
	}
	env := &Service{
		DB:    db,
		GwKey: gwKey,
	}
	http.HandleFunc("/receive", env.handleSms)
	//Connect to database
	log.Fatal(http.ListenAndServe(":8080", nil))
}
