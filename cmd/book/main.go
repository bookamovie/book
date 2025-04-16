package main

import "github.com/bookamovie/book/internal/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	a.Run()
}
