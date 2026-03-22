package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"net/http"

	zh "github.com/alexferl/zerohttp"
	zcautocert "github.com/alexferl/zerohttp-contrib/extensions/autocert"
	"github.com/alexferl/zerohttp-contrib/extensions/http3"
	"github.com/alexferl/zerohttp-contrib/middleware/compress"
	"github.com/alexferl/zerohttp/config"
	"github.com/alexferl/zerohttp/httpx"
	"github.com/alexferl/zerohttp/middleware"
	"golang.org/x/crypto/acme/autocert"

	"alexferlcom/components"
)

//go:embed public
var staticFiles embed.FS

var csp = "default-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self' https://cdnjs.cloudflare.com; font-src 'self' https://cdnjs.cloudflare.com; frame-ancestors 'self'; form-action 'self';"
var hosts = []string{"alexferl.com", "www.alexferl.com"}

func main() {
	local := flag.Bool("local", false, "run locally without TLS on :8080")
	localTLS := flag.Bool("local-tls", false, "run locally with TLS on :8443 (requires localhost+2.pem and localhost+2-key.pem)")
	flag.Parse()

	mgr := zcautocert.New(
		autocert.DirCache("/var/cache/certs"),
		hosts,
	)

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

		h3Server := http3.New(":8443", app)
		app.SetHTTP3Server(h3Server)

		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(httpx.HeaderAltSvc, `h3=":8443"; ma=86400`)
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
					AutocertManager: mgr,
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

		h3Server := http3.NewWithAutocert(":443", app, mgr)
		app.SetHTTP3Server(h3Server)
	}

	if !*local && !*localTLS {
		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(httpx.HeaderAltSvc, `h3=":443"; ma=86400`)
				next.ServeHTTP(w, r)
			})
		})
	}

	app.Use(
		middleware.Compress(config.CompressConfig{
			Algorithms: []config.CompressionAlgorithm{
				"br",
				"zstd",
				config.Gzip,
				config.Deflate,
			},
			Providers: []config.CompressionProvider{
				compress.BrotliProvider{},
				compress.ZstdProvider{},
			},
		}),
		middleware.ETag(),
	)

	app.Files("/public/", staticFiles, "public")

	app.GET("/", zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set(httpx.HeaderContentType, httpx.MIMETextHTMLCharset)
		return components.HomePage().Render(context.Background(), w)
	}))

	app.NotFound(zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusNotFound)
		return components.NotFoundPage().Render(context.Background(), w)
	}))

	if *local {
		log.Fatal(app.Start())
	} else if *localTLS {
		log.Fatal(app.StartTLS("localhost+2.pem", "localhost+2-key.pem"))
	}

	log.Fatal(app.StartAutoTLS())
}
