package guess

import "fmt"

// Guess contains fields and info on a guess, including info to the ORM on how to store a guess in the database.
type Guess struct {
	Total      int64  `gorm:"total"`
	UserMsisdn int64  `gorm:"user_msisdn"`
	MatchID    string `gorm:"match_id"`
}

func (g Guess) respondToGuess() {
	response := fmt.Sprintf("Dit gæt på %d til kampen: %s er registret.", g.Total, g.MatchID)
	fmt.Print(response)
}
