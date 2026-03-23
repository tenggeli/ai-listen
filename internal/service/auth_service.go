package service

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ai-listen/internal/dto"
	"ai-listen/internal/model"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/repository"
	"github.com/redis/go-redis/v9"
)

var chinaMobileRegexp = regexp.MustCompile(`^1\d{10}$`)

type AuthService struct {
	userRepo *repository.UserRepository
	rdb      *redis.Client
	codeTTL  time.Duration
	tokenTTL time.Duration
	appEnv   string
}

func NewAuthService(
	userRepo *repository.UserRepository,
	rdb *redis.Client,
	codeTTL time.Duration,
	tokenTTL time.Duration,
	appEnv string,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		rdb:      rdb,
		codeTTL:  codeTTL,
		tokenTTL: tokenTTL,
		appEnv:   appEnv,
	}
}

func (s *AuthService) SendCode(ctx context.Context, req dto.SendCodeRequest) (*dto.SendCodeResponse, error) {
	if !chinaMobileRegexp.MatchString(req.Mobile) {
		return nil, apperror.BadRequest("手机号格式错误")
	}

	code, err := randomDigits(6)
	if err != nil {
		return nil, apperror.Internal("发送验证码失败", err)
	}

	if err := s.rdb.Set(ctx, codeKey(req.Mobile), code, s.codeTTL).Err(); err != nil {
		return nil, apperror.Internal("发送验证码失败", err)
	}

	resp := &dto.SendCodeResponse{}
	if strings.ToLower(s.appEnv) != "prod" {
		resp.DebugCode = code
	}
	return resp, nil
}

func (s *AuthService) MobileLogin(ctx context.Context, req dto.MobileLoginRequest) (*dto.MobileLoginResponse, error) {
	if !chinaMobileRegexp.MatchString(req.Mobile) {
		return nil, apperror.BadRequest("手机号格式错误")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, apperror.BadRequest("验证码不能为空")
	}

	storedCode, err := s.rdb.Get(ctx, codeKey(req.Mobile)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, apperror.BadRequest("验证码不存在或已过期")
		}
		return nil, apperror.Internal("登录失败", err)
	}
	if strings.TrimSpace(storedCode) != strings.TrimSpace(req.Code) {
		return nil, apperror.BadRequest("验证码错误")
	}
	_ = s.rdb.Del(ctx, codeKey(req.Mobile)).Err()

	now := time.Now()
	user, err := s.userRepo.GetByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, apperror.Internal("登录失败", err)
	}

	if user == nil {
		nickname := fmt.Sprintf("用户%s", req.Mobile[len(req.Mobile)-4:])
		user = &model.User{
			Mobile:      req.Mobile,
			Nickname:    nickname,
			UserStatus:  1,
			LastLoginAt: &now,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, apperror.Internal("登录失败", err)
		}
	} else {
		if err := s.userRepo.UpdateLastLoginAt(ctx, user.ID, now); err != nil {
			return nil, apperror.Internal("登录失败", err)
		}
	}

	token, err := randomToken(32)
	if err != nil {
		return nil, apperror.Internal("登录失败", err)
	}
	if err := s.rdb.Set(ctx, tokenKey(token), strconv.FormatUint(user.ID, 10), s.tokenTTL).Err(); err != nil {
		return nil, apperror.Internal("登录失败", err)
	}

	return &dto.MobileLoginResponse{
		Token:     token,
		ExpiresAt: now.Add(s.tokenTTL).Unix(),
		UserID:    user.ID,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return apperror.BadRequest("token 不能为空")
	}
	if err := s.rdb.Del(ctx, tokenKey(token)).Err(); err != nil {
		return apperror.Internal("退出登录失败", err)
	}
	return nil
}

func codeKey(mobile string) string {
	return fmt.Sprintf("auth:code:%s", mobile)
}

func tokenKey(token string) string {
	return fmt.Sprintf("auth:token:%s", token)
}

func randomDigits(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := crand.Int(crand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}
	return string(result), nil
}

func randomToken(bytesLen int) (string, error) {
	if bytesLen <= 0 {
		return "", fmt.Errorf("invalid bytes length")
	}
	buf := make([]byte, bytesLen)
	if _, err := crand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
