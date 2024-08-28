package recipeschema

import (
	"errors"

	reciparse "github.com/kvalv/reciparse/internal"
)

// NewExtracter returns a reciparse.Extracter that extracts recipe data based on
// a script-tag with type="application/ld+json" in the head of the HTML document.
// See [1] and [2] for details
// [1] https://jsonld.com/recipe/
// [2] https://developers.google.com/search/docs/appearance/structured-data/intro-structured-data
func NewExtracter() reciparse.Extracter {
	return &extracter{}
}

type extracter struct{}

func (n *extracter) Extract(page reciparse.Page) (*reciparse.Recipe, error) {
	schema, err := extractSchema(page)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, r := range schema.RecipeIngredient {
		result = append(result, normalizeLine(r))
	}
	if len(schema.Images) == 0 {
		return nil, errors.New("no image found")
	}
	return &reciparse.Recipe{
		Ingredients: result,
		Name:        schema.Name,
		Image:       schema.Images[0],
	}, nil
}
