package usecases

import (
	"context"
	"strings"

	"shortener/internal/entities"
	"shortener/internal/repositories"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

// Link usecase
type Link interface {
	Add(ctx context.Context, link string) (short string, err error)
	Get(ctx context.Context, link *entities.Link) error
}

// Link interactor
type link struct {
	repo repositories.Link
}

// NewLink is link constructor
func NewLink(r repositories.Link) Link {
	return &link{repo: r}
}

func (l *link) Add(ctx context.Context, link string) (string, error) {
	short, err := shortid.Generate()
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	if err := l.repo.Insert(ctx, entities.Link{Main: link, Shorted: short}); err != nil {
		return "", err
	}

	return short, nil
}

func (l *link) Get(ctx context.Context, link *entities.Link) error {
	if err := l.repo.Select(ctx, link); err != nil {
		return err
	}

	if !strings.HasPrefix(link.Main, "http") {
		link.Main = "http://" + link.Main
	}

	return nil
}
