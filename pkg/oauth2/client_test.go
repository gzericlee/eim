package oauth2

import (
	"fmt"
	"testing"
)

var code string
var client Client
var token Token
var tokenInfo TokenInfo
var userInfo UserInfo

func init() {
	client, _ = NewClient(Config{Version: V5,
		Endpoint:          "http://10.200.21.162/iamsso",
		ClientId:          "zysS5H3v38ZSnrWbieJIov.nFSRsx",
		ClientSecret:      "bqbwS3uCbaQ3LI5iRV",
		RedirectUri:       "",
		LogoutRedirectUri: "",
		Debug:             true,
	})
}

func testClient_GetAccessTokenByUserPassword(t *testing.T) {
	var err error
	token, err = client.GetAccessTokenByUserPassword("admin@default", "Admin@125")
	if err != nil {
		t.Error("GetAccessTokenByUserPassword error:", err)
		return
	}
	fmt.Println("GetAccessTokenByUserPassword Successful:", token)
}

func testClient_GetAccessTokenByClientCredentials(t *testing.T) {
	var err error
	token, err = client.GetAccessTokenByClientCredentials()
	if err != nil {
		t.Error("GetAccessTokenByUserPassword error:", err)
		return
	}
	fmt.Println("GetAccessTokenByUserPassword Successful:", token)
}

func testClient_GetTokenInfo(t *testing.T) {
	var err error
	tokenInfo, err = client.GetTokenInfo(token["access_token"].(string))
	if err != nil {
		t.Error("GetTokenInfo error:", err)
		return
	}
	fmt.Println("GetTokenInfo Successful:", tokenInfo)
}

func testClient_RefreshAccessToken(t *testing.T) {
	var err error
	token, err = client.RefreshAccessToken(token["refresh_token"].(string))
	if err != nil {
		t.Error("RefreshAccessToken error:", err)
		return
	}
	fmt.Println("RefreshAccessToken Successful:", token)
}

func testClient_GetAccessTokenByTokenClientCredentials(t *testing.T) {
	var err error
	token, err = client.GetAccessTokenByTokenClientCredentials(token["access_token"].(string))
	if err != nil {
		t.Error("GetAccessTokenByTokenClientCredentials error:", err)
		return
	}
	fmt.Println("GetAccessTokenByTokenClientCredentials Successful:", token)
}

func testClient_GetCodeByAccessToken(t *testing.T) {
	var err error
	code, err = client.GetCodeByAccessToken(token["access_token"].(string))
	if err != nil {
		t.Error("GetCodeByAccessToken error:", err)
		return
	}
	fmt.Println("GetCodeByAccessToken Successful:", code)
}

func testClient_GetAccessTokenByCode(t *testing.T) {
	var err error
	token, err = client.GetAccessTokenByCode(code)
	if err != nil {
		t.Error("GetAccessTokenByCode error:", err)
		return
	}
	fmt.Println("GetAccessTokenByCode Successful:", token)
}

func testClient_GetUserInfo(t *testing.T) {
	var err error
	userInfo, err = client.GetUserInfo(token["access_token"].(string))
	if err != nil {
		t.Error("GetUserInfo error:", err)
		return
	}
	fmt.Println("GetUserInfo Successful:", userInfo)
}

func testClient_RevokeToken(t *testing.T) {
	err := client.RevokeToken(token["access_token"].(string))
	if err != nil {
		t.Error("RevokeToken error:", err)
		return
	}
	fmt.Println("RevokeToken Successful")
}

func testClient_GetSsoUrl(t *testing.T) {
	t.Log(client.GetSsoLoginUrl())
	t.Log(client.GetSsoLogoutUrl())
}

func TestClient_All(t *testing.T) {
	t.Run("GetAccessTokenByUserPassword", testClient_GetAccessTokenByUserPassword)
	t.Run("GetUserInfo", testClient_GetUserInfo)
	t.Run("RefreshAccessToken", testClient_RefreshAccessToken)
	t.Run("GetCodeByAccessToken", testClient_GetCodeByAccessToken)
	t.Run("GetAccessTokenByCode", testClient_GetAccessTokenByCode)
	t.Run("GetAccessTokenByTokenClientCredentials", testClient_GetAccessTokenByTokenClientCredentials)
	t.Run("GetAccessTokenByClientCredentials", testClient_GetAccessTokenByClientCredentials)
	t.Run("GetTokenInfo", testClient_GetTokenInfo)
	t.Run("GetTokenInfo", testClient_GetTokenInfo)
	t.Run("RevokeToken", testClient_RevokeToken)
	t.Run("GetSsoUrl", testClient_GetSsoUrl)
}
