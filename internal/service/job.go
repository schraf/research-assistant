package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/schraf/research-assistant/internal/models"
	"google.golang.org/api/option"
)

// executeResearchJob executes a Cloud Run Job with the research request
func executeResearchJob(ctx context.Context, logger *slog.Logger, requestId string, topic string) error {
	// Get Cloud Run Job name and region from environment
	jobName := os.Getenv("CLOUD_RUN_JOB_NAME")
	if jobName == "" {
		return fmt.Errorf("CLOUD_RUN_JOB_NAME environment variable is not set")
	}

	region := os.Getenv("CLOUD_RUN_JOB_REGION")
	if region == "" {
		return fmt.Errorf("CLOUD_RUN_JOB_REGION environment variable is not set")
	}

	projectId := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectId == "" {
		return fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	// Create the research request payload
	req := models.ResearchRequest{
		RequestId: requestId,
		Topic:     topic,
	}

	// Marshal to JSON and encode as base64 for CloudEvent data
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal research request: %w", err)
	}

	encodedData := base64.StdEncoding.EncodeToString(payload)

	// Create Cloud Run Jobs client
	client, err := run.NewJobsClient(ctx, option.WithQuotaProject(projectId))
	if err != nil {
		return fmt.Errorf("failed to create Cloud Run Jobs client: %w", err)
	}
	defer client.Close()

	// Build the job execution request
	jobPath := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectId, region, jobName)

	// Create CloudEvent payload for the job
	cloudEventData := map[string]interface{}{
		"data": map[string]interface{}{
			"message": map[string]interface{}{
				"data": encodedData,
			},
		},
	}

	eventData, err := json.Marshal(cloudEventData)
	if err != nil {
		return fmt.Errorf("failed to marshal CloudEvent data: %w", err)
	}

	// Execute the job with the CloudEvent data as environment variable
	eventDataStr := string(eventData)
	reqJob := &runpb.RunJobRequest{
		Name: jobPath,
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Env: []*runpb.EnvVar{
						{
							Name: "CLOUDEVENT_DATA",
							Values: &runpb.EnvVar_Value{
								Value: eventDataStr,
							},
						},
					},
				},
			},
		},
	}

	_, err = client.RunJob(ctx, reqJob)
	if err != nil {
		return fmt.Errorf("failed to execute Cloud Run Job: %w", err)
	}

	logger.Info("job_execution_started",
		slog.String("job_name", jobName),
		slog.String("request_id", requestId),
	)

	return nil
}
