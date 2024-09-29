package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// 用於自定義驗證錯誤消息的映射
var customValidationMessages = map[string]string{
	"required": "Field '%s' is required",
	"email":    "Field '%s' must be a valid email address",
}

// BindJSON 是一個通用的 JSON 綁定和驗證函數
func BindJSON[T any](c *gin.Context, obj *T) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		fmt.Println("Failed to bind JSON", "error", err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 處理驗證錯誤
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				tag := e.Tag()
				if msg, exists := customValidationMessages[tag]; exists {
					errorMessages[field] = fmt.Sprintf(msg, field)
				} else {
					errorMessages[field] = fmt.Sprintf("Validation failed on field '%s' with constraint: %s", field, tag)
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		} else {
			// 處理其他類型的錯誤
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		}

		return false
	}
	return true
}
