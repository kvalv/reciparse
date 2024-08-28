package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/kvalv/reciparse/internal/extractors/recipeschema"
	"github.com/kvalv/reciparse/internal/http"
)

func main() {
	// read url from command line, then print the parsed recipe and ingredients
	if os.Args[1] == "" {
		fmt.Fprintf(os.Stderr, "reciparse <url>\n")
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing url: %v\n", err)
		os.Exit(1)
	}

	pr := http.NewRetriever()
	node, err := pr.Retrieve(*u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching url: %v\n", err)
		os.Exit(1)
	}
	s := recipeschema.NewExtracter()
	recipe, err := s.Extract(*node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting ingredients: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Name: %s\n", recipe.Name)
	fmt.Printf("Image: %s\n", recipe.Image)
	fmt.Println("Ingredients:")
	for _, i := range recipe.Ingredients {
		fmt.Printf("  %s\n", i)
	}

}
