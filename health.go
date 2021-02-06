package health

import "net/http"

var (
	handler = &healthHandler{
		useJSON: true,
		status:  Up,
	}
)

type healthHandler struct {
	useJSON bool
	status  Status
}

func (h *healthHandler) WithJSON(v bool) *healthHandler {
	h.useJSON = v
	return h
}

func (h healthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var status int
	var body []byte
	if h.status == Up {
		status = http.StatusOK
	} else {
		status = http.StatusInternalServerError
	}
	if h.useJSON {
		writer.Header().Set("Content-Type", "application/json")
		body = []byte(`{"status":"` + h.status + `"}`)
	} else {
		body = []byte(h.status)
	}
	writer.WriteHeader(status)
	_, _ = writer.Write(body)
}

func Handler() *healthHandler {
	return handler
}

func SetStatus(status Status) {
	handler.status = status
}
