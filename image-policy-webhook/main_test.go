package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImagePolicyWebhook(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main() // Call the main function to set up the handler
	})

	review := AdmissionReview{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Request: &AdmissionRequest{
			UID: "12345",
			Object: PodSpecWrapper{
				Spec: struct {
					Containers []struct {
						Image string `json:"image"`
					} `json:"containers"`
				}{
					Containers: []struct {
						Image string `json:"image"`
					}{
						{Image: "nginx:latest"},
					},
				},
			},
		},
	}

	body, _ := json.Marshal(review)
	req := httptest.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", resp.Code)
	}

	var response AdmissionReview
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Response == nil || response.Response.Allowed {
		t.Fatalf("Expected Allowed=false, got Allowed=true")
	}

	if response.Response.Status == nil || response.Response.Status.Message != "Image tag 'latest' is not allowed." {
		t.Fatalf("Unexpected status message: %v", response.Response.Status.Message)
	}
}
