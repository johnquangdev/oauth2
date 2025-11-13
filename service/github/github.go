package github

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/johnquangdev/oauth2/repository/interfaces"
	utils "github.com/johnquangdev/oauth2/utils"
	"golang.org/x/oauth2"
)

type Oauth2GithubService struct {
	repo   interfaces.Repo
	config oauth2.Config
	cfg    utils.Config
}

func NewGithubOauth2Service(repo interfaces.Repo, cfg utils.Config) (*Oauth2GithubService, error) {
	config := oauth2.Config{
		ClientID:     cfg.ClientId_GitHub,
		ClientSecret: cfg.ClientSecret_GitHub,
		RedirectURL:  cfg.RedirectUrl_GitHub,
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return &Oauth2GithubService{
		repo:   repo,
		cfg:    cfg,
		config: config,
	}, nil
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (g *Oauth2GithubService) GenerateAuthURL() (*GetAuthURLResponse, error) {
	// Generate random state if not provided
	state, err := GenerateRandomState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Use oauth2 library to generate auth URL
	authURL := g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return &GetAuthURLResponse{
		Url: authURL,
	}, nil
}

func (g *Oauth2GithubService) Exchange(ctx context.Context, code string) (string, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token.AccessToken, nil
}

func (g *Oauth2GithubService) GetUserInfoGithub(ctx context.Context, access_token string) (*GetGithubUserInfoReply, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start new request, error: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request, error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		errResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read google user info during error, error: %v", err)
		}
		return nil, fmt.Errorf("response with failed, status_code: %d details: %v", resp.StatusCode, string(errResponse))
	}

	var userInfo GetGithubUserInfoReply
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode body, error: %v", err)
	}

	return &userInfo, nil
}

// GetUserEmail fetches the user's primary, verified email from GitHub
// GitHub không luôn trả về email trong /user endpoint nên cần gọi /user/emails
func (g *Oauth2GithubService) GetUserEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch emails from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResponse, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(errResponse))
	}

	var emails []EmailInfo
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails response: %w", err)
	}

	// Tìm email chính và đã xác minh
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// Nếu không có primary email, lấy email verified đầu tiên
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found for GitHub user")
}
