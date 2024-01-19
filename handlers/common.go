package handlers

import (
	"blog-api/pkg/logging"
	"net/http"
)

// Logger - общий логгер для использования во всех обработчиках
var Logger = logging.NewLogger().GetLogger()

// HandleError - обработка ошибки, логирование и возврат HTTP-ответа
func HandleError(w http.ResponseWriter, message string, err error, statusCode int) {
	Logger.Error(message, "Error: ", err)
	http.Error(w, message, statusCode)
}
