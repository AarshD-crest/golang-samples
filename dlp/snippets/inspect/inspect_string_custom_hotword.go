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

package inspect

// [START dlp_inspect_string_custom_hot_word]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"google.golang.org/api/option"
)

func inspectStringCustomHotWord(w io.Writer, projectID, textToInspect, customHotWord, infoTypeName string) error {
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewRESTClient(ctx, option.WithCredentialsFile("C:/Users/aarsh.dhokai/Desktop/cred.json"))
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	var contentItem = &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: dlppb.ByteContentItem_TEXT_UTF8,
				Data: []byte(textToInspect),
			},
		},
	}

	// Increase likelihood of matches that have customHotword nearby
	var hotwordRule = &dlppb.InspectionRule_HotwordRule{

		HotwordRule: &dlppb.CustomInfoType_DetectionRule_HotwordRule{
			HotwordRegex: &dlppb.CustomInfoType_Regex{
				Pattern: customHotWord,
			},
			Proximity: &dlppb.CustomInfoType_DetectionRule_Proximity{
				WindowBefore: 50,
			},
			LikelihoodAdjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment{
				Adjustment: &dlppb.CustomInfoType_DetectionRule_LikelihoodAdjustment_FixedLikelihood{
					FixedLikelihood: dlppb.Likelihood_VERY_LIKELY,
				},
			},
		},
	}

	var ruleSet = &dlppb.InspectionRuleSet{
		InfoTypes: []*dlppb.InfoType{
			{Name: infoTypeName}, //"PERSON_NAME"
		},
		Rules: []*dlppb.InspectionRule{
			{
				Type: &dlppb.InspectionRule_HotwordRule{
					HotwordRule: hotwordRule.HotwordRule,
				},
			},
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Item:   contentItem,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{Name: infoTypeName}, //"PERSON_NAME"
			},
			IncludeQuote:  true,
			MinLikelihood: dlppb.Likelihood_VERY_LIKELY,
			RuleSet: []*dlppb.InspectionRuleSet{
				ruleSet,
			},
		},
	}

	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "Receive: %v", err)
		return err
	}

	fmt.Fprintf(w, "Findings: %v\n", len(resp.Result.Findings))
	for _, v := range resp.GetResult().Findings {
		fmt.Fprintf(w, "Quote: %v\n", v.GetQuote())
		fmt.Fprintf(w, "Infotype Name: %v\n", v.GetInfoType().GetName())
		fmt.Fprintf(w, "Likelihood: %v\n", v.GetLikelihood())
	}
	return nil

}

// [END dlp_inspect_string_custom_hotword]
