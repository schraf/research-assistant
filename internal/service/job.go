package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/schraf/research-assistant/internal/models"
	"google.golang.org/api/option"
)

// executeResearchJob executes a Cloud Run Job with the research request
func executeResearchJob(ctx context.Context, logger *slog.Logger, req models.ResearchRequest) error {
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

	// Marshal to JSON and encode as base64
	requestJson, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal research request: %w", err)
	}

	encodedRequest := base64.StdEncoding.EncodeToString(requestJson)

	// Create Cloud Run Jobs client
	client, err := run.NewJobsClient(ctx, option.WithQuotaProject(projectId))
	if err != nil {
		return fmt.Errorf("failed to create Cloud Run Jobs client: %w", err)
	}
	defer client.Close()

	// Build the job execution request
	jobPath := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectId, region, jobName)

	// Execute the job with the request in an environment variable
	reqJob := &runpb.RunJobRequest{
		Name: jobPath,
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Env: []*runpb.EnvVar{
						{
							Name: "RESEARCH_REQUEST",
							Values: &runpb.EnvVar_Value{
								Value: encodedRequest,
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
		slog.String("request_id", req.RequestId),
		slog.String("request", encodedRequest),
	)

	return nil
}
