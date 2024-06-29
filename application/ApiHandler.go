package application

import (
	"fmt"
	"net/http"

	"ward4woods.ca/helpers"
	"ward4woods.ca/models"
)

type ApiHandler struct {
	logger *Logger
}

func NewApiHandler(logger *Logger) *ApiHandler {
	apih := &ApiHandler{
		logger: logger,
	}
	apih.logger.Info("Created ApiHandler successfully.", nil)
	return apih
}

func (apih *ApiHandler) Handle(responseWriter ResponseWriter, w http.ResponseWriter, data interface{}) {
	err := responseWriter(w, data)
	if err != nil {
		apih.logger.Error("Error when trying to write response.", err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	apih.logger.Info(fmt.Sprintf("Now generating response with data: [%+v]", data), nil)
}

// TryWriteError checks if error is nil. If it is not, it writes the error response and returns true.
// If the error is nil, it doesn nothing then returns false.
func (apih *ApiHandler) TryWriteError(errorWriter ErrorWriter, w http.ResponseWriter, err error) bool {
	if err != nil {
		errorWriter(w, err)
		apih.logger.Info("Now serving error page.", nil)
		return true
	}
	return false
}

func GetIdFromApiRequest(r *http.Request) (models.ProductId, error) {
	prefix := "api/products"
	return helpers.GetIdFromRequest(r, prefix)
}
