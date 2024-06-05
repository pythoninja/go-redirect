package api

import (
	"github.com/pythoninja/go-redirect/internal/model"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"github.com/pythoninja/go-redirect/internal/vars"
	"log/slog"
	"net/http"
)

func (h *handler) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	health := model.Health{
		Version:     vars.Version,
		Environment: vars.Environment,
	}

	res, err := h.store.Health.GetDatabaseStatus()
	if err != nil {
		slog.Error(err.Error())
	} else {
		health.DatabaseStatus, health.Status = "connected", "ok"
	}

	// Check the result from 'select 1' query.
	if res != 1 {
		slog.Error("database down")

		health.DatabaseStatus, health.Status = "down", "fail"
	}

	json.Ok(w, r, health)
}
