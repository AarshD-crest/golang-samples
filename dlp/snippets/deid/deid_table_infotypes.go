// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deid

// [START dlp_deidentify_table_infotypes]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"google.golang.org/api/option"
)

func deidentifyTableInfotypes(w io.Writer, projectID string, table *dlppb.Table, columnNames ...string) error {
	// projectId := "your-project-id"
	// table := "your-table-value"
	// columnNames := "PATIENT","FACTOID"

	if table == nil {
		var row1 = &dlppb.Table_Row{
			Values: []*dlppb.Value{
				{Type: &dlppb.Value_StringValue{StringValue: "22"}},
				{Type: &dlppb.Value_StringValue{StringValue: "Jane Austen"}},
				{Type: &dlppb.Value_StringValue{StringValue: "21"}},
				{Type: &dlppb.Value_StringValue{StringValue: "There are 14 kisses in Jane Austen's novels."}},
			},
		}

		var row2 = &dlppb.Table_Row{
			Values: []*dlppb.Value{
				{Type: &dlppb.Value_StringValue{StringValue: "55"}},
				{Type: &dlppb.Value_StringValue{StringValue: "Mark Twain"}},
				{Type: &dlppb.Value_StringValue{StringValue: "75"}},
				{Type: &dlppb.Value_StringValue{StringValue: "Mark Twain loved cats."}},
			},
		}

		var row3 = &dlppb.Table_Row{
			Values: []*dlppb.Value{
				{Type: &dlppb.Value_StringValue{StringValue: "101"}},
				{Type: &dlppb.Value_StringValue{StringValue: "Charles Dickens"}},
				{Type: &dlppb.Value_StringValue{StringValue: "95"}},
				{Type: &dlppb.Value_StringValue{StringValue: "Charles Dickens name was a curse invented by Shakespeare."}},
			},
		}

		table = &dlppb.Table{
			Headers: []*dlppb.FieldId{
				{Name: "AGE"},
				{Name: "PATIENT"},
				{Name: "HAPPINESS SCORE"},
				{Name: "FACTOID"},
			},
			Rows: []*dlppb.Table_Row{
				{Values: row1.Values},
				{Values: row2.Values},
				{Values: row3.Values},
			},
		}
	}

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewRESTClient(ctx, option.WithCredentialsFile("C:/Users/aarsh.dhokai/Desktop/cred.json"))
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify what content you want the service to de-identify.
	var contentItem = &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Table{
			Table: table,
		},
	}

	// Specify how the content should be de-identified.
	// Select type of info to be replaced.
	var infoTypes = []*dlppb.InfoType{
		{Name: "PERSON_NAME"},
	}

	// Specify that findings should be replaced with corresponding info type name.
	var replaceWithInfoTypeConfig = &dlppb.ReplaceWithInfoTypeConfig{}
	var primitiveTransformation = &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{
			ReplaceWithInfoTypeConfig: replaceWithInfoTypeConfig,
		},
	}

	// Associate info type with the replacement strategy
	var infoTypeTransformations = &dlppb.InfoTypeTransformations{
		Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
			{
				InfoTypes:               infoTypes,
				PrimitiveTransformation: primitiveTransformation,
			},
		},
	}

	// Specify fields to be de-identified.
	var f []*dlppb.FieldId
	for _, c := range columnNames {
		f = append(f, &dlppb.FieldId{Name: c})
	}

	// Associate the de-identification and conditions with the specified field.
	var fieldTransformation = &dlppb.FieldTransformation{
		Fields: f,
		Transformation: &dlppb.FieldTransformation_InfoTypeTransformations{
			InfoTypeTransformations: infoTypeTransformations,
		},
	}

	var recordTransformations = &dlppb.RecordTransformations{
		FieldTransformations: []*dlppb.FieldTransformation{
			fieldTransformation,
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
				RecordTransformations: recordTransformations,
			},
		},
		Item: contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return fmt.Errorf("DeidentifyContent: %v", err)
	}

	// Print the results.
	fmt.Fprintf(w, "Table after de-identification : %v", resp.GetItem().GetTable())
	return nil
}

// [END dlp_deidentify_table_infotypes]
