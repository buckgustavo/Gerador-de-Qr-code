package repository

import (
	"database/sql"

	"gerador-qrcode/internal/domain"
)

type CrachaRepo struct {
	db *sql.DB
}

func Novo(db *sql.DB) *CrachaRepo {
	return &CrachaRepo{db: db}
}

func (r *CrachaRepo) IniciarTabela() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS crachas (
			id                 SERIAL PRIMARY KEY,
			nome               TEXT NOT NULL,
			email              TEXT NOT NULL,
			linkedin           TEXT,
			twitter            TEXT,
			instagram          TEXT,
			github             TEXT,
			link_personalizado TEXT,
			consentimento_em   TIMESTAMPTZ NOT NULL,
			criado_em          TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (r *CrachaRepo) Salvar(c *domain.Cracha) error {
	_, err := r.db.Exec(
		`INSERT INTO crachas (nome, email, linkedin, twitter, instagram, github, link_personalizado, consentimento_em)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		c.Nome, c.Email, c.LinkedIn, c.Twitter, c.Instagram, c.GitHub, c.LinkPersonalizado,
	)
	return err
}

func (r *CrachaRepo) DeletarPorEmail(email string) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM crachas WHERE email = $1`, email)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
