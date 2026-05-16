package main

import (
	"crypto/tls"
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gerador-qrcode/internal/handler"
	"gerador-qrcode/internal/repository"
	"gerador-qrcode/internal/service"
	web "gerador-qrcode/web"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	gerador := service.NovoGeradorQR()

	var repo *repository.CrachaRepo
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("falha ao abrir conexao com banco")
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			log.Fatal("banco de dados indisponivel — verifique DATABASE_URL")
		}
		repo = repository.Novo(db)
		if err := repo.IniciarTabela(); err != nil {
			log.Fatalf("falha ao criar tabela: %v", err)
		}
		log.Println("banco de dados conectado")
	} else {
		log.Println("aviso: DATABASE_URL nao definida, dados nao serao persistidos")
	}

	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		log.Println("aviso: ADMIN_KEY nao definida, endpoint de exclusao desabilitado")
	}

	h := handler.Novo(gerador, repo, adminKey)

	webContent, err := fs.Sub(web.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	limiter := handler.NovoLimiter(5)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/gerar", limiter.Middleware(h.Gerar))
	mux.HandleFunc("DELETE /api/dados/{email}", h.DeletarDados)
	mux.Handle("/", http.FileServer(http.FS(webContent)))

	dominio := os.Getenv("DOMAIN")
	if dominio != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(dominio),
		}

		go func() {
			log.Println("redirecionador HTTP→HTTPS em :80")
			if err := http.ListenAndServe(":80", m.HTTPHandler(nil)); err != nil {
				log.Printf("aviso: redirecionador HTTP encerrado: %v", err)
			}
		}()

		srv := &http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
				MinVersion:     tls.VersionTLS12,
			},
		}

		log.Printf("servidor HTTPS em https://%s", dominio)
		log.Fatal(srv.ListenAndServeTLS("", ""))
	} else {
		porta := os.Getenv("PORT")
		if porta == "" {
			porta = "3000"
		}
		log.Printf("servidor HTTP em http://localhost:%s", porta)
		log.Fatal(http.ListenAndServe(":"+porta, mux))
	}
}
