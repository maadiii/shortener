package repositories

import (
	"context"
	"testing"

	"shortener/config"
	"shortener/internal/entities"

	"github.com/stretchr/testify/assert"
)

func TestInsertLink(t *testing.T) {
	cfg := config.Config()
	cfg.DbSession.DBName = cfg.DbSession.TestDBName
	s, err := newTestSession(cfg.DbSession)
	assert.Nil(t, err)
	defer s.Close()
	link := NewLink(s)

	t.Run("When shorted greater than 9 chars", func(t *testing.T) {
		l := entities.Link{Main: "www.test.com", Shorted: "aaaaaaaaaa"}
		err := link.Insert(context.Background(), l)
		assert.Contains(t, err.Error(), "value too long")
	})

	t.Run("When shorted is empty", func(t *testing.T) {
		l := entities.Link{Main: "www.test.com"}
		err := link.Insert(context.Background(), l)
		assert.Contains(t, err.Error(), "non_empty_shorted")
	})

	t.Run("When main link is empty", func(t *testing.T) {
		l := entities.Link{Shorted: "aaaaaaaaa"}
		err := link.Insert(context.Background(), l)
		assert.Contains(t, err.Error(), "non_empty_main")
	})

	t.Run("When OK", func(t *testing.T) {
		l := entities.Link{Main: "www.test.com", Shorted: "aaaaaaaaa"}
		err := link.Insert(context.Background(), l)
		assert.Nil(t, err)
	})
}
