package libp2pmpquic

import (
	"errors"
	"github.com/ZHJBOBO/multipath-quic-go"
	"github.com/ZHJBOBO/multipath-quic-go/qerr"
	"github.com/libp2p/go-libp2p-core/mux"
)

const (
	reset qerr.ErrorCode = 19
)

type stream struct {
	quic.Stream
}

func (s *stream) Read(b []byte) (n int, err error) {
	n, err = s.Stream.Read(b)
	if err != nil && errors.Is(err, &qerr.QuicError{}) {
		err = mux.ErrReset
	}
	return n, err
}

func (s *stream) Write(b []byte) (n int, err error) {
	n, err = s.Stream.Write(b)
	if err != nil && errors.Is(err, &qerr.QuicError{}) {
		err = mux.ErrReset
	}
	return n, err
}

func (s *stream) Reset() error {
	s.Stream.Reset(reset)
	s.Stream.Reset(reset)
	return nil
}

func (s *stream) Close() error {
	//s.Stream.Reset(reset)
	return s.Stream.Close()
}

func (s *stream) CloseRead() error {
	s.Stream.Reset(reset)
	return nil
}

func (s *stream) CloseWrite() error {
	return s.Stream.Close()
}

var _ mux.MuxedStream = &stream{}
