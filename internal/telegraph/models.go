package telegraph

// Account represents a Telegraph account.
type Account struct {
	ShortName   string  `json:"short_name"`
	AuthorName  *string `json:"author_name,omitempty"`
	AuthorURL   *string `json:"author_url,omitempty"`
	AccessToken *string `json:"access_token,omitempty"`
	AuthURL     *string `json:"auth_url,omitempty"`
	PageCount   *int    `json:"page_count,omitempty"`
}

// Page represents a page on Telegraph.
type Page struct {
	Path       string  `json:"path"`
	URL        string  `json:"url"`
	Title      string  `json:"title"`
	Description *string `json:"description,omitempty"`
	AuthorName *string `json:"author_name,omitempty"`
	AuthorURL  *string `json:"author_url,omitempty"`
	ImageURL   *string `json:"image_url,omitempty"`
	Content    Nodes   `json:"content,omitempty"`
	Views      int     `json:"views"`
	CanEdit    *bool   `json:"can_edit,omitempty"`
}

// PageList represents a list of Telegraph articles belonging to an account.
type PageList struct {
	TotalCount int    `json:"total_count"`
	Pages      []Page `json:"pages"`
}

// PageViews represents the number of page views for a Telegraph article.
type PageViews struct {
	Views int `json:"views"`
}

// Node represents a DOM Node. It can be a string (text node) or a NodeElement.
type Node interface{}

// NodeElement represents a DOM element node.
type NodeElement struct {
	Tag      string            `json:"tag"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children Nodes             `json:"children,omitempty"`
}

// Nodes is an array of Node objects.
type Nodes []Node

// APIResponse represents the standard Telegraph API response.
type APIResponse struct {
	OK     bool   `json:"ok"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// CreateAccountRequest represents the request parameters for createAccount.
type CreateAccountRequest struct {
	ShortName  string  `json:"short_name"`
	AuthorName *string `json:"author_name,omitempty"`
	AuthorURL  *string `json:"author_url,omitempty"`
}

// EditAccountInfoRequest represents the request parameters for editAccountInfo.
type EditAccountInfoRequest struct {
	AccessToken string  `json:"access_token"`
	ShortName   *string `json:"short_name,omitempty"`
	AuthorName  *string `json:"author_name,omitempty"`
	AuthorURL   *string `json:"author_url,omitempty"`
}

// GetAccountInfoRequest represents the request parameters for getAccountInfo.
type GetAccountInfoRequest struct {
	AccessToken string   `json:"access_token"`
	Fields      []string `json:"fields,omitempty"`
}

// RevokeAccessTokenRequest represents the request parameters for revokeAccessToken.
type RevokeAccessTokenRequest struct {
	AccessToken string `json:"access_token"`
}

// CreatePageRequest represents the request parameters for createPage.
type CreatePageRequest struct {
	AccessToken   string  `json:"access_token"`
	Title         string  `json:"title"`
	AuthorName    *string `json:"author_name,omitempty"`
	AuthorURL     *string `json:"author_url,omitempty"`
	Content       Nodes   `json:"content"`
	ReturnContent *bool   `json:"return_content,omitempty"`
}

// EditPageRequest represents the request parameters for editPage.
type EditPageRequest struct {
	AccessToken   string  `json:"access_token"`
	Path          string  `json:"path"`
	Title         string  `json:"title"`
	Content       Nodes   `json:"content"`
	AuthorName    *string `json:"author_name,omitempty"`
	AuthorURL     *string `json:"author_url,omitempty"`
	ReturnContent *bool   `json:"return_content,omitempty"`
}

// GetPageRequest represents the request parameters for getPage.
type GetPageRequest struct {
	Path          string `json:"path"`
	ReturnContent *bool  `json:"return_content,omitempty"`
}

// GetPageListRequest represents the request parameters for getPageList.
type GetPageListRequest struct {
	AccessToken string `json:"access_token"`
	Offset      *int   `json:"offset,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

// GetViewsRequest represents the request parameters for getViews.
type GetViewsRequest struct {
	Path  string `json:"path"`
	Year  *int   `json:"year,omitempty"`
	Month *int   `json:"month,omitempty"`
	Day   *int   `json:"day,omitempty"`
	Hour  *int   `json:"hour,omitempty"`
}
