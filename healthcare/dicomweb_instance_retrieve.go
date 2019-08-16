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

package snippets

// [START healthcare_dicomweb_retrieve_instance]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// dicomWebRetrieveInstance retrieves a specific instance.
func dicomWebRetrieveInstance(w io.Writer, projectID, location, datasetID, dicomStoreID, dicomWebPath string, outputFile string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.DicomStores.Studies.Series.Instances

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, dicomStoreID)

	call := storesService.RetrieveInstance(parent, dicomWebPath)
	call.Header().Set("Accept", "application/dicom; transfer-syntax=*")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("RetrieveInstance: %v", err)
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

        if resp.StatusCode > 299 {
                return fmt.Errorf("RetrieveInstance: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
        }

	if err = ioutil.WriteFile(outputFile, respBytes, 0644); err != nil {
		return err
	}

	fmt.Fprintf(w, "DICOM instance retrieved and downloaded to file: %v\n", outputFile)

	return nil
}

// [END healthcare_dicomweb_retrieve_instance]