package api

import (
	"strconv"
	"strings"
)

// Sms struct to unmarshal received webhook into.
type sms struct {
	ID           int64  `json:"id"`
	Msisdn       int64  `json:"msisdn"`
	Receiver     int64  `json:"receiver"`
	Message      string `json:"message"`
	Senttime     int    `json:"senttime"`
	WebhookLabel string `json:"webhook_label"`
}

func (s sms) GuessFromSms() (guess, error) {
	split := strings.Split(s.Message, " ")
	attendanceStr := split[1]
	attendance, err := strconv.ParseInt(attendanceStr, 10, 64)
	if err != nil {
		return guess{}, err
	}
	guess := guess{
		UserMsisdn: s.Msisdn,
		Total:      attendance,
	}
	return guess, nil
}
