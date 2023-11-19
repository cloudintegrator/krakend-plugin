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

}
