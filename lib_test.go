package reciparse_test

import (
	"net/url"
	"testing"

	"github.com/kvalv/reciparse"
)

func TestParseRecipe(t *testing.T) {
	parser := reciparse.New()
	u, err := url.Parse("https://www.godfisk.no/oppskrifter/brosme/ovnsbakt-brosme-med-gronnsaker-og-tomatsalsa/")
	if err != nil {
		t.Fatal(err)
	}
	got, err := parser.ParseRecipe(*u)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Ingredients) != 16 {
		t.Errorf("got %d ingredients, want 16", len(got.Ingredients))
	}
	if got.Name != "Ovnsbakt brosme med grønnsaker og tomatsalsa" {
		t.Errorf("got %s, want Ovnsbakt brosme med grønnsaker og tomatsalsa", got.Name)
	}
}
