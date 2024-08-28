package chatgpt

import (
	"context"
	"fmt"
	"strings"
	"time"

	reciparse "github.com/kvalv/reciparse/internal"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type chatGPTParser struct {
	client         *openai.Client
	requestTimeout time.Duration
	seed           *int
}

func NewChatGPTParser(token string, opts ...gptOption) reciparse.IngredientLineInterpreter {
	p := &chatGPTParser{
		client:         openai.NewClient(token),
		requestTimeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (c *chatGPTParser) Interpret(line reciparse.Ingredient) (*reciparse.StructuredIngredient, error) {
	parsed, err := c.InterpretMany([]reciparse.Ingredient{line})
	if err != nil {
		return nil, err
	}
	return &parsed[0], nil
}

// ParseIngredients implements reciparse.Extracter.
func (c *chatGPTParser) InterpretMany(lines []reciparse.Ingredient) ([]reciparse.StructuredIngredient, error) {

	var choices []string
	for _, u := range reciparse.AllUnits {
		choices = append(choices, string(u))
	}

	schema := jsonschema.Definition{
		Type: "object",
		Properties: map[string]jsonschema.Definition{
			"ingredients": {
				Type: "array",
				Items: &jsonschema.Definition{
					Type: "object",
					Properties: map[string]jsonschema.Definition{
						"name":     {Type: "string"},
						"quantity": {Type: "number"},
						"unit":     {Type: "string", Enum: choices},
					},
					Required:             []string{"name", "quantity", "unit"},
					AdditionalProperties: false,
				},
			},
		},
		Required:             []string{"ingredients"},
		AdditionalProperties: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-2024-08-06",
		Seed:  c.seed,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Parse the ingredients from the recipe. The language may be English or Norwegian. The ingredients are going on a shopping list, so details on how it's prepared (eg 'hakket' / 'chopped') are not needed (but keep 'box' and how it's contained)",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: strings.Join(lines, "\n"),
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "ingredients",
				Schema: &schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("CreateChatCompletion: %w", err)
	}
	var result struct {
		Ingredients []struct {
			Name     string  `json:"name"`
			Quantity float32 `json:"quantity"`
			Unit     string  `json:"unit"`
		} `json:"ingredients"`
	}
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &result)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %w", err)
	}
	parsed := make([]reciparse.StructuredIngredient, len(result.Ingredients))
	for i, r := range result.Ingredients {
		parsed[i] = reciparse.StructuredIngredient{
			Name:     r.Name,
			Quantity: r.Quantity,
			Unit:     reciparse.Unit(r.Unit),
		}
	}
	return parsed, nil
}

type gptOption func(*chatGPTParser)

func WithTimeout(timeout time.Duration) gptOption {
	return func(p *chatGPTParser) {
		p.requestTimeout = timeout
	}
}
func WithSeed(seed int) gptOption {
	return func(p *chatGPTParser) {
		s := seed
		p.seed = &s
	}
}
