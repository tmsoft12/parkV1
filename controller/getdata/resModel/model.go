package resmodel

import modelscar "park/models/modelsCar"

type Response struct {
	Message string              `json:"message"`
	Data    modelscar.Car_Model `json:"data"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}
