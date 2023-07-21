// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package inspect

// [START dlp_inspect_datastore_send_to_scc]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectDataStoreSendToScc inspects sensitive data in a data store and sends the results to Google Cloud Security Command Center (SCC).
func InspectDataStoreSendToScc(w io.Writer, projectID, datastoreNamespace, datastoreKind string) error {
	// projectID := "my-project-id"
	// datastoreNamespace := "your-datastore-namespace"
	// datastoreKind := "your-datastore-kind"

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the Datastore entity to be inspected.
	partitionId := &dlppb.PartitionId{
		ProjectId:   projectID,
		NamespaceId: datastoreNamespace,
	}

	// kindExpr represents an expression specifying a kind or range of kinds for data inspection in DLP.
	kindExpression := &dlppb.KindExpression{
		Name: datastoreKind,
	}

	// Specify datastoreOptions so that It holds the configuration options for inspecting data in
	// Google Cloud Datastore.
	datastoreOptions := &dlppb.DatastoreOptions{
		PartitionId: partitionId,
		Kind:        kindExpression,
	}

	// Specify the storageConfig to represents the configuration settings for inspecting data
	// in different storage types, such as BigQuery and Cloud Storage.
	storageConfig := &dlppb.StorageConfig{
		Type: &dlppb.StorageConfig_DatastoreOptions{
			DatastoreOptions: datastoreOptions,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoTypes := []*dlppb.InfoType{
		{Name: "EMAIL_ADDRESS"},
		{Name: "PERSON_NAME"},
		{Name: "LOCATION"},
		{Name: "PHONE_NUMBER"},
	}

	// The minimum likelihood required before returning a match.
	minLikelihood := dlppb.Likelihood_UNLIKELY

	// The maximum number of findings to report (0 = server maximum).
	findingLimits := &dlppb.InspectConfig_FindingLimits{
		MaxFindingsPerItem: 100,
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes:     infoTypes,
		MinLikelihood: minLikelihood,
		Limits:        findingLimits,
		IncludeQuote:  true,
	}

	// Specify the action that is triggered when the job completes.
	action := &dlppb.Action{
		Action: &dlppb.Action_PublishSummaryToCscc_{
			PublishSummaryToCscc: &dlppb.Action_PublishSummaryToCscc{},
		},
	}

	// Configure the inspection job we want the service to perform.
	inspectJobConfig := &dlppb.InspectJobConfig{
		StorageConfig: storageConfig,
		InspectConfig: inspectConfig,
		Actions: []*dlppb.Action{
			action,
		},
	}

	// Create the request for the job configured above.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: inspectJobConfig,
		},
	}

	// Send the request.
	resp, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		return err
	}

	// Print the result
	fmt.Fprintf(w, "Job created successfully: %v", resp.Name)
	return nil
}

// [END dlp_inspect_datastore_send_to_scc]
