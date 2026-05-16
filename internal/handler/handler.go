package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gerador-qrcode/internal/domain"
	"gerador-qrcode/internal/repository"
	"gerador-qrcode/internal/service"
)

type Handler struct {
	gerador   *service.GeradorQR
	repo      *repository.CrachaRepo
	adminKey  string
}

func Novo(gerador *service.GeradorQR, repo *repository.CrachaRepo, adminKey string) *Handler {
	return &Handler{gerador: gerador, repo: repo, adminKey: adminKey}
}

type requisicaoGerar struct {
	Nome              string `json:"nome"`
	Email             string `json:"email"`
	LinkedIn          string `json:"linkedin"`
	Twitter           string `json:"twitter"`
	Instagram         string `json:"instagram"`
	GitHub            string `json:"github"`
	LinkPersonalizado string `json:"link_personalizado"`
	Consentimento     bool   `json:"consentimento"`
}

type respostaGerar struct {
	Erro              string `json:"erro,omitempty"`
	QRCode            string `json:"qrcode,omitempty"`
	Nome              string `json:"nome,omitempty"`
	Email             string `json:"email,omitempty"`
	LinkedIn          string `json:"linkedin,omitempty"`
	Twitter           string `json:"twitter,omitempty"`
	Instagram         string `json:"instagram,omitempty"`
	GitHub            string `json:"github,omitempty"`
	LinkPersonalizado string `json:"link_personalizado,omitempty"`
}

func (h *Handler) Gerar(w http.ResponseWriter, r *http.Request) {
	var req requisicaoGerar
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escreverJSON(w, http.StatusBadRequest, respostaGerar{Erro: "corpo da requisicao invalido"})
		return
	}

	if !req.Consentimento {
		escreverJSON(w, http.StatusUnprocessableEntity, respostaGerar{Erro: "consentimento e obrigatorio (LGPD art. 7, I)"})
		return
	}

	c := &domain.Cracha{
		Nome:              strings.TrimSpace(req.Nome),
		Email:             strings.TrimSpace(req.Email),
		LinkedIn:          extrairUsuario(req.LinkedIn),
		Twitter:           extrairUsuario(req.Twitter),
		Instagram:         extrairUsuario(req.Instagram),
		GitHub:            extrairUsuario(req.GitHub),
		LinkPersonalizado: strings.TrimSpace(req.LinkPersonalizado),
	}

	if err := c.Validar(); err != nil {
		escreverJSON(w, http.StatusUnprocessableEntity, respostaGerar{Erro: err.Error()})
		return
	}

	qr, err := h.gerador.GerarBase64(c)
	if err != nil {
		escreverJSON(w, http.StatusInternalServerError, respostaGerar{Erro: "falha ao gerar QR code"})
		return
	}

	if h.repo != nil {
		if err := h.repo.Salvar(c); err != nil {
			log.Printf("aviso: falha ao salvar cracha no banco: %v", err)
		}
	}

	escreverJSON(w, http.StatusOK, respostaGerar{
		QRCode:            qr,
		Nome:              c.Nome,
		Email:             c.Email,
		LinkedIn:          c.LinkedIn,
		Twitter:           c.Twitter,
		Instagram:         c.Instagram,
		GitHub:            c.GitHub,
		LinkPersonalizado: c.LinkPersonalizado,
	})
}

func extrairUsuario(val string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return ""
	}
	val = strings.TrimPrefix(val, "@")
	if i := strings.LastIndex(val, "/"); i >= 0 {
		val = val[i+1:]
	}
	return strings.TrimSpace(val)
}

func (h *Handler) DeletarDados(w http.ResponseWriter, r *http.Request) {
	if h.adminKey == "" {
		escreverJSON(w, http.StatusForbidden, respostaGerar{Erro: "endpoint desabilitado"})
		return
	}
	chave := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if chave != h.adminKey {
		escreverJSON(w, http.StatusUnauthorized, respostaGerar{Erro: "nao autorizado"})
		return
	}
	email := strings.TrimSpace(r.PathValue("email"))
	if email == "" {
		escreverJSON(w, http.StatusBadRequest, respostaGerar{Erro: "email nao informado"})
		return
	}
	if h.repo == nil {
		escreverJSON(w, http.StatusServiceUnavailable, respostaGerar{Erro: "banco de dados nao configurado"})
		return
	}
	n, err := h.repo.DeletarPorEmail(email)
	if err != nil {
		escreverJSON(w, http.StatusInternalServerError, respostaGerar{Erro: "falha ao deletar dados"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"mensagem": "nenhum registro encontrado para este email"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "dados removidos com sucesso"})
}

func escreverJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
