package main

import (
	_ "embed"
	_ "image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/drcom.png
var iconPng []byte

func showDialogWindow(application fyne.App, message string, onClosed func()) {
	resource := fyne.NewStaticResource("icon.png", iconPng)
	window := application.NewWindow("Drcom 自动认证")
	window.SetIcon(resource)

	image := canvas.NewImageFromResource(resource)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(64, 64))

	text := widget.NewLabel(message)
	button := widget.NewButton("确定", func() {
		window.Close()
		onClosed()
	})

	window.SetContent(container.NewHBox(
		layout.NewSpacer(),
		image,
		layout.NewSpacer(),
		container.NewVBox(
			layout.NewSpacer(),
			text,
			layout.NewSpacer(),
			button,
			layout.NewSpacer()),
		layout.NewSpacer()))

	window.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyReturn || keyEvent.Name == fyne.KeyEnter {
			window.Close()
			onClosed()
		}
	})

	window.RequestFocus()
	window.Resize(fyne.NewSize(360, window.Content().MinSize().Height+64))
	window.CenterOnScreen()
	window.Show()
}

func showConfigWindow(application fyne.App, onSaved func(config Config), onClosed func()) {
	resource := fyne.NewStaticResource("icon.png", iconPng)
	window := application.NewWindow("Drcom 自动认证")
	window.SetIcon(resource)

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("请输入您的账号")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入您的密码")

	form := widget.NewForm(
		widget.NewFormItem("账号", usernameEntry),
		widget.NewFormItem("密码", passwordEntry),
	)

	usernameEntry.OnSubmitted = func(text string) {
		window.Canvas().Focus(passwordEntry)
	}
	passwordEntry.OnSubmitted = func(password string) {
		window.Hide()
		config := Config{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}
		onSaved(config)
	}

	saveButton := widget.NewButton("保存", func() {
		window.Hide()
		config := Config{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}
		onSaved(config)
	})
	cancelButton := widget.NewButton("取消", func() {
		window.Hide()
		onClosed()
	})

	window.SetContent(container.NewVBox(
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			saveButton,
			layout.NewSpacer(),
			cancelButton,
			layout.NewSpacer()),
		layout.NewSpacer()))

	window.RequestFocus()
	window.Canvas().Focus(usernameEntry)
	window.Resize(fyne.NewSize(360, window.Content().MinSize().Height+32))
	window.CenterOnScreen()
	window.Show()
}

func workflow(application fyne.App) {
	if config, err := loadConfig("config.json"); err == nil {
		if err := loginCampus(config); err == nil {
			os.Exit(0)
		} else {
			showDialogWindow(application, err.Error(), func() {
				os.Exit(1)
			})
		}
	} else {
		showDialogWindow(application, err.Error(), func() {
			showConfigWindow(application,
				func(config Config) {
					if err := saveConfig("config.json", config); err == nil {
						workflow(application)
					} else {
						showDialogWindow(application, err.Error(), func() {
							os.Exit(1)
						})
					}
				},
				func() {
					os.Exit(1)
				},
			)
		})
	}
}

func main() {
	application := app.New()
	workflow(application)
	application.Run()
}
