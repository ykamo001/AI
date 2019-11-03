package request

type RequestKey string

const (
	Origin        = RequestKey("Origin")
	ContentType   = RequestKey("Content-Type")
	Cookie        = RequestKey("Cookie")
	Authorization = RequestKey("Authorization")
	XSignature    = RequestKey("x-signature")
)
