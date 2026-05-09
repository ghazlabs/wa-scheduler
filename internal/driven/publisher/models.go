package wa

import "encoding/json"

type RespSendMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *RespSendMessage) IsSessionExpired() bool {
	if r.Code == "AUTHENTICATION_ERROR" || r.Code == "SESSION_SAVED_ERROR" || r.Code == "INTERNAL_SERVER_ERROR" {
		return true
	}

	return false
}

func (r RespSendMessage) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}
