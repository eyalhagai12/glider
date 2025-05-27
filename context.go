package backend

import "context"

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(*User); ok {
		return user
	}

	return nil
}
