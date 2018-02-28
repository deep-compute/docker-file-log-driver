package http

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/deep-compute/docker-file-log-driver/driver"
)

type startLoggingRequest struct {
	File string
	Info logger.Info
}

type stopLoggingRequest struct {
	File string
}

type capabilitiesResponse struct {
	Err string
	Cap logger.Capability
}

type response struct {
	Err string
}

func Handlers(h *sdk.Handler, d *driver.Driver) {
	h.HandleFunc("/LogDriver.StartLogging", func(w http.ResponseWriter, r *http.Request) {
		var req startLoggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := d.StartLogging(req.File, req.Info)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.StopLogging", func(w http.ResponseWriter, r *http.Request) {
		var req stopLoggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := d.StopLogging(req.File)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&capabilitiesResponse{
			Cap: logger.Capability{ReadLogs: false},
		})
	})
}

func respond(err error, w http.ResponseWriter) {
	var res response
	if err != nil {
		res.Err = err.Error()
	}
	json.NewEncoder(w).Encode(&res)
}
