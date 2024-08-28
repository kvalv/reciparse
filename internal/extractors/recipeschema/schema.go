package recipeschema

import (
	"encoding/json"
	"errors"

	"github.com/PuerkitoBio/goquery"
	reciparse "github.com/kvalv/reciparse/internal"
)

// https://developers.google.com/search/docs/appearance/structured-data/recipe#structured-data-type-definitions
type RecipeSchema struct {
	// required
	Images []string `json:"image"`
	// required
	Name string `json:"name"`

	// The spec says it's an array of strings, but some pages have one long string with commas
	RecipeIngredient []string `json:"recipeIngredient"`

	// should be "Recipe", otherwise it's not a recipe and we ignore it
	Type string `json:"@type"`
}

var (
	ErrStructuredDataNotFound  = errors.New("no ld+json nodes found")
	ErrStructuredDataAmbiguous = errors.New("found multiple ld+json nodes")
)

func extractSchema(p reciparse.Page) (*RecipeSchema, error) {
	var result []RecipeSchema
	visit := func(s RecipeSchema) {
		if s.Type == "Recipe" {
			result = append(result, s)
		}
	}

	goquery.NewDocumentFromNode(p.Contents).Find(`[type="application/ld+json"]`).Each(func(i int, s *goquery.Selection) {
		scriptContents := []byte(s.Text())
		var parsedSingle RecipeSchema
		if err := json.Unmarshal(scriptContents, &parsedSingle); err == nil {
			visit(parsedSingle)
		}

		var parsedMultiple []RecipeSchema
		if err := json.Unmarshal(scriptContents, &parsedMultiple); err == nil {
			for _, recipe := range parsedMultiple {
				visit(recipe)
			}
		}
	})
	if len(result) == 0 {
		return nil, ErrStructuredDataNotFound
	}
	if len(result) > 1 {
		return nil, ErrStructuredDataAmbiguous
	}
	return &result[0], nil
}
