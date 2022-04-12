package dtos

import (
	"github.com/mehdiazizii/fastcontroller"
	"github.com/pkg/errors"
)

type LinkPost struct {
	Link string `json:"link"`
}

func (lp *LinkPost) Validate() error {
	if lp.Link == "" {
		return fastcontroller.ErrValidation("link can not be empty", errors.New("link Post request was empty"))
	}

	return nil
}

type LinkPostResponse struct {
	Short string `json:"shortLink"`
}

type LinkGet struct {
	Link string `json:"link"`
}
