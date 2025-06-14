package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// AdmissionReview and related structs
// Minimal types for decoding/encoding

type AdmissionReview struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Request    *AdmissionRequest  `json:"request,omitempty"`
	Response   *AdmissionResponse `json:"response,omitempty"`
}

type AdmissionRequest struct {
	UID      string               `json:"uid"`
	Kind     GroupVersionKind     `json:"kind"`
	Resource GroupVersionResource `json:"resource"`
	Object   PodSpecWrapper       `json:"object"`
}

type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

type GroupVersionResource struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
}

type AdmissionResponse struct {
	UID     string  `json:"uid"`
	Allowed bool    `json:"allowed"`
	Status  *Status `json:"status,omitempty"`
}

type Status struct {
	Message string `json:"message"`
}

type PodSpecWrapper struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Containers []struct {
			Image string `json:"image"`
		} `json:"containers"`
	} `json:"spec"`
}

func main() {
	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var review AdmissionReview
		fmt.Printf("Received request: %s\n", body)
		if err := json.Unmarshal(body, &review); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		allowed := true
		msg := ""
		if review.Request != nil {
			for _, c := range review.Request.Object.Spec.Containers {
				if isLatestTag(c.Image) {
					allowed = false
					msg = "Image tag 'latest' is not allowed."
					break
				}
			}
		}
		resp := AdmissionReview{
			APIVersion: review.APIVersion,
			Kind:       review.Kind,
			Response: &AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: allowed,
			},
		}
		if !allowed {
			resp.Response.Status = &Status{Message: msg}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	cert, err := tls.LoadX509KeyPair("/etc/ssl/server.crt", "/etc/ssl/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS cert/key: %v", err)
	}
	server := &http.Server{
		Addr:    ":8443",
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	log.Println("Starting HTTPS webhook on :8443 ...")
	// log.Println("Starting HTTP webhook on :8081 ...")
	log.Fatal(server.ListenAndServeTLS("", ""))
	// log.Fatal(server.ListenAndServe())
}

func isLatestTag(image string) bool {
	// Accepts image[:tag], returns true if tag is 'latest' or missing (default is latest)
	// e.g. nginx:latest, nginx, repo/nginx:latest
	colon := -1
	for i := len(image) - 1; i >= 0; i-- {
		if image[i] == ':' {
			colon = i
			break
		}
	}
	if colon == -1 {
		return true // no tag, default is latest
	}
	tag := image[colon+1:]
	return tag == "latest"
}
