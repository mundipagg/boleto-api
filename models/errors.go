package models

// IErrorResponse interface para implementar Error
type IErrorResponse interface {
	Error() string
	ErrorCode() string
}

// DataError objeto de erro
type ArrayDataError struct {
	Error []ErrorResponse `json:"error"`
}

// ErrorResponse objeto de erro
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//NewErrorResponse cria um novo objeto de ErrorReponse com código e mensagem
func NewErrorResponse(code, msg string) ErrorResponse {
	return ErrorResponse{Code: code, Message: msg}
}

// ErrorCode retorna código do erro
func (e ErrorResponse) ErrorCode() string {
	return e.Code
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// Errors coleção de erros
type Errors []ErrorResponse

// NewErrors cria nova coleção de erros vazia
func NewErrors() Errors {
	return []ErrorResponse{}
}

// NewErrorCollection cria nova coleção de erros
func NewErrorCollection(errorResponse ErrorResponse) Errors {
	return []ErrorResponse{errorResponse}
}

// NewSingleErrorCollection cria nova coleção de erros com 1 item
func NewSingleErrorCollection(code, msg string) Errors {
	return NewErrorCollection(NewErrorResponse(code, msg))
}

// GatewayTimeout objeto para erros 404 da aplicação: ex boleto não encontrado
type GatewayTimeout ErrorResponse

//NewGatewayTimeout cria um novo objeto NewGatewayTimeout a partir de uma mensagem original e final
func NewGatewayTimeout(code, msg string) GatewayTimeout {
	return GatewayTimeout{Message: msg, Code: code}
}

//ErrorCode ErrorCode
func (e GatewayTimeout) ErrorCode() string {
	return e.Code
}

//Error Error
func (e GatewayTimeout) Error() string {
	return e.Message
}

//InternalServerError IServerError interface para implementar Error
type InternalServerError ErrorResponse

//NewInternalServerError cria um novo objeto InternalServerError a partir de uma mensagem original e final
func NewInternalServerError(code, msg string) InternalServerError {
	return InternalServerError{Message: msg, Code: code}
}

//ErrorCode Message retorna a mensagem final para o usuário
func (e InternalServerError) ErrorCode() string {
	return e.Code
}

//Error retorna o erro original
func (e InternalServerError) Error() string {
	return e.Message
}

//HttpNotFound interface para implementar Error
type HttpNotFound ErrorResponse

//NewHTTPNotFound cria um novo objeto NewHttpNotFound a partir de uma mensagem original e final
func NewHTTPNotFound(code, msg string) HttpNotFound {
	return HttpNotFound{Message: msg, Code: code}
}

//ErrorCode Message retorna a mensagem final para o usuário
func (e HttpNotFound) ErrorCode() string {
	return e.Code
}

//Error retorna o erro original
func (e HttpNotFound) Error() string {
	return e.Message
}

//FormatError interface para implementar Error
type FormatError ErrorResponse

//NewFormatError cria um novo objeto de FormatError com descrição do erro
func NewFormatError(e string) FormatError {
	return FormatError{Message: e}
}

//Error Retorna um erro code
func (e FormatError) Error() string {
	return e.Message
}

//ErrorCode Retorna um erro code
func (e FormatError) ErrorCode() string {
	return e.Code
}

//Append adiciona mais um erro na coleção
func (e *Errors) Append(code, message string) {
	*e = append(*e, ErrorResponse{Code: code, Message: message})
}
