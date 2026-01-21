package middleware

type сtxKey int

const (
	ctxKeyUUID сtxKey = iota
	ctxKeyBody
	ctxKeyPage
	ctxKeyLimit
)
