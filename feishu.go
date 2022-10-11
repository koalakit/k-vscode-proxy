package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
)

var (
	ErrFeishuAuthenAccessTokenFailed = errors.New("AuthenAccessToken failed")
	ErrFeishuAuthenUserInfoFailed    = errors.New("AuthenUserInfo failed")
)

type AuthenAccessTokenResp struct {
	Body larkext.AuthenAccessTokenRespBody `json:"data"`
}

// 跳转到飞书登陆
func FeishuAuthenRedirect(w http.ResponseWriter, r *http.Request, state string) {
	redirectURL := FeishuAuthenURL(state)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// 获取飞书登陆URL
func FeishuAuthenURL(state string) string {
	return fmt.Sprintf("%s?redirect_uri=%s&app_id=%s&state=%s",
		gAppConfig.FeishuAuthenURL,
		gAppConfig.FeishuRedirectURL,
		gAppConfig.FeishuAppID,
		state)
}

// 飞书回调认证
func FeishuAuth(code string, state string) (user *UserData, err error) {
	accessTokenResp, err := AuthenAccessToken(code)
	if err != nil {
		return
	}

	err = AuthenUserInfo(accessTokenResp.AccessToken)
	if err != nil {
		return
	}

	uid := accessTokenResp.OpenID
	user, err = UserLoad(uid)
	if err != nil {
		user = UserDataNew(uid)
	} else {
		TokenDel(user.UserToken)
	}

	user.FeishuID = accessTokenResp.OpenID
	user.UserToken = TokenNew()
	user.Name = accessTokenResp.Name
	user.AvatarURL = accessTokenResp.AvatarURL
	user.Email = accessTokenResp.Email
	user.Mobile = accessTokenResp.Mobile
	user.Email = accessTokenResp.Email

	// 存储用户数据
	if err = UserSave(user); err != nil {
		return
	}

	// 存储token
	TokenSet(user.UserToken, user.UID, user.BackendAddress)

	return
}

func AuthenAccessToken(code string) (data larkext.AuthenAccessTokenRespBody, err error) {
	client := lark.NewClient(gAppConfig.FeishuAppID, gAppConfig.FeishuAppSecret)
	resp, err := client.Ext.Authen.AuthenAccessToken(context.Background(),
		larkext.NewAuthenAccessTokenReqBuilder().
			Body(larkext.NewAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeAuthorizationCode).
				Code(code).
				Build()).
			Build())
	if err != nil {
		return
	}

	if !resp.Success() {
		err = ErrFeishuAuthenAccessTokenFailed
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	respBody := AuthenAccessTokenResp{}
	if err = resp.JSONUnmarshalBody(&respBody); err != nil {
		return
	}

	data = respBody.Body

	return
}

func AuthenUserInfo(token string) (err error) {
	client := lark.NewClient(gAppConfig.FeishuAppID, gAppConfig.FeishuAppSecret)
	resp, err := client.Ext.Authen.AuthenUserInfo(context.Background(), larkcore.WithUserAccessToken(token))
	if err != nil {
		return
	}

	if !resp.Success() {
		err = ErrFeishuAuthenUserInfoFailed
		return
	}

	return
}
