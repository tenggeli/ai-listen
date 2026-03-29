package audio

import (
	"context"
	"errors"
)

var ErrInvalidPage = errors.New("invalid sounds page")

type HomePageService interface {
	GetHomePage(ctx context.Context, userID string) (HomePage, error)
}
