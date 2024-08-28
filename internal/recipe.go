package reciparse

type Ingredient = string

type Recipe struct {
	Ingredients []Ingredient
	Name        string
	Image       string
}

type Extracter interface {
	Extract(Page) (*Recipe, error)
}

// A StructuredIngredient is an ingredient line that is parsed, so that the Name, Unit and Quantity
// are separated out from the raw text.
type StructuredIngredient struct {
	Name     string
	Quantity float32
	Unit     Unit
}

type IngredientLineInterpreter interface {
	Interpret(Ingredient) (*StructuredIngredient, error)
	InterpretMany([]Ingredient) ([]StructuredIngredient, error)
}
