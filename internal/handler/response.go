package handler

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

func SuccessResponse(message string, result interface{}) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
		Result:  result,
	}
}

func ErrorResponse(status string, message string) BaseResponse {
	return BaseResponse{
		Status:  status,
		Message: message,
	}
}
