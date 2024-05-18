package api

import (
	"github.com/pythoninja/go-redirect/internal/model"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"github.com/pythoninja/go-redirect/internal/vars"
	"log/slog"
	"net/http"
)

var helloResponse = "hello"

func (h *handler) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.store.GetDatabaseStatus()
	if err != nil {
		slog.Error(err.Error())
	}

	var health model.Health

	if res == helloResponse {
		health.DatabaseStatus = "connected"
		health.Status = "ok"
	} else {
		slog.Error("database down")
		health.DatabaseStatus = "down"
		health.Status = "fail"
	}

	health.Version = vars.Version
	health.Environment = vars.Environment

	json.OK(w, r, health)
}
