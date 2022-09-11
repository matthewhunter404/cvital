package db

import "fmt"

const ErrInternal errorType = "internal"
const ErrNotFound errorType = "not_found"
const ErrUniqueViolation errorType = "unique_violation"

type errorType string

func (e errorType) Error() string {
	return fmt.Sprintf(string(e))
}

func WrapError(sentinelError, underlyingError error) error {
	return fmt.Errorf("%w: %v", sentinelError, underlyingError)
}
