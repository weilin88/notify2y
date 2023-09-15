package main

import "os"

var APP_CONFIG *AppConfig

type AppConfig struct {
	Email string
}

func init() {
	APP_CONFIG = new(AppConfig)
	APP_CONFIG.Email = os.Getenv("app_email")
}
