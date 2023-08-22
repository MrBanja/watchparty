package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/mrbanja/watchparty/app"
)

func main() {
	o := app.Options{}
	if err := env.Parse(&o); err != nil {
		log.Fatal("Env error: ", err)
	}

	if err := app.Run(context.Background(), o); err != nil {
		log.Fatal("Running error: ", err)
	}
}
