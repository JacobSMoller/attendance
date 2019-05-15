package sms

import (
	"strconv"
	"strings"

	"github.com/JacobSMoller/attendance/guess"
)

// Sms struct to unmarshal received webhook into.
type Sms struct {
	ID           int64  `json:"id"`
	Msisdn       int64  `json:"msisdn"`
	Receiver     int64  `json:"receiver"`
	Message      string `json:"message"`
	Senttime     int    `json:"senttime"`
	WebhookLabel string `json:"webhook_label"`
}

func (s Sms) guessFromSms() (guess.Guess, error) {
	split := strings.Split(s.Message, " ")
	matchID := split[1]
	// TODO verify that match id is valid.
	attendanceStr := split[2]
	attendance, err := strconv.ParseInt(attendanceStr, 10, 64)
	if err != nil {
		return guess.Guess{}, err
	}
	guess := guess.Guess{
		UserMsisdn: s.Msisdn,
		Total:      attendance,
		MatchID:    matchID,
	}
	return guess, nil
}
