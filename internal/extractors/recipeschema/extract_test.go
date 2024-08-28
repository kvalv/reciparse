package recipeschema_test

import (
	_ "embed"
	"net/url"
	"strings"
	"testing"

	"github.com/kvalv/reciparse/internal/extractors/recipeschema"
	"github.com/kvalv/reciparse/internal/http"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type testCase struct {
	Id       string   `yaml:"id"`
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	URL      string   `yaml:"url"`
	Method   string   `yaml:"method"`
	Expected []string `yaml:"expected"`
}

var (
	CACHE_DIR = "/tmp/reciparse_testdata"
)

//go:embed testdata/cases.yaml
var yamlFile string

func TestExtractors(t *testing.T) {
	var testCases []testCase
	if err := yaml.NewDecoder(strings.NewReader(yamlFile)).Decode(&testCases); err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCases {
		t.Run(tc.Id, func(t *testing.T) {
			src := recipeschema.NewExtracter()
			retr := http.NewRetriever(http.WithCache(CACHE_DIR))
			url, err := url.Parse(tc.URL)
			if err != nil {
				t.Fatal(err)
			}
			page, err := retr.Retrieve(*url)
			if err != nil {
				t.Fatal(err)
			}
			recipe, err := src.Extract(*page)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.Expected, recipe.Ingredients)

			if tc.Name != "" {
				assert.Equal(t, tc.Name, recipe.Name)
			}
			if tc.Image != "" {
				assert.Equal(t, tc.Image, recipe.Image)
			}
		})
	}
}
