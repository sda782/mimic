package types

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
