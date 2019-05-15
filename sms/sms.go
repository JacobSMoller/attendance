package Sms

// Sms struct to unmarshal received webhook into.
type Sms struct {
	ID           int64  `json:"id"`
	Msisdn       int64  `json:"msisdn"`
	Receiver     int64  `json:"receiver"`
	Message      string `json:"message"`
	Senttime     int    `json:"senttime"`
	WebhookLabel string `json:"webhook_label"`
}
