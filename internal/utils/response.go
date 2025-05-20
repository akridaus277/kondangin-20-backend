package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, status string, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendSuccess(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusOK), http.StatusOK, message, data)
}

func SendInternalServerError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, message, data)
}

func SendNotFoundError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, message, data)
}

func SendBadRequestError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, message, data)
}

func SendUnauthorizedError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, message, data)
}

func SendForbiddenError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusText(http.StatusForbidden), http.StatusForbidden, message, data)
}
