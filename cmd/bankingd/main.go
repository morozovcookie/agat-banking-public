package main

import (
	"log"
	stdhttp "net/http"

	"github.com/morozovcookie/agat-banking/zap"
	uberzap "go.uber.org/zap"
)

func main() {
	_, err := uberzap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
}

func makeHTTPHandler(handler stdhttp.Handler, logger *uberzap.Logger) stdhttp.Handler {
	return zap.NewHTTPHandler(handler, logger)
}
