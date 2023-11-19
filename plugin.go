package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
)

var pluginName = "krakend-plugin"
var HandlerRegisterer = registerer(pluginName)

type registerer string

//type arg1 func(name string, handler arg2)
//type arg2 func(ctx context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error)

func (r registerer) RegisterHandlers(f func(name string, handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error))) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {

	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	// Path to be intercepted in extra_config.path property.
	path, _ := config["path"].(string)
	logger.Info("########## Plugin configuration:", config)

	// Handle the request for the path.
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// If the requested path is not what we defined, continue.
		if req.URL.Path != path {
			h.ServeHTTP(w, req)
			return
		}

		// validateRequest validates the request.
		var resp []byte
		token, payment, data, err := validateRequest(w, req)
		if token == "" {
			resp = createResponse("Missing Authorization header", http.StatusForbidden)
		} else if !payment {
			resp = createResponse("Payment cancelled", http.StatusForbidden)
		} else if err != nil {
			resp = createResponse(err.Error(), http.StatusInternalServerError)
		} else {
			sendToNats(token, data)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		//

	}), nil
}

// createResponse creates JSON response to the consumer.
func createResponse(msg string, statusCode int32) []byte {
	resp := make(map[string]string)
	resp["msg"] = msg
	jsonResp, _ := json.Marshal(resp)
	return jsonResp
}
func sendToNats(token string, data BillingData) {

}
func validateRequest(w http.ResponseWriter, req *http.Request) (string, bool, BillingData, error) {
	logger.Info("########## Validate incoming request.")

	// Intercept the Authorization header.
	token := req.Header.Get("Authorization")
	logger.Info("########## Request path:", html.EscapeString(req.URL.Path))
	logger.Info("########## Authorization: ", token)

	// Intercept request payload.
	var data BillingData
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("########## Data:", data.Client)
	logger.Info("########## Payment:", data.Payment)
	return token, data.Payment, data, nil
}

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Info(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type BillingData struct {
	Client  int32 `json:"client"`
	Payment bool  `json:"payment"`
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
