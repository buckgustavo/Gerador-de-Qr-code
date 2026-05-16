package domain

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Cracha struct {
	Nome              string
	Email             string
	LinkedIn          string
	Twitter           string
	Instagram         string
	GitHub            string
	LinkPersonalizado string
}

func (c *Cracha) Validar() error {
	if strings.TrimSpace(c.Nome) == "" {
		return errors.New("nome e sobrenome sao obrigatorios")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("email e obrigatorio")
	}
	if len(strings.Fields(c.Nome)) < 2 {
		return errors.New("informe nome e sobrenome")
	}
	if !emailRegex.MatchString(c.Email) {
		return errors.New("endereco de email invalido")
	}
	return nil
}

func (c *Cracha) ConteudoQR() string {
	if c.LinkPersonalizado != "" {
		return c.LinkPersonalizado
	}

	var urls []string
	if c.LinkedIn != "" {
		urls = append(urls, "https://linkedin.com/in/"+c.LinkedIn)
	}
	if c.Twitter != "" {
		urls = append(urls, "https://twitter.com/"+c.Twitter)
	}
	if c.Instagram != "" {
		urls = append(urls, "https://instagram.com/"+c.Instagram)
	}
	if c.GitHub != "" {
		urls = append(urls, "https://github.com/"+c.GitHub)
	}

	if len(urls) == 1 {
		return urls[0]
	}
	if len(urls) > 1 {
		return strings.Join(urls, "\n")
	}

	return "mailto:" + c.Email
}
