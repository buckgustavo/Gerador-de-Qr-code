# Gerador de QR Code

> Serviço de API REST para geração de QR codes com persistência de dados de participantes, automação de HTTPS via Let's Encrypt e gestão de dados em conformidade com a LGPD.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?style=flat&logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/licença-MIT-green?style=flat)

---

## Visão Geral

Este serviço expõe uma API HTTP que gera QR codes sob demanda, armazena dados de participantes (nome, e-mail, redes sociais) no PostgreSQL e automatiza o provisionamento de certificados TLS via Let's Encrypt. A base de código segue uma arquitetura em camadas (`handler → service → repository → domain`) para manter as responsabilidades separadas e a lógica de negócio testável de forma independente.

Decisões de design relevantes:

- **Servidor HTTP da biblioteca padrão** — sem dependência de router externo, mantendo o binário enxuto.
- **`skip2/go-qrcode`** para codificação QR — amplamente testado em produção, sem dependência de CGO.
- **`lib/pq`** como driver PostgreSQL — estável e consolidado em serviços Go de produção.
- **Conformidade com LGPD art. 18** — endpoint `DELETE /api/dados/:email` dedicado para exercício do direito ao apagamento, protegido por header `ADMIN_KEY`.
- **HTTPS automático** — quando `DOMAIN` está configurado, o servidor sobe na porta `:443` com certificado auto-renovado e redireciona `:80` para HTTPS; sem ele, sobe em HTTP simples para desenvolvimento local.

---

## Arquitetura

```
.
├── cmd/                  # Entrypoints da aplicação (pacotes main)
├── internal/
│   ├── domain/           # Entidades e regras de negócio (sem dependências externas)
│   ├── service/          # Orquestração dos casos de uso
│   ├── repository/       # Camada de acesso a dados (PostgreSQL)
│   └── handler/          # Handlers HTTP e roteamento
├── web/                  # Assets estáticos / frontend (se aplicável)
├── .env.example          # Referência das variáveis de ambiente necessárias
├── go.mod
└── go.sum
```

---

## Requisitos

- Go 1.25+
- PostgreSQL 14+ (com SSL habilitado em produção)
- Domínio público acessível para HTTPS automático (opcional em desenvolvimento local)

---

## Como Executar

**1. Clone o repositório e instale as dependências**

```bash
git clone https://github.com/buckgustavo/Gerador-de-Qr-code.git
cd Gerador-de-Qr-code
go mod download
```

**2. Configure as variáveis de ambiente**

```bash
cp .env.example .env
# Edite o arquivo .env com suas credenciais
```

| Variável       | Obrigatória | Descrição |
|----------------|-------------|-----------|
| `DATABASE_URL` | ✅          | String de conexão PostgreSQL (`postgresql://usuario:senha@host/banco?sslmode=require`) |
| `PORT`         | ❌          | Porta HTTP (padrão: `3000`) |
| `DOMAIN`       | ❌          | Domínio público para HTTPS automático via Let's Encrypt |
| `ADMIN_KEY`    | ❌          | Chave secreta para o endpoint de exclusão de dados (LGPD) |

**3. Executar**

```bash
go run ./cmd/...
```

Para produção, compile um binário estático:

```bash
CGO_ENABLED=0 GOOS=linux go build -o qrcode-server ./cmd/...
./qrcode-server
```

---

## Referência da API

### Gerar QR Code

```
POST /api/qrcode
Content-Type: application/json

{
  "content": "https://exemplo.com.br",
  "name":    "Gustavo Buck",
  "email":   "gustavo@exemplo.com.br"
}
```

Retorna uma imagem PNG com o QR code gerado.

---

### Excluir Dados do Participante (LGPD art. 18)

```
DELETE /api/dados/:email
X-Admin-Key: <ADMIN_KEY>
```

Remove permanentemente todos os dados associados ao e-mail informado.
Este endpoint fica **desabilitado** caso `ADMIN_KEY` não esteja definido.

---

## Privacidade e Conformidade

Os dados dos participantes (nome, e-mail, redes sociais) são armazenados sob as obrigações da **LGPD (Lei nº 13.709/2018)**. Os operadores devem:

- Nunca versionar o arquivo `.env` — ele contém credenciais do banco e chaves secretas.
- Gerar o `ADMIN_KEY` com uma fonte criptograficamente segura: `openssl rand -hex 32`.
- Utilizar `sslmode=require&channel_binding=require` na `DATABASE_URL` em produção.

---

## Contribuindo

1. Faça um fork do repositório e crie uma branch de feature (`git checkout -b feat/minha-feature`).
2. Realize os commits seguindo o padrão [Conventional Commits](https://www.conventionalcommits.org/pt-br/).
3. Abra um Pull Request com uma descrição clara da mudança e sua motivação.

---

## Licença

Distribuído sob a licença MIT. Consulte o arquivo [`LICENSE`](LICENSE) para mais detalhes.
