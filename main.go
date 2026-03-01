package main

import (
	"context"
	"crypto/tls"
	"embed"
	"flag"
	"log"
	"net/http"

	zh "github.com/alexferl/zerohttp"
	"github.com/alexferl/zerohttp/config"
	"github.com/alexferl/zerohttp/middleware"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/crypto/acme/autocert"

	"alexferlcom/components"
)

//go:embed public
var staticFiles embed.FS

var csp = "default-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self' https://cdnjs.cloudflare.com; font-src 'self' https://cdnjs.cloudflare.com; frame-ancestors 'self'; form-action 'self';"
var hosts = []string{"alexferl.com", "www.alexferl.com"}

type http3AutocertServer struct {
	server *http3.Server
}

func (h *http3AutocertServer) ListenAndServeTLS(certFile, keyFile string) error {
	return h.server.ListenAndServeTLS(certFile, keyFile)
}

func (h *http3AutocertServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *http3AutocertServer) Close() error {
	return h.server.Close()
}

func (h *http3AutocertServer) ListenAndServeTLSWithAutocert(manager config.AutocertManager) error {
	if h.server.TLSConfig == nil {
		h.server.TLSConfig = &tls.Config{}
	}
	h.server.TLSConfig.GetCertificate = manager.GetCertificate
	return h.server.ListenAndServeTLS("", "")
}

func main() {
	local := flag.Bool("local", false, "run locally without TLS on :8080")
	flag.Parse()

	manager := &autocert.Manager{
		Cache:      autocert.DirCache("/var/cache/certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
	}

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

		// HTTP/3 disabled for now - needs custom quic-go listener setup
		// h3Server := &http3AutocertServer{
		// 	server: &http3.Server{
		// 		Addr:    ":443",
		// 		Handler: app,
		// 	},
		// }
		// app.SetHTTP3Server(h3Server)
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
	}

	log.Fatal(app.StartAutoTLS())
}
