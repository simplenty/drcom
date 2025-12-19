package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func loginCampus(config *Config) error {
	// 初始化客户端
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	// 构造 POST 请求
	formData := url.Values{
		"DDDDD":  {config.Username},
		"upass":  {config.Password},
		"0MKKey": {"登　录"},
	}
	// 执行网络请求
	response, err := client.PostForm("http://10.0.0.253/0.htm", formData)
	if err != nil {
		return fmt.Errorf("网络请求失败")
	}
	// 自动关闭请求体
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	// 读取请求
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败")
	}
	// 转换编码
	utf8body, err := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBytes)
	if err != nil {
		return fmt.Errorf("编码转换失败")
	}
	// 转换成 UTF-8 编码
	html := string(utf8body)
	// 返回请求状态
	if strings.Contains(html, "您已经成功登录") {
		return nil
	} else if strings.Contains(html, "userid error1") {
		return fmt.Errorf("账号错误")
	} else if strings.Contains(html, "userid error2") {
		return fmt.Errorf("密码错误")
	}
	// 无法处理
	return fmt.Errorf("其他错误")
}
