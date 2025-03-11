package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
)

type Auth struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New returns a new instance of the Auth service
func New(userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		appProvider: appProvider,
		usrProvider: userProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(ctx context.Context, email, passHash string, appID int) (string, error) {
	const op = "auth.Login"

	slog.With(
		slog.String("op", op),
	)

	slog.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			slog.Warn("user not found", slog.Any("error", err))

			return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
		}

		slog.Error("failed to get user", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(passHash)); err != nil {
		slog.Info("invalid credentials", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	slog.Info("successfully logged in", slog.Any("user", user))

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		slog.Error("failed to create token", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, email, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	slog.With(
		slog.String("op", op),
	)

	slog.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		slog.Error("failed to save user", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	slog.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
