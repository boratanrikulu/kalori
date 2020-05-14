package foods

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// Recognize does a prediction and calorie calc for the image classification.
func Recognize(image *bytes.Buffer) (string, string, error) {
	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	location := os.Getenv("GCLOUD_LOCATION")
	modelID := os.Getenv("GCLOUD_MODEL_ID")

	ctx := context.Background()
	client, err := automl.NewPredictionClient(ctx)
	if err != nil {
		return "", "", fmt.Errorf("NewPredictionClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.PredictRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
		Payload: &automlpb.ExamplePayload{
			Payload: &automlpb.ExamplePayload_Image{
				Image: &automlpb.Image{
					Data: &automlpb.Image_ImageBytes{
						ImageBytes: image.Bytes(),
					},
				},
			},
		},
		// Params is additional domain-specific parameters.
		Params: map[string]string{
			// score_threshold is used to filter the result.
			"score_threshold": "0.6",
		},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("Predict: %v", err)
	}

	payloads := resp.GetPayload()

	if len(payloads) == 0 {
		return "", "", fmt.Errorf("There is no predicted class name")
	}

	foodName := payloads[0].GetDisplayName()
	foodName = strings.ReplaceAll(foodName, "_", " ")
	foodName = strings.Title(strings.ToLower(foodName))
	foodCalorie, _ := Calorie(foodName)

	// TODO fix this messy code lines.
	// Implement a json to keep kcals.
	if foodName == "Sarma" {
		foodCalorie = "106 kcal"
	}
	if foodName == "Cigkofte" {
		foodCalorie = "181 kcal"
	}

	if foodCalorie == "" {
		return "", "", fmt.Errorf("We can not find the calorie for the image")
	}

	return foodName, foodCalorie, nil
}
