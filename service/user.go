package service

import (
	"context"
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/OVillas/e-commercer-api/secure"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type userService struct {
	i               *do.Injector
	userRespository domain.UserRepository
	sessionService  domain.SessionService
	emailService    domain.EmailService
}

func NewUserService(i *do.Injector) (domain.UserService, error) {
	userRepository, err := do.Invoke[domain.UserRepository](i)
	if err != nil {
		return nil, err
	}

	sessionService, err := do.Invoke[domain.SessionService](i)
	if err != nil {
		return nil, err
	}

	emailService, err := do.Invoke[domain.EmailService](i)
	if err != nil {
		return nil, err
	}

	return &userService{
		i:               i,
		userRespository: userRepository,
		sessionService:  sessionService,
		emailService:    emailService,
	}, nil
}

func (u *userService) Create(ctx context.Context, userPayload domain.UserPayLoad) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	user, err := u.userRespository.GetByEmail(ctx, userPayload.Email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return err
	}

	if user != nil {
		log.Warn("User already exists")
		return domain.ErrUserAlreadyExists
	}

	passwordHash, err := secure.Hash(userPayload.Password)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return domain.ErrHashingPassword
	}

	user = userPayload.ToUser(string(passwordHash))

	if err := u.userRespository.Create(ctx, *user); err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("User created successfully")
	return nil
}

func (u *userService) SignIn(ctx context.Context, signInPayload domain.SignInPayLoad) (*domain.SessionResponse, error) {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "SignIn"),
	)

	log.Info("Initializing user sign in process")

	user, err := u.userRespository.GetByEmail(ctx, signInPayload.Email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	if user == nil {
		log.Warn("User not found")
		return nil, domain.ErrUserNotFound
	}

	if err := secure.CheckPassword(user.PasswordHash, signInPayload.Password); err != nil {
		log.Warn("Invalid password")
		return nil, domain.ErrInvalidPassword
	}

	token, err := u.sessionService.Create(ctx, *user)
	if err != nil {
		log.Error("Failed to create token", slog.String("error", err.Error()))
		return nil, err
	}

	if !user.EmailConfirmed {
		log.Warn("User email not confirmed")

		if err := u.emailService.SendConfirmationCode(ctx, *user); err != nil {
			log.Error("Failed to send OTP email", slog.String("error", err.Error()))
			return nil, err
		}

		return &domain.SessionResponse{
			Token: token,
		}, domain.ErrEmailNotConfirmed
	}

	log.Info("User signed in successfully")

	return &domain.SessionResponse{
		Token: token,
	}, nil
}

func (u *userService) CheckEmailConfirmation(ctx context.Context, email string) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "CheckEmailConfirmation"),
	)

	log.Info("Initializing email confirmation check process")

	user, err := u.userRespository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return err
	}

	if user == nil {
		log.Warn("User not found")
		return domain.ErrUserNotFound
	}

	if !user.EmailConfirmed {
		log.Warn("User email not confirmed")
		return domain.ErrEmailNotConfirmed
	}

	log.Info("check User email confirmed executed successfully")
	return nil
}

func (u *userService) UpdateName(ctx context.Context, name string) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "UpdateName"),
	)

	log.Info("Initializing user name update process")

	user, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || user == nil {
		return domain.ErrUserNotFoundInContext
	}

	if user.Name == name {
		log.Warn("New name is the same as the old name")
		return domain.ErrNameIsSame
	}

	if err := u.userRespository.UpdateName(ctx, user.UserID, name); err != nil {
		log.Error("Failed to update user name", slog.String("error", err.Error()))
		return err
	}

	if err := u.sessionService.Update(ctx); err != nil {
		log.Error("Failed to update token", slog.String("error", err.Error()))
		return err
	}

	log.Info("User name updated successfully")
	return nil
}

func (u *userService) UpdatePassword(ctx context.Context, updatePasswordPayload domain.UpdatePasswordPayload) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "UpdatePassword"),
	)

	log.Info("Initializing user password update process")

	userResponse, ok := ctx.Value(middleware.UserKey).(*domain.UserResponse)
	if !ok || userResponse == nil {
		log.Warn("User not found in context")
		return domain.ErrUserNotFoundInContext
	}

	userID, err := uuid.Parse(userResponse.ID)
	if err != nil {
		log.Error("Failed to parse user ID", slog.String("error", err.Error()))
		return err
	}

	user, err := u.userRespository.GetByID(ctx, userID)
	if err != nil {
		log.Error("Failed to get user by ID", slog.String("error", err.Error()))
		return err
	}

	if err := secure.CheckPassword(user.PasswordHash, updatePasswordPayload.OldPassword); err != nil {
		log.Warn("Invalid old password")
		return domain.ErrInvalidOldPassword
	}

	if updatePasswordPayload.OldPassword == updatePasswordPayload.NewPassword {
		log.Warn("New password is the same as the old password")
		return domain.ErrPasswordIsSame
	}

	passwordHash, err := secure.Hash(updatePasswordPayload.NewPassword)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return domain.ErrHashingPassword
	}

	if err := u.userRespository.UpdatePassword(ctx, userID, string(passwordHash)); err != nil {
		log.Error("Failed to update user password", slog.String("error", err.Error()))
		return err
	}

	if err := u.sessionService.Update(ctx); err != nil {
		log.Error("Failed to update token", slog.String("error", err.Error()))
		return err
	}

	log.Info("User password updated successfully")
	return nil
}

func (u *userService) GetUserInfo(ctx context.Context) (*domain.UserResponse, error) {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "GetUserInfo"),
	)

	log.Info("Initializing get user info process")

	user, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || user == nil {
		return nil, domain.ErrUserNotFoundInContext
	}
	log.Info("user info requested successfully")
	return user.ToResponse(), nil
}

func (u *userService) ResendCode(ctx context.Context, resendCodePayload domain.ResendCodePayload) error {
	log := slog.With(
		slog.String("service", "user"),
		slog.String("func", "ResendCode"),
	)

	log.Info("Initializing resend code process")

	userSession, _ := ctx.Value(middleware.UserKey).(*domain.Session)

	var email string
	if userSession != nil {
		email = userSession.Email
	} else {
		email = resendCodePayload.Email
	}

	user, err := u.userRespository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return err
	}

	if err := u.emailService.SendConfirmationCode(ctx, *user); err != nil {
		log.Error("Failed to send OTP email", slog.String("error", err.Error()))
		return err
	}

	log.Info("confirmatio code send/resend successfully")
	return nil
}
