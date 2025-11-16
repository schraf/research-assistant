package telegraph

import "context"

// Client defines the interface for Telegraph API operations.
type Client interface {
	// CreateAccount creates a new Telegraph account.
	// Returns an Account object with the regular fields and an additional access_token field.
	CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error)

	// EditAccountInfo updates information about a Telegraph account.
	// Pass only the parameters that you want to edit.
	EditAccountInfo(ctx context.Context, req EditAccountInfoRequest) (*Account, error)

	// GetAccountInfo gets information about a Telegraph account.
	GetAccountInfo(ctx context.Context, req GetAccountInfoRequest) (*Account, error)

	// RevokeAccessToken revokes access_token and generates a new one.
	// Returns an Account object with new access_token and auth_url fields.
	RevokeAccessToken(ctx context.Context, req RevokeAccessTokenRequest) (*Account, error)

	// CreatePage creates a new Telegraph page.
	CreatePage(ctx context.Context, req CreatePageRequest) (*Page, error)

	// EditPage edits an existing Telegraph page.
	EditPage(ctx context.Context, req EditPageRequest) (*Page, error)

	// GetPage gets a Telegraph page.
	GetPage(ctx context.Context, req GetPageRequest) (*Page, error)

	// GetPageList gets a list of pages belonging to a Telegraph account.
	// Returns a PageList object, sorted by most recently created pages first.
	GetPageList(ctx context.Context, req GetPageListRequest) (*PageList, error)

	// GetViews gets the number of views for a Telegraph article.
	// By default, the total number of page views will be returned.
	GetViews(ctx context.Context, req GetViewsRequest) (*PageViews, error)
}
