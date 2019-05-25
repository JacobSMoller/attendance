package match

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Match to be stored in the database.
type Match struct {
	ID         uint32     `gorm:"id"`
	Tournament string     `gorm:"tournament"`
	StartTime  *time.Time `gorm:"start_time"`
	Spectators *int       `gorm:"spectators"`
	State      string     `gorm:"state"`
	Referee    string     `gorm:"referee"`
	HomeTeam   string     `gorm:"home_team"`
	AwayTeam   string     `gorm:"away_team"`
}

// returns date at beginning of day.
func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// TodaysMatch returns todays match if any is found, otherwise err.
func TodaysMatch(db *gorm.DB) (*Match, error) {
	var match Match
	dayStart := bod(time.Now().UTC())
	dayEnd := dayStart.Add(24 * time.Hour)
	result := db.Table("match").Select("*").Where("start_time BETWEEN ? AND ?", dayStart, dayEnd).Scan(&match)
	if result.Error != nil {
		return nil, result.Error
	}
	return &match, nil
}
