package unit

import (
	"errors"
	"testing"

	"github.com/your-org/mythic-sdk-go/pkg/mythic"
)

func TestError(t *testing.T) {
	baseErr := errors.New("base error")
	err := &mythic.Error{
		Op:      "TestOperation",
		Err:     baseErr,
		Message: "additional context",
	}

	expected := "TestOperation: additional context: base error"
	if err.Error() != expected {
		t.Errorf("Error.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestErrorWithoutMessage(t *testing.T) {
	baseErr := errors.New("base error")
	err := &mythic.Error{
		Op:  "TestOperation",
		Err: baseErr,
	}

	expected := "TestOperation: base error"
	if err.Error() != expected {
		t.Errorf("Error.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestErrorWithoutErr(t *testing.T) {
	err := &mythic.Error{
		Op: "TestOperation",
	}

	expected := "TestOperation"
	if err.Error() != expected {
		t.Errorf("Error.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestErrorUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	err := &mythic.Error{
		Op:  "TestOperation",
		Err: baseErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Error.Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestWrapError(t *testing.T) {
	baseErr := errors.New("base error")
	wrapped := mythic.WrapError("TestOp", baseErr, "context")

	if wrapped == nil {
		t.Fatal("WrapError() returned nil")
	}

	var mythicErr *mythic.Error
	if !errors.As(wrapped, &mythicErr) {
		t.Error("WrapError() did not return *mythic.Error")
	}

	if mythicErr.Op != "TestOp" {
		t.Errorf("WrapError() Op = %q, want %q", mythicErr.Op, "TestOp")
	}

	if mythicErr.Message != "context" {
		t.Errorf("WrapError() Message = %q, want %q", mythicErr.Message, "context")
	}

	if !errors.Is(wrapped, baseErr) {
		t.Error("WrapError() lost base error in chain")
	}
}

func TestWrapErrorNil(t *testing.T) {
	wrapped := mythic.WrapError("TestOp", nil, "context")
	if wrapped != nil {
		t.Errorf("WrapError(nil) = %v, want nil", wrapped)
	}
}

func TestErrorConstants(t *testing.T) {
	// Just verify they exist and are non-nil
	errs := []error{
		mythic.ErrNotAuthenticated,
		mythic.ErrInvalidConfig,
		mythic.ErrAuthenticationFailed,
		mythic.ErrTimeout,
		mythic.ErrNotFound,
		mythic.ErrInvalidResponse,
		mythic.ErrConnectionFailed,
	}

	for _, err := range errs {
		if err == nil {
			t.Error("Error constant is nil")
		}
		if err.Error() == "" {
			t.Error("Error constant has empty message")
		}
	}
}
