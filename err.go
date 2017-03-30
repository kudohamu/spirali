package spirali

import "errors"

// Various errors the spirali might return.
var (
	ErrUnknownDriver = errors.New("unknown driver")
	ErrEnvNotFound   = errors.New("env not found in config")
)
