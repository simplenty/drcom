package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func loginCampus(config *Config) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	urlStr := "http://10.0.0.252/cgi-bin/ace_web_auth.cgi"

	formData := url.Values{
		"web_jumpto":   {""},
		"orig_referer": {"http://www.qq.com/"},
		"username":     {config.Username},
		"userpwd":      {config.Password},
		"login_page":   {""},
		"temp_account": {"0"},
		"path":         {"-6"},
		"u_type":       {"1"},
	}

	response, err := client.PostForm(urlStr, formData)
	if err != nil {
		return fmt.Errorf("网络请求失败")
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败")
	}

	html := string(bodyBytes)
	redirect := extractRedirect(html)

	if strings.Contains(redirect, "login_online_detail") {
		return nil
	} else if strings.Contains(redirect, "reason=27") {
		return fmt.Errorf("账号或密码错误")
	}

	return fmt.Errorf("其他错误")
}

func extractRedirect(html string) string {
	start := strings.Index(html, "location")
	if start == -1 {
		return ""
	}

	firstQuote := strings.IndexAny(html[start:], `"'`)
	if firstQuote == -1 {
		return ""
	}
	firstQuote += start

	secondQuote := strings.IndexAny(html[firstQuote+1:], `"'`)
	if secondQuote == -1 {
		return ""
	}
	secondQuote += firstQuote + 1
	return html[firstQuote+1 : secondQuote]
}
