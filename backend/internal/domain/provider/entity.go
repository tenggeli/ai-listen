package provider

import "errors"

var ErrInvalidProviderTransition = errors.New("invalid provider review transition")

const (
	ReviewStatusSubmitted          = "submitted"
	ReviewStatusUnderReview        = "under_review"
	ReviewStatusApproved           = "approved"
	ReviewStatusRejected           = "rejected"
	ReviewStatusSupplementRequired = "supplement_required"
)

type Provider struct {
	ID           string
	DisplayName  string
	CityCode     string
	Bio          string
	ReviewStatus string
}

func (p *Provider) Approve() error {
	if p.ReviewStatus == ReviewStatusApproved {
		return nil
	}
	if p.ReviewStatus == ReviewStatusRejected {
		return ErrInvalidProviderTransition
	}
	p.ReviewStatus = ReviewStatusApproved
	return nil
}

func (p *Provider) Reject() error {
	if p.ReviewStatus == ReviewStatusRejected {
		return nil
	}
	if p.ReviewStatus == ReviewStatusApproved {
		return ErrInvalidProviderTransition
	}
	p.ReviewStatus = ReviewStatusRejected
	return nil
}

func (p *Provider) RequireSupplement() error {
	if p.ReviewStatus == ReviewStatusApproved {
		return ErrInvalidProviderTransition
	}
	p.ReviewStatus = ReviewStatusSupplementRequired
	return nil
}
