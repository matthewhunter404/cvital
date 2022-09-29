package domain

import "fmt"

const ErrInternal errorType = "internal"
const ErrLoginFailed errorType = "login_failed"
const ErrNotFound errorType = "not_found"
const ErrAlreadyExists errorType = "already_exists"

type errorType string

func (e errorType) Error() string {
	return fmt.Sprintf(string(e))
}

func WrapError(sentinelError, underlyingError error) error {
	return fmt.Errorf("%w: %v", sentinelError, underlyingError)
}
