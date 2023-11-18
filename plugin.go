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
	logger.Info(fmt.Sprintf("The plugin is now hijacking the path %s", path))

	// Handle the request for the path.
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// If the requested path is not what we defined, continue.
		if req.URL.Path != path {
			h.ServeHTTP(w, req)
			return
		}

		// The path has to be hijacked:
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
		auth := req.Header.Get("Authorization")
		if auth == "" {
			logger.Info("Please send Authorization header")
			return
		}

		logger.Info("########## Request path:", html.EscapeString(req.URL.Path))
		logger.Info("########## Authorization: ", auth)

		// Send the token to NATS.
		sendDataToNats(w, req)
	}), nil
}

func sendDataToNats(w http.ResponseWriter, req *http.Request) {
	logger.Info("########## Sending data to NATS.")
	var data BillingData
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("########## Data:", data.Client)
	logger.Info("########## Payment:", data.Payment)
	if data.Payment {
		logger.Info("########## Payment will be sent to NATS.")
		fmt.Fprintf(w, "Payment will be sent to NATS.")
	} else {
		logger.Info("########## Payment cancelled.")
		fmt.Fprintf(w, "Payment cancelled.")
	}
}

func main() {
	//p, err := plugin.Open("/Users/anupam.gogoi.br/projects/go/krakend-plugin/plugin/krakend-plugin.plugin")
	//if err != nil {
	//
	//}
	//fmt.Println(p)
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
	Client  string `json:"client"`
	Payment bool   `json:"payment"`
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
