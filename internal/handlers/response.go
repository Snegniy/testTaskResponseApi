package handlers

import (
	"encoding/json"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func writeJSON[T any](w http.ResponseWriter, info T) {
	response, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		logger.Error("error writing JSON", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		logger.Error("error writing JSON", zap.Error(err))
	}
}
