package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"
)

type ProxyApp struct {
}

func (app *ProxyApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Log(r.URL.Path)

	isWebsocket := r.Header.Get("Upgrade") == "websocket"

	if !isWebsocket {
		if r.URL.Path == "/api/auth/feishu" {
			var query = r.URL.Query()
			var code = query.Get("code")
			var state = query.Get("state")
			if len(state) <= 0 {
				state = "/"
			}

			user, err := FeishuAuth(code, state)
			if err != nil {
				RenderErrorHtml(w, fmt.Sprintf(`[ERROR] %v`, err))
				return
			}

			TokenCookie(w, user.UserToken)
			http.Redirect(w, r, state, http.StatusTemporaryRedirect)

			return
		}

		if r.URL.Path == "/login" {
			query := r.URL.Query()
			redirectURL := query.Get("redirect")
			if len(redirectURL) <= 0 {
				redirectURL = "/"
			}

			RenderLoginHtml(w, "飞书登陆", FeishuAuthenURL(redirectURL))
			return
		}
	}

	LogDebug("验证cookie")

	// var err error
	var token string

	// 获取token
	{
		cookie, err := r.Cookie(gAppConfig.Cookie)
		if err != nil || cookie == nil || len(cookie.Value) <= 0 {
			http.Redirect(w, r, fmt.Sprintf("/login?redirect=%s", r.RequestURI), http.StatusTemporaryRedirect)
			return
		}

		token = cookie.Value
	}

	// 验证cookie
	// token, err := r.Cookie(gAppConfig.Cookie)
	// if err != nil || token == nil || len(token.Value) <= 0 {
	// 	// FeishuAuthenRedirect(w, r, r.RequestURI)
	// 	http.Redirect(w, r, fmt.Sprintf("/login?redirect=%s", r.RequestURI), http.StatusTemporaryRedirect)
	// 	return
	// }

	uid, backendURL, ok := TokenGet(token)
	if !ok {
		if isWebsocket {
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/login?redirect=%s", r.RequestURI), http.StatusTemporaryRedirect)
		return
	}

	if !isWebsocket {
		// 管理员操作
		if r.URL.Path == "/admin/execute" {
			query := r.URL.Query()
			command := query.Get("command")

			Log("[ADMIN] ", command)

			// 默认id为1是管理员
			if uid != 1 {
				RenderErrorHtml(w, fmt.Sprintf(`[ERROR] %v`, "没有设置操作码或者操作码不对"))
				return
			}

			// 更新token
			if command == "token-reload" {
				account := query.Get("account")
				if len(account) <= 0 {
					RenderErrorHtml(w, fmt.Sprintf(`[ERROR] %v`, "没有设置账号"))
					return
				}

				user, err := UserLoad(account)
				if err != nil {
					RenderErrorHtml(w, fmt.Sprintf(`[ERROR] %v`, err))
					return
				}

				err = TokenSet(user.UserToken, user.ID, user.BackendAddress)
				if err != nil {
					RenderErrorHtml(w, fmt.Sprintf(`[ERROR] %v`, err))
					return
				}

				RenderErrorHtml(w, fmt.Sprint(`[INFO] `, user.ID, user.BackendAddress))
				return
			}
		}
	}

	// 后端vscode未配置
	if len(backendURL) <= 0 {
		RenderLoginHtml(w, "后端VSCODE未配置, 重新登陆", r.RequestURI)
		return
	}

	// user, err := UserLoad(uid)
	// if err != nil {
	// 	if isWebsocket {
	// 		return
	// 	}

	// 	RenderErrorHtml(w, fmt.Sprintf("用户数据加载失败, 请联系管理员: %v", err))
	// 	return
	// }

	// // 检测token
	// if user.UserToken != token {
	// 	if isWebsocket {
	// 		return
	// 	}

	// 	RenderErrorHtml(w, "登陆过期，请重新登陆")
	// 	LogDebug("登陆过期，请重新登陆")
	// 	return
	// }

	if len(backendURL) <= 0 {
		RenderErrorHtml(w, "开发环境尚未配置, 请联系管理员")
		LogDebug("开发环境尚未配置, 请联系管理员")
		return
	}

	// 反向代理到vscode server
	targetURL, _ := url.Parse(backendURL)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		reply := fmt.Sprintf("[HTTP] proxy error: %v", err)
		Log(reply)
		rw.Write([]byte(reply))
	}

	proxy.ServeHTTP(w, r)
}

func RenderLoginHtml(w http.ResponseWriter, label, redirectURL string) {
	// t, _ := template.ParseFiles("./html/login.html")
	t, _ := template.New("login").Parse(htmlLogin)
	t.Execute(w, map[string]string{"Label": label, "URL": redirectURL})
}

func RenderErrorHtml(w http.ResponseWriter, message string) {
	// t, _ := template.ParseFiles("./html/error.html")
	t, _ := template.New("error").Parse(htmlError)
	t.Execute(w, map[string]string{"Message": message})
}
