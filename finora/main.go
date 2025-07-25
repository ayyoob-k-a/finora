package main

import (
	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/di"
)

func main() {
	// Enable detailed logging

	config := configs.GetConfig()
	//set up otp

	err := di.InitDI(config)
	if err != nil {
		panic("failed to initialize dependency injection: " + err.Error())
	}

	// Test email with better HTML template

}
