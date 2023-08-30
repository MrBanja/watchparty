package ssr

import (
	"html/template"

	"github.com/mrbanja/watchparty/party"

	"go.uber.org/zap"
)

type Service struct {
	publicAddress string
	logger        *zap.Logger

	tmpl  *template.Template
	party *party.Party
}

func New(p *party.Party, publicAddr string, logger *zap.Logger) *Service {
	return &Service{
		publicAddress: publicAddr,
		logger:        logger,
		party:         p,
	}
}

func (s *Service) MustBuildTemplate() {
	s.tmpl = template.Must(template.ParseFiles("./static/status_badge.html"))
}

type m map[string]interface{}

const (
	statusBadgeTmpl string = "status_badge.html"
)
