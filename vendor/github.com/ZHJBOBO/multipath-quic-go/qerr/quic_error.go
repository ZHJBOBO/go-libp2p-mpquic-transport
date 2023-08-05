package qerr

import (
	"fmt"
	"github.com/ZHJBOBO/multipath-quic-go/internal/protocol"
	"github.com/ZHJBOBO/multipath-quic-go/internal/utils"
)

// ErrorCode can be used as a normal error without reason.
type ErrorCode uint32

func (e ErrorCode) Error() string {
	return e.String()
}

// A QuicError consists of an error code plus a error reason
type QuicError struct {
	ErrorCode    ErrorCode
	ErrorMessage string
}

// Error creates a new QuicError instance
func Error(errorCode ErrorCode, errorMessage string) *QuicError {
	return &QuicError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}

func (e *QuicError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode.String(), e.ErrorMessage)
}

func (e *QuicError) Timeout() bool {
	switch e.ErrorCode {
	case NetworkIdleTimeout,
		HandshakeTimeout,
		TimeoutsWithOpenStreams:
		return true
	}
	return false
}

// ToQuicError converts an arbitrary error to a QuicError. It leaves QuicErrors
// unchanged, and properly handles `ErrorCode`s.
func ToQuicError(err error) *QuicError {
	switch e := err.(type) {
	case *QuicError:
		return e
	case ErrorCode:
		return Error(e, "")
	}
	utils.Errorf("Internal error: %v", err)
	return Error(InternalError, err.Error())
}

type TransportError struct {
	Remote       bool
	FrameType    uint64
	ErrorCode    ErrorCode
	ErrorMessage string
}
type ApplicationError struct {
	Remote       bool
	ErrorCode    ErrorCode
	ErrorMessage string
}

// A StatelessResetToken is a stateless reset token.
type StatelessResetToken [16]byte

// A StatelessResetError occurs when we receive a stateless reset.
type StatelessResetError struct {
	Token StatelessResetToken
}
type HandshakeTimeoutError struct{}

// A VersionNegotiationError occurs when the client and the server can't agree on a QUIC version.
type VersionNegotiationError struct {
	Ours   []protocol.VersionNumber
	Theirs []protocol.VersionNumber
}
type IdleTimeoutError struct{}
