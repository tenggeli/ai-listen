package audio

import (
	"context"
	"errors"
)

var ErrInvalidPage = errors.New("invalid sounds page")
var ErrInvalidCategory = errors.New("invalid sounds category")
var ErrInvalidInput = errors.New("invalid input")
var ErrSoundNotFound = errors.New("sound not found")

type TrackQuery struct {
	CategoryKey string
	PageNo      int
	PageSize    int
}

// HomePageService is kept for compatibility with older mock implementations.
type HomePageService interface {
	GetHomePage(ctx context.Context, userID string) (HomePage, error)
}

type Repository interface {
	ListCategories(ctx context.Context) ([]Category, error)
	ListTracks(ctx context.Context, query TrackQuery) ([]Track, int, error)
	ListAdminTracks(ctx context.Context, query AdminTrackQuery) ([]AdminTrack, int, error)
	GetAdminTrackByID(ctx context.Context, trackID string) (AdminTrack, error)
	CreateAdminTrack(ctx context.Context, input AdminTrack) error
	UpdateAdminTrack(ctx context.Context, input AdminTrack) error
	UpdateAdminTrackStatus(ctx context.Context, trackID string, status string) error
}

const (
	SoundStatusActive   = "active"
	SoundStatusInactive = "inactive"
)

type AdminTrack struct {
	ID            string
	CategoryKey   string
	Title         string
	PlayCountText string
	DurationText  string
	Emoji         string
	Author        string
	SortOrder     int
	Status        string
}

type AdminTrackQuery struct {
	CategoryKey string
	Status      string
	Keyword     string
	PageNo      int
	PageSize    int
}
