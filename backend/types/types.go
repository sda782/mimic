package types

import "context"

type User struct {
	ID           int
	Name         string
	PasswordHash string
	SessionToken string
}

type Upload struct {
	ID        int
	UserID    int
	FileName  string
	FilePath  string
	ShortCode string
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func GetUser(ctx context.Context) User {
	return ctx.Value("user").(User)
}
