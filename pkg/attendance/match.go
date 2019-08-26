package api

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Match to be stored in the database.
type Match struct {
	ID         uint32     `gorm:"id" json:"id"`
	Tournament string     `gorm:"tournament" json:"tournament"`
	StartTime  *time.Time `gorm:"start_time" json:"start_time"`
	Spectators *int       `gorm:"spectators" json:"spectators"`
	State      string     `gorm:"state" json:"state"`
	Referee    string     `gorm:"referee" json:"referee"`
	HomeTeam   string     `gorm:"home_team" json:"home_team"`
	AwayTeam   string     `gorm:"away_team" json:"away_team"`
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

// GetMatchByID looks up matches by their id.
func GetMatchByID(db *gorm.DB, matchID string) (*Match, error) {
	var match Match
	result := db.Table("match").Select("*").Where("id = ?", matchID).Scan(&match)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, fmt.Errorf("Match not found")
		}
		return nil, result.Error
	}
	return &match, nil
}
