package main

import "os"

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		showErrorDialog("Drcom 自动认证", err.Error())
		openConfigGUI("config.json")
		config, err = loadConfig("config.json")
		if err != nil {
			showErrorDialog("Drcom 自动认证", err.Error())
			os.Exit(1)
		}
	}

	err = loginCampus(config)
	if err != nil {
		showErrorDialog("Drcom 自动认证", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
