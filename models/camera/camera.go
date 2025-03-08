package camera

import "time"

type CapturedEventData struct {
	EventID          string    `json:"EventId"`
	EventDescription string    `json:"EventDescription"`
	EventComment     string    `json:"EventComment"`
	ChannelName      string    `json:"ChannelName"`
	CapturedTime     time.Time `json:"captured_time"`
}

type CameraType string

const (
	Inside  CameraType = "inside"
	Outside CameraType = "outside"
)

type CapturedEventDataE struct {
	EventID          string    `json:"EventId"`
	EventDescription string    `json:"EventDescription"`
	EventComment     string    `json:"EventComment"`
	ChannelName      string    `json:"ChannelName"`
	CapturedTime     time.Time `json:"captured_time"`
	ChannelId        string    `json:"ChannelId"`
}

type Cameras struct {
	Id   int        `json:"id"`
	Name string     `json:"name"`
	Type CameraType `json:"type"`
}

type CamFix struct {
	Id          int        `json:"id"`
	ChannelName string     `json:"ChannelName"`
	ChannelId   string     `json:"ChannelId"`
	Type        CameraType `json:"type"`
}
