// Copyright 2021 Google LLC
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

package buckets

// [START storage_set_public_access_prevention_enforced]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setPublicAccessPreventionEnforced sets public access prevention to
// "enforced" for the bucket.
func setPublicAccessPreventionEnforced(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	setPublicAccessPrevention := storage.BucketAttrsToUpdate{
		PublicAccessPrevention: storage.PublicAccessPreventionEnforced,
	}
	if _, err := bucket.Update(ctx, setPublicAccessPrevention); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Public access prevention is 'enforced' for %v", bucketName)
	return nil
}

// [END storage_set_public_access_prevention_enforced]
