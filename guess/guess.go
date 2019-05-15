package guess

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Guess contains fields and info on a guess, including info to the ORM on how to store a guess in the database.
type Guess struct {
	Total      int64  `gorm:"total"`
	UserMsisdn int64  `gorm:"user_msisdn"`
	MatchID    string `gorm:"match_id"`
}

// Sends sms
func (g Guess) RespondToGuess() {
	response := fmt.Sprintf("Dit gæt på %d til kampen: %s er registret.", g.Total, g.MatchID)
	fmt.Println(response)
}

// GuessFromDB check if a guess already exists for match and user.
func (g Guess) GuessExists(db *gorm.DB) error {
	result := db.Table("guess").Where("user_msisdn = ? and match_id = ?", g.UserMsisdn, g.MatchID).First(&Guess{})
	if !result.RecordNotFound() {
		return fmt.Errorf("You have already guessed on match %q", g.MatchID)
	}
	return nil
}
