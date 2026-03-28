package identity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput            = errors.New("invalid input")
	ErrInvalidCredential       = errors.New("invalid credential")
	ErrGenderChangeNeedConfirm = errors.New("gender change requires confirmation")
	ErrUserNotFound            = errors.New("user not found")
)

const (
	LoginChannelSMS    = "sms"
	LoginChannelWechat = "wechat"
)

type UserAccount struct {
	UserID               string
	Phone                string
	WechatOpenID         string
	RegisterSource       string
	Status               string
	DisplayName          string
	AvatarURL            string
	AgeRange             string
	City                 string
	Bio                  string
	Gender               string
	ProfileCompleted     bool
	MBTI                 string
	InterestTags         []string
	PersonalitySkipped   bool
	PersonalityCompleted bool
}

type User struct {
	ID                   string
	Phone                string
	Nickname             string
	AvatarURL            string
	RegisterSource       string
	Status               string
	ProfileCompleted     bool
	PersonalityCompleted bool
}

type UserProfile struct {
	UserID   string
	Gender   string
	AgeRange string
	City     string
	Bio      string
}

type UserPersonalityProfile struct {
	UserID       string
	MBTI         string
	InterestTags []string
	Skipped      bool
}

type UserIdentity struct {
	UserID           string
	LoginChannel     string
	AccessToken      string
	RefreshToken     string
	ExpiresInSeconds int64
	DisplayName      string
	AvatarURL        string
	IsNewUser        bool
	ProfileCompleted bool
}

func NewUserAccount(userID string, phone string, wechatOpenID string) (UserAccount, error) {
	if userID == "" {
		return UserAccount{}, ErrInvalidInput
	}
	suffix := userID
	if len(userID) > 4 {
		suffix = userID[len(userID)-4:]
	}
	return UserAccount{
		UserID:           userID,
		Phone:            phone,
		WechatOpenID:     wechatOpenID,
		RegisterSource:   sourceByAccount(phone, wechatOpenID),
		Status:           "active",
		DisplayName:      "user_" + suffix,
		AvatarURL:        "",
		ProfileCompleted: false,
	}, nil
}

func (a UserAccount) ToUser() User {
	return User{
		ID:                   a.UserID,
		Phone:                a.Phone,
		Nickname:             a.DisplayName,
		AvatarURL:            a.AvatarURL,
		RegisterSource:       a.RegisterSource,
		Status:               a.Status,
		ProfileCompleted:     a.ProfileCompleted,
		PersonalityCompleted: a.PersonalityCompleted,
	}
}

func (a UserAccount) ToProfile() UserProfile {
	return UserProfile{
		UserID:   a.UserID,
		Gender:   a.Gender,
		AgeRange: a.AgeRange,
		City:     a.City,
		Bio:      a.Bio,
	}
}

func (a UserAccount) ToPersonalityProfile() UserPersonalityProfile {
	tags := make([]string, 0, len(a.InterestTags))
	tags = append(tags, a.InterestTags...)
	return UserPersonalityProfile{
		UserID:       a.UserID,
		MBTI:         a.MBTI,
		InterestTags: tags,
		Skipped:      a.PersonalitySkipped,
	}
}

func (a *UserAccount) CompleteBasicProfile(nickname string, avatarURL string, gender string, ageRange string, city string, bio string, genderChangeConfirmed bool) error {
	nickname = strings.TrimSpace(nickname)
	gender = strings.TrimSpace(strings.ToLower(gender))
	ageRange = strings.TrimSpace(ageRange)
	city = strings.TrimSpace(city)
	bio = strings.TrimSpace(bio)
	avatarURL = strings.TrimSpace(avatarURL)

	if nickname == "" || gender == "" {
		return ErrInvalidInput
	}
	if err := a.ConfirmGenderOnce(gender, genderChangeConfirmed); err != nil {
		return err
	}

	a.DisplayName = nickname
	if avatarURL != "" {
		a.AvatarURL = avatarURL
	}
	a.AgeRange = ageRange
	a.City = city
	a.Bio = bio
	a.ProfileCompleted = a.IsBasicProfileCompleted()
	return nil
}

func (a *UserAccount) ConfirmGenderOnce(newGender string, genderChangeConfirmed bool) error {
	if newGender == "" {
		return ErrInvalidInput
	}
	if a.Gender == "" || a.Gender == newGender {
		a.Gender = newGender
		return nil
	}
	if !genderChangeConfirmed {
		return ErrGenderChangeNeedConfirm
	}
	a.Gender = newGender
	return nil
}

func (a UserAccount) IsBasicProfileCompleted() bool {
	return strings.TrimSpace(a.DisplayName) != "" &&
		strings.TrimSpace(a.Gender) != "" &&
		strings.TrimSpace(a.AgeRange) != "" &&
		strings.TrimSpace(a.City) != ""
}

func (a *UserAccount) UpdatePersonality(mbti string, interestTags []string) error {
	cleanMBTI := strings.ToUpper(strings.TrimSpace(mbti))
	cleanTags := normalizeTags(interestTags)

	if cleanMBTI == "" && len(cleanTags) == 0 {
		return ErrInvalidInput
	}

	if cleanMBTI != "" {
		if !isValidMBTI(cleanMBTI) {
			return ErrInvalidInput
		}
		a.MBTI = cleanMBTI
	}
	if len(cleanTags) > 0 {
		a.InterestTags = cleanTags
	}
	a.PersonalitySkipped = false
	a.PersonalityCompleted = a.IsPersonalityCompleted()
	return nil
}

func (a *UserAccount) SkipForNow() {
	a.PersonalitySkipped = true
	a.PersonalityCompleted = true
}

func (a UserAccount) IsPersonalityCompleted() bool {
	return a.PersonalitySkipped || strings.TrimSpace(a.MBTI) != "" || len(a.InterestTags) > 0
}

func BuildUserIdentity(account UserAccount, loginChannel string, now time.Time, isNewUser bool) UserIdentity {
	return UserIdentity{
		UserID:           account.UserID,
		LoginChannel:     loginChannel,
		AccessToken:      "mock_at_" + account.UserID + "_" + now.Format("20060102150405"),
		RefreshToken:     "mock_rt_" + account.UserID + "_" + now.Format("20060102150405"),
		ExpiresInSeconds: 7200,
		DisplayName:      account.DisplayName,
		AvatarURL:        account.AvatarURL,
		IsNewUser:        isNewUser,
		ProfileCompleted: account.ProfileCompleted,
	}
}

func sourceByAccount(phone string, openID string) string {
	if strings.TrimSpace(phone) != "" {
		return LoginChannelSMS
	}
	if strings.TrimSpace(openID) != "" {
		return LoginChannelWechat
	}
	return "unknown"
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		clean := strings.TrimSpace(tag)
		if clean == "" {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		result = append(result, clean)
	}
	return result
}

func isValidMBTI(mbti string) bool {
	if len(mbti) != 4 {
		return false
	}
	if !strings.ContainsRune("IE", rune(mbti[0])) {
		return false
	}
	if !strings.ContainsRune("NS", rune(mbti[1])) {
		return false
	}
	if !strings.ContainsRune("TF", rune(mbti[2])) {
		return false
	}
	if !strings.ContainsRune("JP", rune(mbti[3])) {
		return false
	}
	return true
}
