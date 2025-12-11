package main

import (
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/dlgs"
)

func showErrorDialog(title string, message string) {
	_, err := dlgs.Error(title, "错误："+message)
	if err != nil {
		os.Exit(1)
	}
}

func openConfigGUI(path string) {
	// 创建应用窗口
	application := app.New()
	window := application.NewWindow("Drcom 自动认证")
	// 定义输入框组件
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("请输入用户名")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")
	loginURLEntry := widget.NewEntry()
	loginURLEntry.SetText("http://10.0.0.253/0.htm")
	timeoutEntry := widget.NewEntry()
	timeoutEntry.SetText("10")
	// 定义表单组合
	form := widget.NewForm(
		widget.NewFormItem("用户名", usernameEntry),
		widget.NewFormItem("密码", passwordEntry),
		widget.NewFormItem("URL", loginURLEntry),
		widget.NewFormItem("超时", timeoutEntry),
	)
	// 定义保存按钮
	saveButton := widget.NewButton("保存", func() {
		// 获取文本框中的超时秒数
		timeout := 10
		if parsed, err := strconv.Atoi(timeoutEntry.Text); err == nil {
			timeout = parsed
		}
		// 写入配置项
		config := &Config{
			Username:       usernameEntry.Text,
			Password:       passwordEntry.Text,
			LoginURL:       loginURLEntry.Text,
			TimeoutSeconds: timeout,
		}
		// 保存配置
		err := saveConfig(path, config)
		if err != nil {
			showErrorDialog("Drcom 自动认证", err.Error())
		}
		// 关闭窗口
		window.Close()
	})
	// 定义取消按钮
	cancelButton := widget.NewButton("取消", func() {
		// 关闭窗口
		window.Close()
	})
	// 定义按钮栏布局
	buttonBar := container.NewHBox(
		layout.NewSpacer(),
		saveButton,
		layout.NewSpacer(),
		cancelButton,
		layout.NewSpacer(),
	)
	// 定义总布局
	content := container.NewVBox(
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		buttonBar,
		layout.NewSpacer(),
	)
	// 初始化应用窗口
	window.SetContent(content)
	window.Resize(fyne.NewSize(420, window.Content().MinSize().Height+40))
	window.CenterOnScreen()
	window.Canvas().Focus(usernameEntry)
	window.ShowAndRun()
}
