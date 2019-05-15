package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Sms struct to unmarshal inbound message json.
type Sms struct {
	ID           int64  `json:"id"`
	Msisdn       int64  `json:"msisdn"`
	Receiver     int64  `json:"receiver"`
	Message      string `json:"message"`
	Senttime     int    `json:"senttime"`
	WebhookLabel string `json:"webhook_label"`
}

func (s Sms) guessFromSms() (Guess, error) {
	split := strings.Split(s.Message, " ")
	matchID := split[1]
	// TODO verify that match id is valid.
	attendanceStr := split[2]
	attendance, err := strconv.ParseInt(attendanceStr, 10, 64)
	if err != nil {
		return Guess{}, err
	}
	guess := Guess{
		UserMsisdn: s.Msisdn,
		Total:      attendance,
		MatchID:    matchID,
	}
	return guess, nil
}

func (g Guess) respondToGuess() {
	response := fmt.Sprintf("Dit gæt på %d til kampen: %s er registret.", g.Total, g.MatchID)
	fmt.Print(response)
}

// Guess will contain values gorm should write to database for a guess.
type Guess struct {
	Total      int64  `gorm:"total"`
	UserMsisdn int64  `gorm:"user_msisdn"`
	MatchID    string `gorm:"match_id"`
}

func handleSms(w http.ResponseWriter, r *http.Request) {
	//conect to db
	db, err := gorm.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=attendance password=docker sslmode=disable")
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
	var sms Sms
	err = json.Unmarshal(b, &sms)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	guess, err := sms.guessFromSms()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	result := db.Table("guess").Create(&guess)
	if result.Error != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	guess.respondToGuess()
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
