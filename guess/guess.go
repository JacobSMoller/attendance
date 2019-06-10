package guess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Guess contains fields and info on a guess, including info to the ORM on how to store a guess in the database.
type Guess struct {
	Total      int64  `gorm:"total"`
	UserMsisdn int64  `gorm:"user_msisdn"`
	MatchID    uint32 `gorm:"match_id"`
}

type mtsms struct {
	Message    string             `json:"message"`
	Sender     string             `json:"sender"`
	Recipients []map[string]int64 `json:"recipients"`
}

// SendMtsms sends a sms message to given number.
func SendMtsms(message, key string, msisdn int64) error {
	sms := mtsms{
		Message: message,
		Sender:  "Attendance",
	}
	sms.Recipients = []map[string]int64{
		map[string]int64{
			"msisdn": msisdn,
		},
	}
	jsonSms, err := json.Marshal(sms)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://gatewayapi.com/rest/mtsms", bytes.NewBuffer(jsonSms))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(key, "")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non 200 (%q) code received from GatewayAPI", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil

}

// RespondToGuess sends a response to the msisdn of the guess, stating that guess is received.
func (g Guess) RespondToGuess(key, homeTeam, awayTeam string) error {
	message := fmt.Sprintf("Dit gæt på %d til dagens kamp mellem %s og %s er registreret.", g.Total, homeTeam, awayTeam)
	err := SendMtsms(message, key, g.UserMsisdn)
	if err != nil {
		return err
	}
	return nil
}

// GuessExists check if a guess already exists for match and user.
func (g Guess) GuessExists(db *gorm.DB, key string) error {
	var currentGuess Guess
	result := db.Table("guess").Select("total").Where("user_msisdn = ? and match_id = ?", g.UserMsisdn, g.MatchID).Scan(&currentGuess)
	if !result.RecordNotFound() {
		message := fmt.Sprintf("Du har allerede gættet på %d tilskuere til dagens kamp.", currentGuess.Total)
		SendMtsms(message, key, g.UserMsisdn)
		return fmt.Errorf("User %d has already guessed on match %d", g.UserMsisdn, g.MatchID)
	}
	return nil
}
