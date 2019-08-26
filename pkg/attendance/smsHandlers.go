package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ReceiveSms handles incoming smses from gateway.
func (s *Service) ReceiveSms(w http.ResponseWriter, r *http.Request) {
	err := s.VerifyGWRequest(r)
	if err != nil {
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
	var sms sms
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
	match, err := TodaysMatch(s.DB)
	if err != nil {
		fmt.Println(err.Error())
		err = SendMtsms("No match today", s.GwKey, newGuess.UserMsisdn)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
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
		http.Error(w, result.Error.Error(), 500)
		return
	}
	err = newGuess.RespondToGuess(s.GwKey, match.HomeTeam, match.AwayTeam)
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
