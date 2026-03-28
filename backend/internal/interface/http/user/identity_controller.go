package user

import (
	"encoding/json"
	"errors"
	"net/http"

	app "listen/backend/internal/application/identity"
	domain "listen/backend/internal/domain/identity"
)

type IdentityController struct {
	loginSMSUC        app.LoginBySMSUseCase
	loginWechatUC     app.LoginByWechatUseCase
	getProfileUC      app.GetUserProfileUseCase
	saveProfileUC     app.SaveUserProfileUseCase
	savePersonalityUC app.SaveUserPersonalityUseCase
	skipPersonalityUC app.SkipUserPersonalityUseCase
}

func NewIdentityController(
	loginSMSUC app.LoginBySMSUseCase,
	loginWechatUC app.LoginByWechatUseCase,
	getProfileUC app.GetUserProfileUseCase,
	saveProfileUC app.SaveUserProfileUseCase,
	savePersonalityUC app.SaveUserPersonalityUseCase,
	skipPersonalityUC app.SkipUserPersonalityUseCase,
) IdentityController {
	return IdentityController{
		loginSMSUC:        loginSMSUC,
		loginWechatUC:     loginWechatUC,
		getProfileUC:      getProfileUC,
		saveProfileUC:     saveProfileUC,
		savePersonalityUC: savePersonalityUC,
		skipPersonalityUC: skipPersonalityUC,
	}
}

func (c IdentityController) HandleSMSLogin(w http.ResponseWriter, r *http.Request) {
	var body LoginBySMSRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.loginSMSUC.Execute(r.Context(), app.LoginBySMSInput{
		Phone:             body.Phone,
		VerifyCode:        body.VerifyCode,
		AgreementAccepted: body.AgreementAccepted,
	})
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildIdentityResponse(output))
}

func (c IdentityController) HandleWechatMockLogin(w http.ResponseWriter, r *http.Request) {
	var body LoginByWechatRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.loginWechatUC.Execute(r.Context(), app.LoginByWechatInput{
		AuthCode:          body.AuthCode,
		AgreementAccepted: body.AgreementAccepted,
	})
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildIdentityResponse(output))
}

func (c IdentityController) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	output, err := c.getProfileUC.Execute(r.Context(), app.GetUserProfileInput{UserID: userID})
	if err != nil {
		writeIdentityError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildMeResponse(output))
}

func (c IdentityController) HandleSaveProfile(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)

	var body SaveProfileRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.saveProfileUC.Execute(r.Context(), app.SaveUserProfileInput{
		UserID:                userID,
		Nickname:              body.Nickname,
		AvatarURL:             body.AvatarURL,
		Gender:                body.Gender,
		AgeRange:              body.AgeRange,
		City:                  body.City,
		Bio:                   body.Bio,
		GenderChangeConfirmed: body.GenderChangeConfirmed,
	})
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildMeResponse(app.GetUserProfileOutput{
		User:        output.User,
		Profile:     output.Profile,
		Personality: output.Personality,
	}))
}

func (c IdentityController) HandleSavePersonality(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)

	var body SavePersonalityRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.savePersonalityUC.Execute(r.Context(), app.SaveUserPersonalityInput{
		UserID:       userID,
		MBTI:         body.MBTI,
		InterestTags: body.InterestTags,
	})
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildMeResponse(app.GetUserProfileOutput{
		User:        output.User,
		Profile:     output.Profile,
		Personality: output.Personality,
	}))
}

func (c IdentityController) HandleSkipPersonality(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	output, err := c.skipPersonalityUC.Execute(r.Context(), app.SkipUserPersonalityInput{
		UserID: userID,
	})
	if err != nil {
		writeIdentityError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildMeResponse(app.GetUserProfileOutput{
		User:        output.User,
		Profile:     output.Profile,
		Personality: output.Personality,
	}))
}

func buildIdentityResponse(output app.LoginOutput) UserIdentityResponseDTO {
	return UserIdentityResponseDTO{
		UserID:           output.Identity.UserID,
		LoginChannel:     output.Identity.LoginChannel,
		AccessToken:      output.Identity.AccessToken,
		RefreshToken:     output.Identity.RefreshToken,
		ExpiresInSeconds: output.Identity.ExpiresInSeconds,
		DisplayName:      output.Identity.DisplayName,
		AvatarURL:        output.Identity.AvatarURL,
		IsNewUser:        output.Identity.IsNewUser,
		ProfileCompleted: output.Identity.ProfileCompleted,
	}
}

func writeIdentityError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredential):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrGenderChangeNeedConfirm):
		writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}

func buildMeResponse(output app.GetUserProfileOutput) UserMeResponseDTO {
	return UserMeResponseDTO{
		UserID:               output.User.ID,
		Phone:                output.User.Phone,
		Nickname:             output.User.Nickname,
		AvatarURL:            output.User.AvatarURL,
		RegisterSource:       output.User.RegisterSource,
		Status:               output.User.Status,
		ProfileCompleted:     output.User.ProfileCompleted,
		Gender:               output.Profile.Gender,
		AgeRange:             output.Profile.AgeRange,
		City:                 output.Profile.City,
		Bio:                  output.Profile.Bio,
		PersonalityCompleted: output.User.PersonalityCompleted,
		MBTI:                 output.Personality.MBTI,
		InterestTags:         output.Personality.InterestTags,
		PersonalitySkipped:   output.Personality.Skipped,
	}
}
