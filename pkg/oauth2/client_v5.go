package oauth2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gzericlee/eim/pkg/httputil"
	"github.com/gzericlee/eim/pkg/maputil"
)

type ClientV5 struct {
	ClientV3
}

func (its *ClientV5) GetTokenInfo(accessToken string) (TokenInfo, error) {
	var err error
	auth := identityManager.Get(accessToken)

	if auth.TokenInfo != nil && (time.Now().Unix() < auth.ExpiresTime.Unix()) {
		return auth.TokenInfo, nil

	} else {
		postUrl := fmt.Sprintf("%v/oauth2/introspect?token=%v", its.clientConfig.Endpoint, accessToken)

		headers := http.Header{
			"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
		}

		auth.TokenInfo, err = httputil.DoRequest[TokenInfo](context.Background(), postUrl, http.MethodPost, headers, nil, its.clientConfig.Debug)
		if err != nil {
			return nil, err
		}

		auth.AccessToken = accessToken
		auth.ExpiresTime = time.Unix(int64(auth.TokenInfo["exp"].(float64)), 0)

		identityManager.Set(accessToken, auth)
	}

	return auth.TokenInfo, nil
}

func (its *ClientV5) GetCodeByAccessToken(accessToken string) (string, error) {
	postUrl := fmt.Sprintf("%v/oauth2/authzcode?client_id=%v&access_token=%v", its.clientConfig.Endpoint, its.clientConfig.ClientId, accessToken)

	result, err := httputil.DoRequest[map[string]interface{}](context.Background(), postUrl, http.MethodPost, nil, nil, its.clientConfig.Debug)
	if err != nil {
		return "", err
	}

	return maputil.GetString(result, "code", ""), err
}

func (its *ClientV5) RevokeToken(accessToken string) error {
	identityManager.Delete(accessToken)

	postUrl := fmt.Sprintf("%v/oauth2/revoke", its.clientConfig.Endpoint)

	headers := http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", its.clientConfig.ClientId, its.clientConfig.ClientSecret)))},
		"Content-Type":  []string{"application/x-www-form-urlencoded"},
	}

	params := url.Values{}
	params.Set("token", accessToken)
	params.Set("token_type_hint", "access_token")

	_, err := httputil.DoRequest[any](context.Background(), postUrl, http.MethodPost, headers, params, its.clientConfig.Debug)
	return err
}
