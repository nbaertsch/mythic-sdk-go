package mythic

import "fmt"

// Error represents a Mythic SDK error.
type Error struct {
	// Op is the operation that failed
	Op string

	// Err is the underlying error
	Err error

	// Message provides additional context
	Message string
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Op, e.Err)
	}
	return e.Op
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// Common error types
var (
	// ErrNotAuthenticated indicates the client is not authenticated
	ErrNotAuthenticated = fmt.Errorf("not authenticated")

	// ErrInvalidConfig indicates the configuration is invalid
	ErrInvalidConfig = fmt.Errorf("invalid configuration")

	// ErrAuthenticationFailed indicates authentication failed
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")

	// ErrTimeout indicates a request timed out
	ErrTimeout = fmt.Errorf("request timeout")

	// ErrNotFound indicates a resource was not found
	ErrNotFound = fmt.Errorf("not found")

	// ErrInvalidResponse indicates an unexpected response from the server
	ErrInvalidResponse = fmt.Errorf("invalid response")

	// ErrConnectionFailed indicates connection to the server failed
	ErrConnectionFailed = fmt.Errorf("connection failed")
)

// WrapError wraps an error with an operation and optional message.
func WrapError(op string, err error, message string) error {
	if err == nil {
		return nil
	}
	return &Error{
		Op:      op,
		Err:     err,
		Message: message,
	}
}
