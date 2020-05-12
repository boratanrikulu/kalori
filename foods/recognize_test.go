package foods

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestRecognize(t *testing.T) {
	testImages := []string{
		"test_images/hamburger.jpg",
		"test_images/pizza.jpg",
		"test_images/lentil_soup.jpg",
	}

	for _, value := range testImages {
		b := getTestingPlainImage(value)
		result, calorie, err := Recognize(bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(result, calorie)
	}
}

func getTestingPlainImage(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error occur while reading test image: %v", err)
	}

	return b
}
