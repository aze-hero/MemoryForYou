package service

import (
	"github.com/zhaozeguang/timecapsule-api/internal/config"
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/internal/repository"
	pkgauth "github.com/zhaozeguang/timecapsule-api/pkg/auth"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type LoginResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         interface{} `json:"user"`
}

func (s *AuthService) Login(provider, code string) (*LoginResult, error) {
	var oauthUser *pkgauth.OAuthUserInfo
	var err error

	switch provider {
	case "google":
		oauthUser, err = pkgauth.ExchangeGoogleToken(
			code,
			s.cfg.GoogleClientID,
			s.cfg.GoogleClientSecret,
			s.cfg.GoogleRedirectURL,
		)
	case "github":
		oauthUser, err = pkgauth.ExchangeGitHubToken(
			code,
			s.cfg.GitHubClientID,
			s.cfg.GitHubClientSecret,
			s.cfg.GitHubRedirectURL,
		)
	default:
		return nil, gorm.ErrInvalidData
	}
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.FindByProvider(provider, oauthUser.ProviderID)
	if err == nil {
		return s.generateTokens(existingUser)
	}

	newUser := &model.User{
		Email:      oauthUser.Email,
		Name:       oauthUser.Name,
		AvatarURL:  oauthUser.AvatarURL,
		Provider:   provider,
		ProviderID: oauthUser.ProviderID,
	}
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return s.generateTokens(newUser)
}

func (s *AuthService) generateTokens(user *model.User) (*LoginResult, error) {
	accessToken, err := pkgauth.GenerateAccessToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := pkgauth.GenerateRefreshToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWTAccessExpiry,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*LoginResult, error) {
	claims, err := pkgauth.ValidateToken(refreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, pkgauth.ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(user)
}

func (s *AuthService) GetMe(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) DevLogin(email, name string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		user = &model.User{
			Email:      email,
			Name:       name,
			Provider:   "dev",
			ProviderID: email,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	}

	return s.generateTokens(user)
}
