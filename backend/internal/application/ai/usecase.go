package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	domain "listen/backend/internal/domain/ai"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) string
}

type GetRemainingMatchInput struct {
	UserID string
}

type GetHomeOverviewInput struct {
	UserID string
}

type GetHomeOverviewOutput struct {
	Overview domain.HomeOverview
}

type GetRemainingMatchOutput struct {
	Remaining int
}

type SubmitMatchInput struct {
	UserID    string
	InputText string
}

type SubmitMatchOutput struct {
	Remaining  int
	Candidates []domain.MatchCandidate
}

type CreateSessionInput struct {
	UserID    string
	SceneType string
}

type CreateSessionOutput struct {
	SessionID string
}

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	Session domain.Session
}

type AppendMessageInput struct {
	SessionID   string
	SenderType  string
	ContentText string
}

type AppendMessageOutput struct {
	Session domain.Session
}

type GetRemainingMatchUseCase struct {
	quotaRepo domain.MatchQuotaRepository
	clock     Clock
}

type GetAiHomeUseCase struct {
	quotaRepo    domain.MatchQuotaRepository
	homeOverview domain.HomeOverviewService
	clock        Clock
}

func NewGetRemainingMatchUseCase(quotaRepo domain.MatchQuotaRepository, clock Clock) GetRemainingMatchUseCase {
	return GetRemainingMatchUseCase{quotaRepo: quotaRepo, clock: clock}
}

func NewGetAiHomeUseCase(quotaRepo domain.MatchQuotaRepository, homeOverview domain.HomeOverviewService, clock Clock) GetAiHomeUseCase {
	return GetAiHomeUseCase{quotaRepo: quotaRepo, homeOverview: homeOverview, clock: clock}
}

func (u GetRemainingMatchUseCase) Execute(ctx context.Context, input GetRemainingMatchInput) (GetRemainingMatchOutput, error) {
	if input.UserID == "" {
		return GetRemainingMatchOutput{}, domain.ErrInvalidInput
	}
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, dateKey(u.clock.Now()))
	if err != nil {
		return GetRemainingMatchOutput{}, err
	}
	return GetRemainingMatchOutput{Remaining: quota.Remaining()}, nil
}

func (u GetAiHomeUseCase) Execute(ctx context.Context, input GetHomeOverviewInput) (GetHomeOverviewOutput, error) {
	if input.UserID == "" {
		return GetHomeOverviewOutput{}, domain.ErrInvalidInput
	}

	now := u.clock.Now()
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, dateKey(now))
	if err != nil {
		return GetHomeOverviewOutput{}, err
	}

	overview, err := u.homeOverview.BuildOverview(ctx, domain.HomeOverviewServiceInput{
		UserID:    input.UserID,
		Now:       now,
		Remaining: quota.Remaining(),
	})
	if err != nil {
		return GetHomeOverviewOutput{}, err
	}
	return GetHomeOverviewOutput{Overview: overview}, nil
}

type SubmitMatchUseCase struct {
	quotaRepo    domain.MatchQuotaRepository
	matchService domain.MatchService
	clock        Clock
}

func NewSubmitMatchUseCase(quotaRepo domain.MatchQuotaRepository, matchService domain.MatchService, clock Clock) SubmitMatchUseCase {
	return SubmitMatchUseCase{quotaRepo: quotaRepo, matchService: matchService, clock: clock}
}

func (u SubmitMatchUseCase) Execute(ctx context.Context, input SubmitMatchInput) (SubmitMatchOutput, error) {
	if input.UserID == "" || input.InputText == "" {
		return SubmitMatchOutput{}, domain.ErrInvalidInput
	}

	date := dateKey(u.clock.Now())
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, date)
	if err != nil {
		return SubmitMatchOutput{}, err
	}

	if quota.Remaining() == 0 {
		return SubmitMatchOutput{}, domain.ErrDailyLimitReached
	}

	result, err := u.matchService.GenerateCandidates(ctx, input.InputText)
	if err != nil {
		return SubmitMatchOutput{}, err
	}

	if !result.IsEmpty() {
		if err := quota.Consume(); err != nil {
			return SubmitMatchOutput{}, err
		}
		if err := u.quotaRepo.Save(ctx, quota); err != nil {
			return SubmitMatchOutput{}, err
		}
	}

	return SubmitMatchOutput{Remaining: quota.Remaining(), Candidates: result.Candidates}, nil
}

type CreateAiSessionUseCase struct {
	sessionRepo domain.SessionRepository
	idGenerator IDGenerator
}

func NewCreateAiSessionUseCase(sessionRepo domain.SessionRepository, idGenerator IDGenerator) CreateAiSessionUseCase {
	return CreateAiSessionUseCase{sessionRepo: sessionRepo, idGenerator: idGenerator}
}

func (u CreateAiSessionUseCase) Execute(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error) {
	sessionID := u.idGenerator.NewID("sess")
	session, err := domain.NewSession(sessionID, input.UserID, input.SceneType)
	if err != nil {
		return CreateSessionOutput{}, err
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		return CreateSessionOutput{}, err
	}
	return CreateSessionOutput{SessionID: sessionID}, nil
}

type GetAiSessionUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewGetAiSessionUseCase(sessionRepo domain.SessionRepository) GetAiSessionUseCase {
	return GetAiSessionUseCase{sessionRepo: sessionRepo}
}

func (u GetAiSessionUseCase) Execute(ctx context.Context, input GetSessionInput) (GetSessionOutput, error) {
	session, err := u.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return GetSessionOutput{}, err
	}
	return GetSessionOutput{Session: session}, nil
}

type AppendAiMessageUseCase struct {
	sessionRepo domain.SessionRepository
	replySvc    domain.ReplyService
	clock       Clock
	replyCfg    replyConfig
	lockPool    *sessionLockPool
}

type replyConfig struct {
	timeout    time.Duration
	maxRetries int
}

type assistantReply struct {
	Content     string
	ActionCard  *domain.ActionCard
	SafetyLevel string
}

func NewAppendAiMessageUseCase(sessionRepo domain.SessionRepository, replySvc domain.ReplyService, clock Clock) AppendAiMessageUseCase {
	return AppendAiMessageUseCase{
		sessionRepo: sessionRepo,
		replySvc:    replySvc,
		clock:       clock,
		replyCfg: replyConfig{
			timeout:    800 * time.Millisecond,
			maxRetries: 2,
		},
		lockPool: newSessionLockPool(),
	}
}

func (u AppendAiMessageUseCase) Execute(ctx context.Context, input AppendMessageInput) (AppendMessageOutput, error) {
	unlock := u.lockPool.Lock(input.SessionID)
	defer unlock()

	session, err := u.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return AppendMessageOutput{}, err
	}
	if err := session.AppendMessage(input.SenderType, input.ContentText, u.clock.Now()); err != nil {
		return AppendMessageOutput{}, err
	}

	if input.SenderType == "user" {
		reply := u.buildAssistantReply(ctx, session, input.ContentText)
		if err := session.AppendMessageWithMeta("assistant", reply.Content, u.clock.Now(), reply.ActionCard, reply.SafetyLevel); err != nil {
			return AppendMessageOutput{}, err
		}
	}

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return AppendMessageOutput{}, err
	}
	return AppendMessageOutput{Session: session}, nil
}

func (u AppendAiMessageUseCase) generateReplyWithRetry(ctx context.Context, session domain.Session, userMessage string) (string, error) {
	if u.replySvc == nil {
		return "", errors.New("reply service is nil")
	}

	var lastErr error
	for attempt := 0; attempt <= u.replyCfg.maxRetries; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, u.replyCfg.timeout)
		reply, err := u.replySvc.GenerateReply(attemptCtx, session, userMessage)
		cancel()
		if err == nil {
			trimmed := strings.TrimSpace(reply)
			if trimmed != "" {
				return trimmed, nil
			}
			err = domain.ErrInvalidInput
		}
		lastErr = err
	}

	if lastErr == nil {
		lastErr = domain.ErrInvalidInput
	}
	return "", lastErr
}

func (u AppendAiMessageUseCase) buildAssistantReply(ctx context.Context, session domain.Session, userMessage string) assistantReply {
	content := strings.TrimSpace(userMessage)
	if isSensitiveMessage(content) {
		return assistantReply{
			Content: "你提到的内容可能和自我伤害风险有关，这很重要。请先保证你身边环境安全，并尽快联系身边可信任的人；如果风险正在发生，请立即联系当地紧急援助电话。",
			ActionCard: &domain.ActionCard{
				Action:      "go_sound",
				Title:       "先做 3 分钟稳定情绪",
				Description: "进入声音页，先让呼吸慢下来，再决定下一步。",
				Route:       "/sound",
				ButtonText:  "去声音",
			},
			SafetyLevel: "high",
		}
	}

	replyText, err := u.generateReplyWithRetry(ctx, session, userMessage)
	if err != nil {
		replyText = buildFallbackReply(userMessage)
	}

	return assistantReply{
		Content:     replyText,
		ActionCard:  pickActionCard(content),
		SafetyLevel: "normal",
	}
}

func isSensitiveMessage(content string) bool {
	if content == "" {
		return false
	}

	sensitiveKeywords := []string{
		"自杀",
		"轻生",
		"不想活",
		"活不下去",
		"结束生命",
		"伤害自己",
		"自残",
	}
	negations := []string{"不", "没", "不会", "不是", "并非", "没有", "不想"}

	for _, keyword := range sensitiveKeywords {
		index := strings.Index(content, keyword)
		if index < 0 {
			continue
		}
		prefix := content[:index]
		prefixRunes := []rune(prefix)
		windowStart := len(prefixRunes) - 6
		if windowStart < 0 {
			windowStart = 0
		}
		nearbyPrefix := string(prefixRunes[windowStart:])
		negated := false
		for _, neg := range negations {
			if strings.Contains(nearbyPrefix, neg) {
				negated = true
				break
			}
		}
		if !negated {
			return true
		}
	}
	return false
}

func pickActionCard(content string) *domain.ActionCard {
	type actionRule struct {
		keywords []string
		card     domain.ActionCard
	}

	rules := []actionRule{
		{
			keywords: []string{"睡不着", "失眠", "放松", "白噪音", "冥想", "声音"},
			card: domain.ActionCard{
				Action:      "go_sound",
				Title:       "试试声音放松",
				Description: "白噪音和呼吸引导可以先帮你稳住状态。",
				Route:       "/sound",
				ButtonText:  "去声音",
			},
		},
		{
			keywords: []string{"下单", "服务", "咨询", "陪伴服务"},
			card: domain.ActionCard{
				Action:      "go_service",
				Title:       "看看可选服务",
				Description: "可以先浏览服务方，再决定是否发起邀约。",
				Route:       "/services",
				ButtonText:  "去服务",
			},
		},
		{
			keywords: []string{"找人", "匹配", "推荐", "搭子", "陪聊", "孤单"},
			card: domain.ActionCard{
				Action:      "go_match",
				Title:       "为你匹配陪伴对象",
				Description: "回到首页快速匹配，先看最适合你的 3 位推荐。",
				Route:       "/home",
				ButtonText:  "去匹配",
			},
		},
	}

	for _, rule := range rules {
		for _, keyword := range rule.keywords {
			if strings.Contains(content, keyword) {
				card := rule.card
				return &card
			}
		}
	}
	return nil
}

type sessionLockPool struct {
	mu    sync.Mutex
	locks map[string]*sessionRefLock
}

type sessionRefLock struct {
	mu   sync.Mutex
	refs int
}

func newSessionLockPool() *sessionLockPool {
	return &sessionLockPool{locks: make(map[string]*sessionRefLock)}
}

func (p *sessionLockPool) Lock(sessionID string) func() {
	key := sessionID
	if key == "" {
		key = "_empty_session"
	}

	p.mu.Lock()
	ref, ok := p.locks[key]
	if !ok {
		ref = &sessionRefLock{}
		p.locks[key] = ref
	}
	ref.refs++
	p.mu.Unlock()

	ref.mu.Lock()
	return func() {
		ref.mu.Unlock()
		p.mu.Lock()
		ref.refs--
		if ref.refs == 0 {
			delete(p.locks, key)
		}
		p.mu.Unlock()
	}
}

type MockReplyService struct{}

func NewMockReplyService() MockReplyService {
	return MockReplyService{}
}

func (MockReplyService) GenerateReply(_ context.Context, _ domain.Session, userMessage string) (string, error) {
	content := strings.TrimSpace(userMessage)
	switch {
	case strings.Contains(content, "压力"), strings.Contains(content, "累"):
		return "听起来你已经扛了很久。先不用急着解决全部问题，我们可以先把最难受的那一刻说出来。", nil
	case strings.Contains(content, "睡"):
		return "睡不着往往不是你不够努力，而是情绪还没被放下。你愿意和我说说，今晚最卡在心里的事吗？", nil
	case strings.Contains(content, "孤单"), strings.Contains(content, "说说话"):
		return "那种想找个人说话的瞬间很真实，我在。你想到哪就说到哪，我们慢慢来。", nil
	default:
		return "我听见你了。现在不用组织得很完整，先从你最想说的一句开始就好。", nil
	}
}

func buildFallbackReply(_ string) string {
	return "我在这里，已经收到你的消息。刚刚有点忙乱，但不会离开。你愿意再多说一点吗？"
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type TimestampIDGenerator struct {
	clock Clock
}

func NewTimestampIDGenerator(clock Clock) TimestampIDGenerator {
	return TimestampIDGenerator{clock: clock}
}

func (g TimestampIDGenerator) NewID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, g.clock.Now().UnixNano())
}

func dateKey(now time.Time) string {
	return now.Format("2006-01-02")
}

func greetingPeriod(now time.Time) string {
	period := "晚上好"
	switch hour := now.Hour(); {
	case hour < 11:
		period = "早上好"
	case hour < 14:
		period = "中午好"
	case hour < 18:
		period = "下午好"
	}

	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return period + " · " + weekdays[int(now.Weekday())]
}
