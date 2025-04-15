package handlers

import "log"

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func HandleError(c *gin.Context, err error, statusCode int) {
	log.Printf("Error: %v", err)

	c.JSON(statusCode, ErrorResponse{
		Message: err.Error(),
		Code:    statusCode,
	})
}
