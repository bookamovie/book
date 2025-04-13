package main

import "github.com/xoticdsign/bookamovie/internal/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
