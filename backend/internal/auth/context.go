package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

type contextKey string

const userContextKey contextKey = "auth_user"

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

// MustUserID is safe only inside handlers mounted behind RequireUser.
func MustUserID(ctx context.Context) uuid.UUID {
	user, ok := UserFromContext(ctx)
	if !ok {
		panic("auth: user missing from context")
	}
	return user.ID
}

func RequireUser(uc *Usecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				httpx.Error(w, r, httpx.Unauthenticated("authentication required"))
				return
			}

			user, err := uc.Authenticate(r.Context(), cookie.Value)
			if err != nil {
				httpx.Error(w, r, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
		})
	}
}
