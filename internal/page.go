package reciparse

import (
	"net/url"

	"golang.org/x/net/html"
)

type Page struct {
	Contents *html.Node
}

type PageRetriever interface {
	Retrieve(url.URL) (*Page, error)
}
