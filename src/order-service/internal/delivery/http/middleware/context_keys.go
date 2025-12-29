package my_middleware

type CtxKey int

const (
	ctxKeyUUID CtxKey = iota
	ctxKeyUserId
	ctxKeyPage
	ctxKeyLimit
	ctxKeyBody
)
