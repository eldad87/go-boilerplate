package middleware

type contextKey string

func (c contextKey) String() string {
	return string(c)
}
