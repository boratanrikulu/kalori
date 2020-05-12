package foods

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestCalorie(t *testing.T) {
	result, err := Calorie("Hamburger")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
