package response

import (
	"net/http"

	"talk-backend/internal/http/dto"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, dto.ErrorResponse{
		Code:  code,
		Error: message,
	})
}

func ErrorWithDetails(c *gin.Context, status int, code, message, details string) {
	c.JSON(status, dto.ErrorResponse{
		Code:    code,
		Error:   message,
		Details: details,
	})
}

func InvalidBody(c *gin.Context, err error) {
	ErrorWithDetails(c, http.StatusBadRequest, CodeInvalidRequest, MsgInvalidRequestBody, err.Error())
}
