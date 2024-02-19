package middlewares

type JwtLevel string

const (
	WriteLevel JwtLevel = "write"
	ReadLevel  JwtLevel = "read"
)
