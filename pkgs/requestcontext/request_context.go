package requestcontext

import (
	"context"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	TokenKey   contextKey = "user_token"
	SessionKey contextKey = "session_id"
)

type RequestContext interface {
	SetUserID(ctx context.Context, userID string) context.Context
	SetToken(ctx context.Context, token string) context.Context
	SetSessionID(ctx context.Context, sessionID string) context.Context
	GetUserID(ctx context.Context) (string, bool)
	GetToken(ctx context.Context) (string, bool)
	GetSessionID(ctx context.Context) (string, bool)
}

type requestContext struct {
	i         *di.Injector
	UserIDKey contextKey
	TokenKey  contextKey
}

func NewRequestContext(i *di.Injector) (RequestContext, error) {
	return &requestContext{
		i:         i,
		UserIDKey: UserIDKey,
		TokenKey:  TokenKey,
	}, nil
}

func (r *requestContext) SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, r.UserIDKey, userID)
}

func (r *requestContext) SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, r.TokenKey, token)
}

func (r *requestContext) SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionKey, sessionID)
}

func (r *requestContext) GetUserID(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(r.UserIDKey).(string)
	return UserID, ok
}

func (r *requestContext) GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(r.TokenKey).(string)
	return token, ok
}

func (r *requestContext) GetSessionID(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(SessionKey).(string)
	return sessionID, ok
}
