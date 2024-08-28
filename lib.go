package reciparse

import (
	"fmt"
	"net/url"

	reciparse "github.com/kvalv/reciparse/internal"
	"github.com/kvalv/reciparse/internal/extractors/recipeschema"
	"github.com/kvalv/reciparse/internal/http"
)

type parser struct {
	pr reciparse.PageRetriever
	ex reciparse.Extracter
}

func New() *parser {
	return &parser{
		pr: http.NewRetriever(),
		ex: recipeschema.NewExtracter(),
	}
}

func (p *parser) ParseRecipe(url url.URL) (*reciparse.Recipe, error) {
	page, err := p.pr.Retrieve(url)
	if err != nil {
		return nil, fmt.Errorf("reciparse: could not fetch page: %w", err)
	}
	recipe, err := p.ex.Extract(*page)
	if err != nil {
		return nil, fmt.Errorf("reciparse: could not extract: %w", err)
	}
	return recipe, nil
}
