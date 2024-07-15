package oauth2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gzericlee/eim/pkg/httputil"
)

type ClientV3 struct {
	clientConfig *Config
}

func (its *ClientV3) setConfig(config *Config) {
	its.clientConfig = config
}

func (its *ClientV3) GetAccessTokenByUserPassword(loginId, password string) (Token, error) {
	postUrl := fmt.Sprintf("%v/oauth2/token?grant_type=password&username=%v&password=%v", its.clientConfig.Endpoint, loginId, url.QueryEscape(password))

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
	}

	return httputil.DoRequest[Token](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
}

func (its *ClientV3) GetAccessTokenByClientCredentials() (Token, error) {
	postUrl := fmt.Sprintf("%v/oauth2/token?grant_type=client_credentials", its.clientConfig.Endpoint)

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
	}

	return httputil.DoRequest[Token](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
}

func (its *ClientV3) GetAccessTokenByTokenClientCredentials(accessToken string) (Token, error) {
	postUrl := fmt.Sprintf("%v/oauth2/token?grant_type=token_client_credentials&access_token=%v", its.clientConfig.Endpoint, accessToken)

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
	}

	return httputil.DoRequest[Token](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
}

func (its *ClientV3) GetAccessTokenByCode(code string) (Token, error) {
	postUrl := fmt.Sprintf("%v/oauth2/token?grant_type=authorization_code&code=%s", its.clientConfig.Endpoint, code)

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
	}

	return httputil.DoRequest[Token](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
}

func (its *ClientV3) GetTokenInfo(accessToken string) (TokenInfo, error) {
	var err error
	auth := identityManager.Get(accessToken)

	if auth.TokenInfo != nil && len(auth.TokenInfo["user_id"].(string)) > 0 && (time.Now().Unix() < auth.ExpiresTime.Unix()) {
		return auth.TokenInfo, nil

	} else {
		postUrl := fmt.Sprintf("%v/oauth2/tokeninfo?access_token=%v", its.clientConfig.Endpoint, accessToken)

		headers := http.Header{
			"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
		}

		auth.TokenInfo, err = httputil.DoRequest[TokenInfo](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
		if err != nil {
			return nil, err
		}

		auth.AccessToken = accessToken
		auth.ExpiresTime = time.Now().Add(time.Second * time.Duration(auth.TokenInfo["expires_in"].(float64)))

		identityManager.Set(accessToken, auth)
	}

	return auth.TokenInfo, nil
}

func (its *ClientV3) GetUserInfo(accessToken string) (UserInfo, error) {
	var err error
	auth := identityManager.Get(accessToken)

	if auth.UserInfo != nil && len(auth.UserInfo["sub"].(string)) > 0 {
		return auth.UserInfo, nil

	} else {
		postUrl := fmt.Sprintf("%v/oauth2/userinfo", its.clientConfig.Endpoint)

		headers := http.Header{
			"Authorization": []string{"Bearer " + accessToken},
		}

		auth.UserInfo, err = httputil.DoRequest[UserInfo](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
		if err != nil {
			return nil, err
		}

		identityManager.Set(accessToken, auth)
	}

	return auth.UserInfo, nil
}

func (its *ClientV3) GetCodeByAccessToken(accessToken string) (string, error) {
	_ = accessToken
	return "", ErrUnsupported
}

func (its *ClientV3) RefreshAccessToken(accessToken string) (Token, error) {
	postUrl := fmt.Sprintf("%v/oauth2/token?grant_type=refresh_token&refresh_token=%v", its.clientConfig.Endpoint, accessToken)

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
	}

	return httputil.DoRequest[Token](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
}

func (its *ClientV3) GetSsoLoginUrl() string {
	params := map[string]interface{}{
		"client_id":     its.clientConfig.ClientId,
		"redirect_uri":  its.clientConfig.RedirectUri,
		"response_type": "code id_token",
		"scope":         "openid",
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, fmt.Sprintf("%v", value))
	}

	return fmt.Sprintf("%s/oauth2/authorize?%s", its.clientConfig.Endpoint, values.Encode())
}

func (its *ClientV3) GetSsoLogoutUrl() string {
	params := map[string]interface{}{
		"post_logout_redirect_uri": its.clientConfig.LogoutRedirectUri,
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, fmt.Sprintf("%v", value))
	}

	return fmt.Sprintf("%s/oauth2/logout?%s", its.clientConfig.Endpoint, values.Encode())
}

func (its *ClientV3) RevokeToken(accessToken string) error {
	identityManager.Delete(accessToken)
	return nil
}
