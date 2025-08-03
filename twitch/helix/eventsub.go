package helix

type EventsubSubCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type EventsubSubWebhookTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type EventsubSub struct {
	ID        string          `json:"id"`
	Status    string          `json:"status"`
	Type      EventsubSubType `json:"type"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
	} `json:"transport"`
}
