package modelscar

type Car_Model struct {
	ID            int     `json:"id"`
	Car_number    string  `json:"car_number"`
	Start_time    string  `json:"start_time"`
	End_time      string  `json:"end_time"`
	Total_payment float64 `json:"total_payment"`
	Status        string  `json:"status"`
	Reason        string  `json:"reason"`
	Image_Url     string  `json:"image_url"`
	ParkNo        string  `json:"park_no"`
	Duration      int     `json:"duration"`
	User_id       string  `json:"user_id"`
	PayStatus     bool    `json:"paystatus"`
	CameraID      string  `json:"cameraid"`
	CamToken      string  `json:"ChannelId"`
}

type CarUpdate struct {
	Reason        string  `json:"reason"`
	Total_payment float64 `json:"total_payment"`
}
