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

type Cameras struct {
	Id   int        `json:"id"`
	Name string     `json:"name"`
	Type CameraType `json:"type"`
}

const (
	Inside  CameraType = "inside"
	Outside CameraType = "outside"
)
