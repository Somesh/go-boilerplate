package types

import (
	"time"
)

type ExtraConfiguration map[string]interface{}
type TrackingInfo map[string]interface{}

type Address struct{}


type StatusMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type StatusResponse struct {
	Result  string        `json:"result"`
	Message StatusMessage `json:"message"`
	Code    int64         `json:"code"`
}
