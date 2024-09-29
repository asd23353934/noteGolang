package apperrors

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseHandler struct{}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

func (h *ResponseHandler) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (h *ResponseHandler) Error(c *gin.Context, err error) {
	var statusCode int
	var errorMsg string
	log.Println(err)
	switch e := err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound
		errorMsg = e.Error()
	default:
		statusCode = http.StatusInternalServerError
		errorMsg = "Internal Server Error"
	}

	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode,
		Message: errorMsg,
	})
}
