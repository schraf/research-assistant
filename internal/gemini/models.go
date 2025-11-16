package gemini

type ModelIdentifier string

const (
	ModelIdentifierPro       ModelIdentifier = "gemini-pro-latest"
	ModelIdentifierFlash     ModelIdentifier = "gemini-flash-latest"
	ModelIdentifierFlashLite ModelIdentifier = "gemini-flash-lite-latest"
)

type ChatRole string

const (
	ChatRoleUser  ChatRole = "user"
	ChatRoleModel ChatRole = "model"
)

type ChatContent struct {
	Role ChatRole
	Text string
}

type ChatHistory []ChatContent
