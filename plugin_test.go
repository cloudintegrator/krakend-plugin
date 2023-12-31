package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateResponse(t *testing.T) {
	b := createResponse("Test", http.StatusCreated)

	var resp map[string]string
	json.Unmarshal(b, &resp)
	if resp["msg"] == "Test" {
		t.Logf("Success")
	} else {
		t.Errorf("Fail")
	}
}

func TestSendToNats(t *testing.T) {
	config := make(map[string]interface{})
	config["nats_url"] = "nats://127.0.0.1:4222"
	config["nats_topic"] = "BILLING.TOPIC"
	data := BillingData{
		Client:  0,
		Payment: false,
	}
	err := sendToNats("xxxxx", data, config)
	if err != nil {
		t.Logf("Success")
	} else {
		t.Errorf("Fail")
	}
}
