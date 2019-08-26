package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// GetMatch used to handle /match/{id}.
func (s *Service) GetMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	match, err := GetMatchByID(s.DB, id)
	if err != nil {
		if err.Error() == "Match not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
	return
}

// CreateMatch used to handle /match/create
func (s *Service) CreateMatch(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var match Match
	err = json.Unmarshal(b, &match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO Unescape string fields before storing.
	result := s.DB.Table("match").Create(&match)
	if result.Error != nil {
		fmt.Println("failed to create match")
	}
	return
}

// UpdateMatch used to handle /match/update
func (s *Service) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var match Match
	err = json.Unmarshal(b, &match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(match)
	// TODO Unescape string fields before storing.
	result := s.DB.Table("match").Save(&match)
	if result.Error != nil {
		fmt.Println("failed to create match")
	}
	message := fmt.Sprintf("Match done todays attendance: %d", *match.Spectators)
	SendMtsms(message, s.GwKey, 4528725485)
	// TODO SEND SMS REPORT.
	// var guesses []Guess
	// result = db.Table("guess").Select("*").Where("match_id = ?", match.ID).Scan(&guesses)
	// if result.Error != nil {
	// 	fmt.Printf("Could not retrieve guesses")
	// }
	return
}
