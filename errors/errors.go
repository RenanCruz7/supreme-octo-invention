package errors

import "fmt"

// AppError representa um erro customizado da aplicação
type AppError struct {
	Code    string // Código do erro para identificação (ex: "NOT_FOUND", "FORBIDDEN")
	Message string // Mensagem legível ao usuário
	Err     error  // Erro original (pode ser nil)
}

// Error implementa a interface error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap permite usar errors.Is() e errors.As()
func (e *AppError) Unwrap() error {
	return e.Err
}

// Constantes de códigos de erro
const (
	CodeNotFound      = "NOT_FOUND"
	CodeForbidden     = "FORBIDDEN"
	CodeConflict      = "CONFLICT"
	CodeBadRequest    = "BAD_REQUEST"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeInternalError = "INTERNAL_ERROR"
	CodeValidation    = "VALIDATION_ERROR"
)

// Construtores de erro

// ErrNotFound retorna um erro de recurso não encontrado
func ErrNotFound(message string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
	}
}

// ErrNotFoundWithErr retorna um erro de recurso não encontrado com o erro original
func ErrNotFoundWithErr(message string, err error) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
		Err:     err,
	}
}

// ErrForbidden retorna um erro de acesso proibido
func ErrForbidden(message string) *AppError {
	return &AppError{
		Code:    CodeForbidden,
		Message: message,
	}
}

// ErrConflict retorna um erro de conflito (ex: email já cadastrado)
func ErrConflict(message string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: message,
	}
}

// ErrBadRequest retorna um erro de requisição inválida
func ErrBadRequest(message string) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: message,
	}
}

// ErrBadRequestWithErr retorna um erro de requisição inválida com o erro original
func ErrBadRequestWithErr(message string, err error) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: message,
		Err:     err,
	}
}

// ErrUnauthorized retorna um erro de não autorizado
func ErrUnauthorized(message string) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

// ErrUnauthorizedWithErr retorna um erro de não autorizado com o erro original
func ErrUnauthorizedWithErr(message string, err error) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
		Err:     err,
	}
}

// ErrInternal retorna um erro interno do servidor
func ErrInternal(message string) *AppError {
	return &AppError{
		Code:    CodeInternalError,
		Message: message,
	}
}

// ErrInternalWithErr retorna um erro interno com o erro original
func ErrInternalWithErr(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternalError,
		Message: message,
		Err:     err,
	}
}

// ErrValidation retorna um erro de validação
func ErrValidation(message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
	}
}

// IsAppError verifica se um erro é do tipo AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extrai um AppError de um erro genérico
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
