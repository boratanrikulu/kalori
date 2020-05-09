package foods

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestloadModel(t *testing.T) {
	err := loadModel()
	if err != nil {
		t.Error("Could not load the model.")
	}
}

func TestRecognize(t *testing.T) {
	b := getTestingPlainImage()
	result, probability := Recognize(bytes.NewBuffer(b), "jpg")
	fmt.Printf("It might be an %v (%v)\n", result, probability)
}

func getTestingPlainImage() []byte {
	b, err := ioutil.ReadFile("lentil_soup.jpg")
	if err != nil {
		log.Fatalf("Error occur while reading test image: %v", err)
	}

	return b
}
