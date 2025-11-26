package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/schraf/research-assistant/internal/models"
	"google.golang.org/genai"
)

type client struct {
	genaiClient *genai.Client
	system      *string
	logger      *slog.Logger
}

func NewClient(ctx context.Context) (Client, error) {
	genaiClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &client{
		genaiClient: genaiClient,
		system:      nil,
	}, nil
}

//--========================================================--
//--=== RESOURCES IMPLEMENTATION
//--========================================================--

func (c *client) Ask(ctx context.Context, mode models.ResourceMode, persona string, request string) (*string, error) {
	var model ModelIdentifier

	switch mode {
	case models.ResourceModeMinimal:
		model = ModelIdentifierFlashLite
	case models.ResourceModeBasic:
		model = ModelIdentifierFlash
	case models.ResourceModePro:
		model = ModelIdentifierPro
	}

	c.SetSystemInstruction(persona)
	return c.GenerateText(ctx, model, request)
}

func (c *client) StructuredAsk(ctx context.Context, mode models.ResourceMode, persona string, request string, schema models.Schema) (json.RawMessage, error) {
	var model ModelIdentifier

	switch mode {
	case models.ResourceModeMinimal:
		model = ModelIdentifierFlashLite
	case models.ResourceModeBasic:
		model = ModelIdentifierFlash
	case models.ResourceModePro:
		model = ModelIdentifierPro
	}

	c.SetSystemInstruction(persona)
	return c.GenerateJson(ctx, model, request, schema)
}

//--========================================================--
//--=== INTERFACE IMPLEMENTATION
//--========================================================--

func (c *client) EnableLogging(filename string) (func(), error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return func() {}, fmt.Errorf("failed to open log file: %w", err)
	}

	handler := slog.NewJSONHandler(file, nil)
	c.logger = slog.New(handler)

	return func() { file.Close() }, nil
}

func (c *client) SetSystemInstruction(instruction string) {
	if c.logger != nil {
		c.logger.Info("set_system_instruction",
			slog.String("system", instruction),
		)
	}

	c.system = &instruction
}

func (c *client) GenerateText(ctx context.Context, model ModelIdentifier, prompt string) (*string, error) {
	config := c.contentConfig()
	result, err := c.genaiClient.Models.GenerateContent(ctx, string(model), genai.Text(prompt), config)
	if err != nil {
		return nil, err
	}

	if c.logger != nil {
		c.logger.Info("generate_text",
			slog.String("prompt", prompt),
			slog.String("response", result.Text()),
		)
	}

	responseText := result.Text()
	return &responseText, nil
}

func (c *client) GenerateJson(ctx context.Context, model ModelIdentifier, prompt string, schema map[string]any) (json.RawMessage, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: schema,
	}

	if c.system != nil {
		config.SystemInstruction = genai.NewContentFromText(*c.system, genai.RoleModel)
	}

	result, err := c.genaiClient.Models.GenerateContent(ctx, string(model), genai.Text(prompt), config)
	if err != nil {
		return nil, err
	}

	responseText := result.Text()

	var responseJson json.RawMessage

	err = json.Unmarshal([]byte(responseText), &responseJson)
	if err != nil {
		return nil, err
	}

	if c.logger != nil {
		c.logger.Info("generate_json",
			slog.String("prompt", prompt),
			slog.Any("response", responseJson),
		)
	}

	return responseJson, nil
}

func (c *client) Chat(ctx context.Context, model ModelIdentifier, history ChatHistory, message string) (*string, error) {
	chatContent := []*genai.Content{}

	for _, content := range history {
		var role genai.Role

		switch content.Role {
		case ChatRoleUser:
			role = genai.RoleUser
		case ChatRoleModel:
			role = genai.RoleModel
		default:
			return nil, fmt.Errorf("unknown role in chat history: %s", string(content.Role))
		}

		chatContent = append(chatContent, genai.NewContentFromText(content.Text, role))
	}

	config := c.contentConfig()

	chat, err := c.genaiClient.Chats.Create(ctx, string(model), config, chatContent)
	if err != nil {
		return nil, err
	}

	result, err := chat.SendMessage(ctx, genai.Part{Text: message})
	if err != nil {
		return nil, err
	}

	responseText := result.Text()
	return &responseText, nil
}

func (c *client) contentConfig() *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
				URLContext:   &genai.URLContext{},
			},
		},
	}

	if c.system != nil {
		config.SystemInstruction = genai.NewContentFromText(*c.system, genai.RoleModel)
	}

	return config
}
