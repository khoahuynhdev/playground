package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestK8sIntegration(t *testing.T) {
	// Step 1: Create a KIND cluster
	cmd := exec.Command("kind", "create", "cluster", "--name", "image-policy-test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create KIND cluster: %v", err)
	}
	defer func() {
		exec.Command("kind", "delete", "cluster", "--name", "image-policy-test").Run()
	}()

	// Step 2: Deploy the webhook server
	cmd = exec.Command("kubectl", "apply", "-f", "webhook-deployment.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to deploy webhook server: %v", err)
	}

	// Step 3: Apply a sample deployment with nginx:latest
	deployment := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest`

	file, err := os.Create("nginx-deployment.yaml")
	if err != nil {
		t.Fatalf("Failed to create deployment file: %v", err)
	}
	defer os.Remove("nginx-deployment.yaml")
	file.WriteString(deployment)
	file.Close()

	cmd = exec.Command("kubectl", "apply", "-f", "nginx-deployment.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err == nil {
		t.Fatalf("Expected deployment to be rejected, but it was applied successfully")
	}
}
