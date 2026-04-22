package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse é a estrutura de resposta de erro
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HandleError trata um erro e envia a resposta apropriada
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	appErr, ok := err.(*AppError)
	if !ok {
		// Se não for um AppError, trata como erro interno
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    CodeInternalError,
			Message: "erro interno do servidor",
		})
		return
	}

	// Mapeia o código de erro para o status HTTP apropriado
	statusCode := mapErrorToStatus(appErr.Code)
	c.JSON(statusCode, ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	})
}

// mapErrorToStatus converte um código de erro para o status HTTP correspondente
func mapErrorToStatus(code string) int {
	switch code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeForbidden:
		return http.StatusForbidden
	case CodeConflict:
		return http.StatusConflict
	case CodeBadRequest, CodeValidation:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeInternalError:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
