package http

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	reciparse "github.com/kvalv/reciparse/internal"
	"golang.org/x/net/html"
)

type cache struct {
	dir string
}

func (c *cache) Get(url url.URL) (*html.Node, error) {
	sh := sha256.New()
	sh.Write([]byte(url.String()))
	if _, err := os.Stat(fmt.Sprintf("%s/%x", c.dir, sh.Sum(nil))); os.IsNotExist(err) {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("%s/%x", c.dir, sh.Sum(nil)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return html.Parse(f)
}
func (c *cache) Write(url url.URL, contents *html.Node) error {
	if _, err := os.Stat(c.dir); os.IsNotExist(err) {
		os.Mkdir(c.dir, 0755)
	}
	sh := sha256.New()
	sh.Write([]byte(url.String()))
	f, err := os.Create(fmt.Sprintf("%s/%x", c.dir, sh.Sum(nil)))
	if err != nil {
		return err
	}
	defer f.Close()
	return html.Render(f, contents)
}

type HTTPFetcher struct {
	client *http.Client
	cache  *cache
}

// Fetch implements Fetcher.
func (h *HTTPFetcher) Retrieve(url url.URL) (*reciparse.Page, error) {
	if h.cache != nil {
		node, err := h.cache.Get(url)
		if err == nil {
			return &reciparse.Page{Contents: node}, nil
		}
	}
	node, err := h.retrieve(url)
	if err != nil {
		return nil, err
	}
	if h.cache != nil {
		err = h.cache.Write(url, node.Contents)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (h *HTTPFetcher) retrieve(url url.URL) (*reciparse.Page, error) {
	res, err := h.client.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d", res.StatusCode)
	}
	contents, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	return &reciparse.Page{Contents: contents}, nil
}

type Option func(*HTTPFetcher)

func WithCache(dir string) Option {
	return func(h *HTTPFetcher) {
		h.cache = &cache{dir: dir}
	}
}

func NewRetriever(opts ...Option) reciparse.PageRetriever {
	client := &HTTPFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}
