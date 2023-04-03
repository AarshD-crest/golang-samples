// Copyright 2019 Google LLC
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

package inspect

// [START dlp_inspect_hotword_rule]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// inspectWithHotWordRules inspect the data with hotword rule, it uses a custom regex
// with hotword rule to increase the likelihood match
func inspectWithHotWordRules(w io.Writer, projectID, textToInspect, customRegexPattern, hotWordRegexPattern, infoTypeName string) error {
	//projectID := "my-project-id"
	//textToInspect := "Patient's MRN 444-5-22222 and just a number 333-2-33333"
	//customRegexPattern := "[1-9]{3}-[1-9]{1}-[1-9]{5}"
	//hotWordRegexPattern := "(?i)(mrn|medical)(?-i)"
	//infoTypeName := "C_MRN"
	ctx := context.Background()
	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Specify the type and content to be inspected.
	var contentItem = &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: dlppb.ByteContentItem_TEXT_UTF8,
				Data: []byte(textToInspect),
			},
		},
	}

	// Construct the custom regex detectors
	var customInfoType = &dlppb.CustomInfoType{
		InfoType: &dlppb.InfoType{
			Name: infoTypeName,
		},
		// Specify the regex pattern the inspection will look for.
		Type: &dlppb.CustomInfoType_Regex_{
			Regex: &dlppb.CustomInfoType_Regex{
				Pattern: customRegexPattern,
			},
		},
		MinLikelihood: dlppb.Likelihood_POSSIBLE,
	}

	var inspectionRuleSet = &dlppb.InspectionRuleSet{
		Rules: []*dlppb.InspectionRule{
			{
				// Construct hotword rule.
				Type: &dlppb.InspectionRule_HotwordRule{
					HotwordRule: &dlppb.CustomInfoType_DetectionRule_HotwordRule{
						HotwordRegex: &dlppb.CustomInfoType_Regex{
							Pattern: hotWordRegexPattern,
						},
						// Specify a window around a finding to apply a detection rule.
						Proximity: &dlppb.CustomInfoType_DetectionRule_Proximity{
							WindowBefore: int32(10),
						},
						// Specify hotword likelihood adjustment.
						LikelihoodAdjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
							Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
								FixedLikelihood: dlppb.Likelihood_VERY_LIKELY,
							},
						},
					},
				},
			},
		},
		InfoTypes: []*dlppb.InfoType{
			customInfoType.InfoType,
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:   contentItem,
		// Construct the configuration for the Inspect request.
		InspectConfig: &dlppb.InspectConfig{
			CustomInfoTypes: []*dlppb.CustomInfoType{
				customInfoType,
			},
			// Construct rule set for the inspect config.
			RuleSet: []*dlppb.InspectionRuleSet{
				inspectionRuleSet,
			},
			IncludeQuote: true,
		},
	}

	// Send the request.
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "Receive: %v", err)
		return err
	}

	// Parse the response and process results
	fmt.Fprintf(w, "Findings: %v\n", len(resp.Result.Findings))
	for _, v := range resp.GetResult().Findings {
		fmt.Fprintf(w, "Quote: %v\n", v.GetQuote())
		fmt.Fprintf(w, "InfoType Name: %v\n", v.GetInfoType().GetName())
		fmt.Fprintf(w, "Likelihood: %v\n", v.GetLikelihood())
	}
	return nil

}

// [END dlp_inspect_hotword_rule]
