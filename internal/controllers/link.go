package controllers

import (
	"net/http"

	"shortener/internal/entities"
	"shortener/internal/presenters/dtos"
	"shortener/internal/usecases"

	"github.com/mehdiazizii/fastcontroller"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Link struct {
	fastcontroller.Controller
	service usecases.Link
}

func NewLink(c fastcontroller.Controller, s usecases.Link) *Link {
	return &Link{c, s}
}

func (l *Link) Post(ctx *fastcontroller.Context) error {
	var in dtos.LinkPost
	if err := l.DecodeJson(ctx, &in); err != nil {
		return err
	}

	short, err := l.service.Add(ctx, in.Link)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return l.ResponseWithJson(
		ctx, http.StatusOK, dtos.LinkPostResponse{Short: short},
	)
}

func (l *Link) Get(ctx *fastcontroller.Context) error {
	value := ctx.UserValue("link")
	shortLink, ok := value.(string)
	if !ok {
		return fastcontroller.ErrValidation("invalid link", errors.New("string type assertion"))
	}

	link := entities.Link{Shorted: shortLink}
	if err := l.service.Get(ctx, &link); err != nil {
		return err
	}

	ctx.Redirect(link.Main, fasthttp.StatusMovedPermanently)

	return nil
}
