package gcputil

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"os"
	"os/exec"
	"strings"
)

func GetDefaultProjectID(ctx context.Context) (string, error) {
	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/compute")
	if err != nil {
		return "", fmt.Errorf("failed to get GCP project from default credentials: %w", err)
	}

	gcpProjectID := credentials.ProjectID
	if gcpProjectID != "" {
		return gcpProjectID, nil
	}

	out := bytes.Buffer{}
	getProjectCmd := exec.CommandContext(ctx, "gcloud", "config", "get", "project")
	getProjectCmd.Stderr = os.Stderr
	getProjectCmd.Stdout = &out
	if err := getProjectCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get GCP project by running 'gcloud config get project': %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}
