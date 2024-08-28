package reciparse

type Unit string

var (
	Piece      Unit = "stk"
	Gram       Unit = "g"
	MilliLitre Unit = "ml"
	Decilitre  Unit = "dl"
	Tablespoon Unit = "ss"
)

var AllUnits = []Unit{Piece, Gram, MilliLitre, Decilitre, Tablespoon}
