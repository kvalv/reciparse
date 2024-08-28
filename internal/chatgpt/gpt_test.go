package chatgpt

import (
	"os"
	"testing"

	reciparse "github.com/kvalv/reciparse/internal"
	"github.com/stretchr/testify/assert"
)

func TestGPTParser(t *testing.T) {
	token := os.Getenv("OPENAI_KEY")
	p := NewChatGPTParser(token, WithSeed(42))
	ingredients, err := p.InterpretMany([]string{
		" 1 boks chilibønner (á 390 g) ",
		" 1 boks hermetiske kikerter (á 290 g) ",
		" 1 glass salsa (á 230 g) ",
		"4 stk. egg",
		"1 stk. vårløk",
		"2 ss hakket frisk koriander",
		"1 stk. avokado",
		"1 dl lettrømme",
		"4 stk. hvetetortilla ",
	})
	if err != nil {
		t.Fatalf("ParseIngredients: %v", err)
	}

	expected := []reciparse.StructuredIngredient{
		{Name: "chilibønner", Quantity: 1, Unit: reciparse.Piece},
		{Name: "hermetiske kikerter", Quantity: 1, Unit: reciparse.Piece},
		{Name: "salsa", Quantity: 1, Unit: reciparse.Piece},
		{Name: "egg", Quantity: 4, Unit: reciparse.Piece},
		{Name: "vårløk", Quantity: 1, Unit: reciparse.Piece},
		{Name: "frisk koriander", Quantity: 2, Unit: reciparse.Tablespoon},
		{Name: "avokado", Quantity: 1, Unit: reciparse.Piece},
		{Name: "lettrømme", Quantity: 1, Unit: reciparse.Decilitre},
		{Name: "hvetetortilla", Quantity: 4, Unit: reciparse.Piece},
	}
	assert.Equal(t, expected, ingredients)
}
