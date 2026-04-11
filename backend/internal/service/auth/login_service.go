package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
)

type DefaultAdminSeed struct {
	Enabled     bool
	Username    string
	Password    string
	DisplayName string
	Email       string
}

type LoginInput struct {
	Username string
	Password string
}

type RefreshInput struct {
	RefreshToken string
}

type AuthUser struct {
	ID            string
	Username      string
	DisplayName   string
	PlatformRoles []string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         AuthUser
}

type LoginService struct {
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
	roleRepo     *repository.PlatformRoleRepository
	passwordSvc  *PasswordService
	tokenSvc     *TokenService
	defaultAdmin DefaultAdminSeed
}

func NewLoginService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	roleRepo *repository.PlatformRoleRepository,
	passwordSvc *PasswordService,
	tokenSvc *TokenService,
	defaultAdmin DefaultAdminSeed,
) *LoginService {
	return &LoginService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		roleRepo:     roleRepo,
		passwordSvc:  passwordSvc,
		tokenSvc:     tokenSvc,
		defaultAdmin: defaultAdmin,
	}
}

func (s *LoginService) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	username := strings.TrimSpace(in.Username)
	password := in.Password
	if username == "" || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status != domain.UserStatusActive {
		return nil, ErrUserDisabled
	}
	if !s.passwordSvc.Verify(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueAndPersistSession(ctx, user, 0)
}

func (s *LoginService) Refresh(ctx context.Context, in RefreshInput) (*LoginResult, error) {
	refreshToken := strings.TrimSpace(in.RefreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidRefresh
	}

	claims, err := s.tokenSvc.ParseAndValidate(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefresh
	}

	activeSession, err := s.sessionRepo.GetActiveByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefresh
		}
		return nil, err
	}

	if activeSession.UserID == 0 || activeSession.UserID != claims.UserID {
		return nil, ErrInvalidRefresh
	}

	user, err := s.userRepo.GetByID(ctx, activeSession.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefresh
		}
		return nil, err
	}
	if user.Status != domain.UserStatusActive {
		return nil, ErrUserDisabled
	}

	return s.issueAndPersistSession(ctx, user, activeSession.ID)
}

func (s *LoginService) EnsureDefaultAdmin(ctx context.Context) error {
	if s.roleRepo != nil {
		if err := s.roleRepo.EnsureDefaults(ctx); err != nil {
			return err
		}
	}

	if !s.defaultAdmin.Enabled {
		return nil
	}
	if strings.TrimSpace(s.defaultAdmin.Username) == "" || strings.TrimSpace(s.defaultAdmin.Password) == "" {
		return nil
	}

	if existing, err := s.userRepo.GetByUsername(ctx, s.defaultAdmin.Username); err == nil {
		if s.roleRepo != nil && existing != nil && existing.ID != 0 {
			return s.roleRepo.EnsureUserRoleByName(ctx, existing.ID, "platform-admin")
		}
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := s.passwordSvc.Hash(s.defaultAdmin.Password)
	if err != nil {
		return err
	}

	// For development convenience only: bootstrap a default admin account when absent.
	admin := &domain.User{
		Username:     s.defaultAdmin.Username,
		DisplayName:  s.defaultAdmin.DisplayName,
		Email:        s.defaultAdmin.Email,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
	}
	err = s.userRepo.Create(ctx, admin)
	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		if s.roleRepo == nil {
			return nil
		}
		existing, getErr := s.userRepo.GetByUsername(ctx, s.defaultAdmin.Username)
		if getErr != nil {
			return nil
		}
		return s.roleRepo.EnsureUserRoleByName(ctx, existing.ID, "platform-admin")
	}
	if err != nil {
		return err
	}
	if s.roleRepo != nil && admin.ID != 0 {
		return s.roleRepo.EnsureUserRoleByName(ctx, admin.ID, "platform-admin")
	}
	return nil
}

func (s *LoginService) issueAndPersistSession(
	ctx context.Context,
	user *domain.User,
	revokeSessionID uint64,
) (*LoginResult, error) {
	if user == nil || user.ID == 0 {
		return nil, gorm.ErrInvalidData
	}

	accessToken, err := s.tokenSvc.IssueAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenSvc.IssueRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	accessClaims, err := s.tokenSvc.ParseAndValidate(accessToken)
	if err != nil {
		return nil, err
	}
	refreshClaims, err := s.tokenSvc.ParseAndValidate(refreshToken)
	if err != nil {
		return nil, err
	}

	expiresAt, expiresIn := claimsExpiry(refreshClaims), claimsTTLSeconds(accessClaims)
	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}

	if revokeSessionID == 0 {
		err = s.sessionRepo.Create(ctx, session)
	} else {
		err = s.sessionRepo.Rotate(ctx, revokeSessionID, session)
	}
	if err != nil {
		return nil, err
	}

	result := &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User: AuthUser{
			ID:          strconv.FormatUint(user.ID, 10),
			Username:    user.Username,
			DisplayName: user.DisplayName,
		},
	}

	if s.roleRepo != nil {
		roles, err := s.roleRepo.ListByUserID(ctx, user.ID)
		if err == nil {
			result.User.PlatformRoles = make([]string, 0, len(roles))
			for _, role := range roles {
				if strings.TrimSpace(role.Name) != "" {
					result.User.PlatformRoles = append(result.User.PlatformRoles, role.Name)
				}
			}
		}
	}

	return result, nil
}

func claimsExpiry(claims *Claims) time.Time {
	if claims == nil || claims.ExpiresAt == nil {
		return time.Now()
	}
	return claims.ExpiresAt.Time
}

func claimsTTLSeconds(claims *Claims) int64 {
	if claims == nil || claims.ExpiresAt == nil || claims.IssuedAt == nil {
		return 0
	}
	seconds := int64(claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time).Seconds())
	if seconds < 0 {
		return 0
	}
	return seconds
}
