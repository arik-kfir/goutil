package gcputil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"os"
	"os/exec"
	"strings"
)

func getGoogleCloudCLIProject(ctx context.Context) (string, error) {
	out := bytes.Buffer{}
	getProjectCmd := exec.CommandContext(ctx, "gcloud", "config", "get", "project")
	getProjectCmd.Stderr = os.Stderr
	getProjectCmd.Stdout = &out
	if err := getProjectCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get GCP project from 'gcloud config get project': %w", err)
	} else {
		return strings.TrimSpace(out.String()), nil
	}
}

func GetDefaultProjectID(ctx context.Context) (string, error) {
	credentials, adcErr := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/compute")
	if adcErr != nil {
		if p, cliErr := getGoogleCloudCLIProject(ctx); cliErr != nil {
			return "", fmt.Errorf("failed to get GCP project: %w", errors.Join(adcErr, cliErr))
		} else {
			return p, nil
		}
	}

	gcpProjectID := credentials.ProjectID
	if gcpProjectID != "" {
		return gcpProjectID, nil
	}

	return getGoogleCloudCLIProject(ctx)
}
