package web

type ResponseWithData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithoutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func JSONResponseWithData(code int, status string, message string, data interface{}) ResponseWithData {
	jsonResponse := ResponseWithData{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}

	return jsonResponse
}

func JSONResponseWithoutData(code int, status string, message string) ResponseWithoutData {
	jsonResponse := ResponseWithoutData{
		Code:    code,
		Status:  status,
		Message: message,
	}

	return jsonResponse
}
