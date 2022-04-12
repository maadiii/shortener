package repositories

import (
	"context"

	"shortener/internal/entities"

	"github.com/jackc/pgx/v4"
	"github.com/mehdiazizii/fastcontroller"
	"github.com/pkg/errors"
)

type Link interface {
	Insert(context.Context, entities.Link) error
	Select(ctx context.Context, link *entities.Link) error
}

type link struct {
	session *DbSession
}

func NewLink(s *DbSession) Link {
	return &link{session: s}
}

func (l *link) Insert(ctx context.Context, link entities.Link) error {
	_, err := l.session.Exec(ctx, "INSERT INTO link(main, shorted) VALUES($1, $2)", link.Main, link.Shorted)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (l *link) Select(ctx context.Context, link *entities.Link) error {
	err := l.session.QueryRow(ctx, "SELECT main FROM link WHERE shorted = $1", link.Shorted).Scan(&link.Main)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fastcontroller.ErrNotFound("link", errors.Wrap(err, ""))
		}

		return errors.Wrap(err, "")
	}

	return nil
}
