package middleware

type CtxKey int

const (
	ctxKeyUUID CtxKey = iota
	ctxKeyUserId
	ctxKeyBody
)
