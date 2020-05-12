package foods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type foodSearchResult struct {
	FoodSearchCriteria struct {
		PageNumber      int  `json:"pageNumber"`
		RequireAllWords bool `json:"requireAllWords"`
	} `json:"foodSearchCriteria"`
	TotalHits   int `json:"totalHits"`
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	Foods       []struct {
		FdcID                  int    `json:"fdcId"`
		Description            string `json:"description"`
		AdditionalDescriptions string `json:"additionalDescriptions"`
		DataType               string `json:"dataType"`
		FoodCode               string `json:"foodCode"`
		PublishedDate          string `json:"publishedDate"`
		FoodNutrients          []struct {
			NutrientID     int     `json:"nutrientId"`
			NutrientName   string  `json:"nutrientName"`
			NutrientNumber string  `json:"nutrientNumber"`
			UnitName       string  `json:"unitName"`
			Value          float64 `json:"value"`
		} `json:"foodNutrients"`
		AllHighlightFields string  `json:"allHighlightFields"`
		Score              float64 `json:"score"`
	} `json:"foods"`
}

// Calorie find food's calories and returns it as string.
func Calorie(foodName string) (string, error) {
	u, _ := url.Parse("https://api.nal.usda.gov/fdc/v1/search")
	q, _ := url.ParseQuery(u.RawQuery)

	q.Add("query", foodName)
	q.Add("dataType", "Branded")
	q.Add("sortBy", "fdcId")
	q.Add("sortOrder", "desc")
	q.Add("pageSize", "1")
	u.RawQuery = q.Encode()

	// Sets the request.
	req, _ := http.NewRequest("GET", fmt.Sprint(u), nil)

	// Sends the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	foodSearchResult := new(foodSearchResult)
	json.Unmarshal(body, &foodSearchResult)

	if len(foodSearchResult.Foods) == 0 {
		return "", fmt.Errorf("We can not find the calorie for the food.")
	}

	food := foodSearchResult.Foods[0]

	energy := ""
	for _, nutrient := range food.FoodNutrients {
		if nutrient.NutrientName == "Energy" {
			energy = fmt.Sprint(nutrient.Value, " kcal")
			break
		}
	}

	if energy == "" {
		return "", fmt.Errorf("We can not find the calorie for the food.")
	}

	return energy, nil
}
