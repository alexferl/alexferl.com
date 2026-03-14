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

type autocertManager struct {
	*autocert.Manager
}

func (a *autocertManager) Hostnames() []string {
	return hosts
}

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
	return nil
}

func (h *http3AutocertServer) ListenAndServeTLSWithAutocert(manager config.AutocertManager) error {
	tlsConfig := &tls.Config{
		GetCertificate: manager.GetCertificate,
		NextProtos:     []string{"h3"},
	}
	h.server.TLSConfig = tlsConfig

	err := h.server.ListenAndServe()
	if err != nil {
		log.Printf("[ERROR] HTTP/3 server failed: %v", err)
	}
	return err
}

func main() {
	local := flag.Bool("local", false, "run locally without TLS on :8080")
	localTLS := flag.Bool("local-tls", false, "run locally with TLS on :8443 (requires localhost+2.pem and localhost+2-key.pem)")
	flag.Parse()

	manager := &autocertManager{
		Manager: &autocert.Manager{
			Cache:      autocert.DirCache("/var/cache/certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hosts...),
		},
	}

	var app *zh.Server

	if *local {
		app = zh.New(
			config.Config{
				Addr: "localhost:8080",
				SecurityHeaders: config.SecurityHeadersConfig{
					ContentSecurityPolicy: csp,
				},
			},
		)
	} else if *localTLS {
		app = zh.New(
			config.Config{
				Addr: "localhost:8080",
				TLS: config.TLSConfig{
					Addr:     "localhost:8443",
					CertFile: "localhost+2.pem",
					KeyFile:  "localhost+2-key.pem",
				},
				SecurityHeaders: config.SecurityHeadersConfig{
					ContentSecurityPolicy: csp,
				},
			},
		)

		h3Server := &http3.Server{
			Addr:    ":8443",
			Handler: app,
		}
		app.SetHTTP3Server(h3Server)

		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Alt-Svc", `h3=":8443"; ma=86400`)
				next.ServeHTTP(w, r)
			})
		})
	} else {
		app = zh.New(
			config.Config{
				Addr: ":80",
				TLS: config.TLSConfig{
					Addr: ":443",
				},
				Extensions: config.ExtensionsConfig{
					AutocertManager: manager,
				},
				SecurityHeaders: config.SecurityHeadersConfig{
					ContentSecurityPolicy: csp,
					StrictTransportSecurity: config.StrictTransportSecurity{
						MaxAge:         31536000,
						PreloadEnabled: true,
					},
				},
			},
		)

		h3Server := &http3AutocertServer{
			server: &http3.Server{
				Addr:    ":443",
				Handler: app,
			},
		}
		app.SetHTTP3Server(h3Server)
	}

	if !*local && !*localTLS {
		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Alt-Svc", `h3=":443"; ma=86400`)
				next.ServeHTTP(w, r)
			})
		})
	}

	app.Use(
		middleware.Compress(),
		middleware.ETag(),
	)

	app.Files("/public/", staticFiles, "public")

	app.GET("/", zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return components.HomePage().Render(context.Background(), w)
	}))

	app.NotFound(zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return components.NotFoundPage().Render(context.Background(), w)
	}))

	if *local {
		log.Fatal(app.Start())
	} else if *localTLS {
		log.Println("Starting local TLS server on https://localhost:8443")
		log.Fatal(app.StartTLS("localhost+2.pem", "localhost+2-key.pem"))
	}

	log.Fatal(app.StartAutoTLS())
}
