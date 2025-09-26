package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"net/http"

	zh "github.com/alexferl/zerohttp"
	"github.com/alexferl/zerohttp/config"
	"github.com/alexferl/zerohttp/middleware"

	"alexferlcom/components"
)

//go:embed public
var staticFiles embed.FS

var csp = "default-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self' https://cdnjs.cloudflare.com; font-src 'self' https://cdnjs.cloudflare.com; frame-ancestors 'self'; form-action 'self';"
var hosts = []string{"alexferl.com", "www.alexferl.com"}

func main() {
	local := flag.Bool("local", false, "run locally without TLS on :8080")
	flag.Parse()

	manager := zh.NewAutocertManager("/var/cache/certs", hosts...)

	var app *zh.Server

	if *local {
		app = zh.New(
			config.WithAddr("localhost:8080"),
			config.WithSecurityHeadersOptions(
				config.WithSecurityHeadersCSP(csp),
			),
		)
	} else {
		app = zh.New(
			config.WithAddr(":80"),
			config.WithTLSAddr(":443"),
			config.WithAutocertManager(manager),
			config.WithSecurityHeadersOptions(
				config.WithSecurityHeadersCSP(csp),
				config.WithSecurityHeadersHSTS(
					config.WithHSTSPreload(true),
					config.WithHSTSMaxAge(31536000),
				),
			),
		)
	}

	app.Use(middleware.Compress())

	app.Files("/public/", staticFiles, "public")

	app.GET("/", zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return components.HomePage().Render(context.Background(), w)
	}))

	app.NotFound(zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return components.NotFoundPage().Render(context.Background(), w)
	}))

	if *local {
		log.Fatal(app.Start())
	} else {
		log.Fatal(app.StartAutoTLS(hosts...))
	}
}
