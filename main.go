package main

import (
	"embed"
	"log"

	zh "github.com/alexferl/zerohttp"
	"github.com/alexferl/zerohttp/config"
	"github.com/alexferl/zerohttp/middleware"
)

//go:embed public
var files embed.FS

const csp = "default-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self' https://cdnjs.cloudflare.com; font-src 'self' https://cdnjs.cloudflare.com; frame-ancestors 'self'; form-action 'self';"

var hosts = []string{"alexferl.com", "www.alexferl.com"}

func main() {
	manager := zh.NewAutocertManager(
		"/var/cache/certs",
		hosts...,
	)

	app := zh.New(
		config.WithAddr(":80"),
		config.WithTLSAddr(":443"),
		config.WithAutocertManager(manager),
		config.WithSecurityHeadersOptions(
			config.WithSecurityHeadersCSP(csp),
		),
	)

	app.Use(middleware.Compress())

	app.Static(files, "public")

	log.Fatal(app.StartAutoTLS(hosts...))
}
