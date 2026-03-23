package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	zh "github.com/alexferl/zerohttp"
	zcautocert "github.com/alexferl/zerohttp-contrib/extensions/autocert"
	zchttp3 "github.com/alexferl/zerohttp-contrib/extensions/http3"
	zccompress "github.com/alexferl/zerohttp-contrib/middleware/compress"
	"github.com/alexferl/zerohttp/httpx"
	"github.com/alexferl/zerohttp/middleware/compress"
	"github.com/alexferl/zerohttp/middleware/etag"
	"github.com/alexferl/zerohttp/middleware/securityheaders"
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
			zh.Config{
				Addr: "localhost:8080",
				SecurityHeaders: securityheaders.Config{
					ContentSecurityPolicy: csp,
				},
			},
		)
	} else if *localTLS {
		app = zh.New(
			zh.Config{
				Addr: "localhost:8080",
				TLS: zh.TLSConfig{
					Addr:     "localhost:8443",
					CertFile: "localhost+2.pem",
					KeyFile:  "localhost+2-key.pem",
				},
				SecurityHeaders: securityheaders.Config{
					ContentSecurityPolicy: csp,
				},
			},
		)

		h3Server := zchttp3.New(":8443", app)
		app.SetHTTP3Server(h3Server)

		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add(httpx.HeaderAltSvc, `h3=":8443"; ma=86400`)
				next.ServeHTTP(w, r)
			})
		})
	} else {
		app = zh.New(
			zh.Config{
				Addr: ":80",
				TLS: zh.TLSConfig{
					Addr: ":443",
				},
				Extensions: zh.ExtensionsConfig{
					AutocertManager: mgr,
				},
				SecurityHeaders: securityheaders.Config{
					ContentSecurityPolicy: csp,
					StrictTransportSecurity: securityheaders.StrictTransportSecurity{
						MaxAge:         31536000,
						PreloadEnabled: true,
					},
				},
			},
		)

		h3Server := zchttp3.NewWithAutocert(":443", app, mgr)
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
		compress.New(compress.Config{
			Algorithms: []compress.Algorithm{
				"br",
				"zstd",
				compress.Gzip,
				compress.Deflate,
			},
			Providers: []compress.Provider{
				zccompress.BrotliProvider{},
				zccompress.ZstdProvider{},
			},
		}),
		etag.New(),
	)

	app.Files("/public/", staticFiles, "public")

	app.GET("/", zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set(httpx.HeaderContentType, httpx.MIMETextHTMLCharset)
		return components.HomePage().Render(r.Context(), w)
	}))

	app.NotFound(zh.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set(httpx.HeaderContentType, httpx.MIMETextHTMLCharset)
		w.WriteHeader(http.StatusNotFound)
		return components.NotFoundPage().Render(r.Context(), w)
	}))

	if *local {
		log.Fatal(app.Start())
	} else if *localTLS {
		log.Fatal(app.StartTLS("localhost+2.pem", "localhost+2-key.pem"))
	}

	log.Fatal(app.StartAutoTLS())
}
