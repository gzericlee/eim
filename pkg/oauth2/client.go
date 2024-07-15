package oauth2

import (
	"fmt"
	"time"
)

type Version string

const (
	V3 Version = "3"
	V5 Version = "5"
)

var registry = make(map[Version]Client)

func init() {
	registry[V3] = &ClientV3{}
	registry[V5] = &ClientV5{}
}

type Config struct {
	Version           Version
	Endpoint          string
	ClientId          string
	ClientSecret      string
	RedirectUri       string
	LogoutRedirectUri string
	Timeout           time.Duration
	Debug             bool
}

type Client interface {
	setConfig(config *Config)

	//密码模式
	GetAccessTokenByUserPassword(loginId, password string) (Token, error)

	//客户端模式
	GetAccessTokenByClientCredentials() (Token, error)

	//内部服务模式
	GetAccessTokenByTokenClientCredentials(accessToken string) (Token, error)

	//授权码模式
	GetAccessTokenByCode(code string) (Token, error)

	//获取授权码
	GetCodeByAccessToken(accessToken string) (string, error)

	//获取用户信息
	GetUserInfo(accessToken string) (UserInfo, error)

	//获取Token信息
	GetTokenInfo(accessToken string) (TokenInfo, error)

	//刷新Token
	RefreshAccessToken(accessToken string) (Token, error)

	//获取SSO登陆地址
	GetSsoLoginUrl() string

	//获取SSO注销地址
	GetSsoLogoutUrl() string

	//注销Token
	RevokeToken(accessToken string) error
}

func NewClient(config *Config) (Client, error) {
	iamClient := registry[config.Version]
	if iamClient == nil {
		return nil, fmt.Errorf("unsupported Version %v", config.Version)
	}
	if config.Timeout == 0 {
		config.Timeout = time.Minute * 2
	}
	iamClient.setConfig(config)
	return iamClient, nil
}
