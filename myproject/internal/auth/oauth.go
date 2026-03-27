package auth

// OAuthProvider 定义第三方 OAuth 登录接口
// 生产环境集成 GitHub、Google、微信等
type OAuthProvider interface {
	GetAuthURL(state string) string
	Exchange(code string) (*OAuthUser, error)
}

type OAuthUser struct {
	ProviderID string
	Provider   string
	Email      string
	Name       string
	Avatar     string
}

// GitHubProvider 示例（需要 oauth2 库）
type GitHubProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	return "https://github.com/login/oauth/authorize?client_id=" + p.ClientID +
		"&redirect_uri=" + p.RedirectURL + "&state=" + state + "&scope=user:email"
}

func (p *GitHubProvider) Exchange(code string) (*OAuthUser, error) {
	// TODO: 使用 golang.org/x/oauth2 交换 token 并获取用户信息
	return nil, nil
}