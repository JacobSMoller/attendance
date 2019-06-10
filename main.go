package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JacobSMoller/attendance/guess"
	"github.com/JacobSMoller/attendance/match"
	"github.com/JacobSMoller/attendance/sms"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
)

// Config contains variables for configuring the services parsed from env.
type Config struct {
	DB        *gorm.DB
	GwKey     string `required:"true" split_words:"true"`
	AuthToken string `required:"true" split_words:"true"`
	DbHost    string `required:"true" split_words:"true"`
	DbName    string `required:"true" split_words:"true"`
	DbUser    string `required:"true" split_words:"true"`
	DbPw      string `required:"true" split_words:"true"`
}

var cfg Config

func handleSms(w http.ResponseWriter, r *http.Request) {
	gwJwt := r.Header.Get("X-Gwapi-Signature")
	token, err := jwt.Parse(gwJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(cfg.AuthToken), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Could not verify gwapi signature for request")
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

	newGuess, err := sms.GuessFromSms()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get todays match.
	match, err := match.TodaysMatch(cfg.DB)
	if err != nil {
		fmt.Println(err.Error())
		err = guess.SendMtsms("No match today", cfg.GwKey, newGuess.UserMsisdn)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	newGuess.MatchID = match.ID

	err = newGuess.GuessExists(cfg.DB, cfg.GwKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result := cfg.DB.Table("guess").Create(&newGuess)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), 500)
		return
	}
	err = newGuess.RespondToGuess(cfg.GwKey, match.HomeTeam, match.AwayTeam)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(newGuess)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	err := envconfig.Process("attendance", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	dbConnectString := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbName, cfg.DbPw)
	//conect to db
	db, err := gorm.Open(
		"postgres",
		dbConnectString,
	)
	db.LogMode(true)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	cfg.DB = db
	http.HandleFunc("/receive", handleSms)
	//Serve on port
	log.Fatal(http.ListenAndServe(":8080", nil))
}
