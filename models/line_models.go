package models

type FollowEvent struct {
	Source struct {
		Type   string `json:"type"`
		UserID string `json:"userId"`
	} `json:"source"`
	Timestamp       int64  `json:"timestamp"`
	Mode            string `json:"mode"`
	WebhookEventID  string `json:"webhookEventId"`
	DeliveryContext struct {
		IsRedelivery bool `json:"isRedelivery"`
	} `json:"deliveryContext"`
	ReplyToken string `json:"replyToken"`
	Follow     struct {
		IsUnblocked bool `json:"isUnblocked"`
	} `json:"follow"`
	Type string `json:"type"`
}
