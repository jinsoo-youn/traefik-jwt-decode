package decoder_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jinsoo-youn/traefik-jwt-decode/decoder"
	dt "github.com/jinsoo-youn/traefik-jwt-decode/decodertest"
)

func TestCacheAllResponses(t *testing.T) {
	tests := map[string]struct {
		token *decoder.Token
		err   error
	}{
		"TokenAndNoError": {token: &decoder.Token{Expiration: time.Now().Add(time.Hour)}, err: nil},
		"TokenAndError":   {token: &decoder.Token{Expiration: time.Now().Add(time.Hour)}, err: fmt.Errorf("some error")},
		"NoTokenAndError": {token: nil, err: fmt.Errorf("some other error")},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			delegate := newMock(func(ctx context.Context, raw string) (*decoder.Token, error) {
				return tc.token, tc.err
			})
			dec := decoder.NewCachedJwtDecoder(dt.Cache, delegate)
			getAndCompareCached(t, name, dec, delegate, tc.token, tc.err, 1)
			time.Sleep(100 * time.Millisecond)
			getAndCompareCached(t, name, dec, delegate, tc.token, tc.err, 1)
		})
	}
}

func getAndCompareCached(t *testing.T, name string, dec decoder.TokenDecoder, delegate *decoderMock, expectedToken *decoder.Token, expectedError error, expectedCalls int) {
	token, err := dec.Decode(dt.Ctx(), name)
	dt.Report(t, err != expectedError, "got unexpected error %v from cache expected %v", err, expectedError)
	dt.Report(t, token != expectedToken, "got unexpected token %v from cache expected %v", token, expectedToken)
	dt.Report(t, delegate.calls != expectedCalls, "incorrect number of calls to delegate %d expected %d", delegate.calls, expectedCalls)
}

type DecodeFunc func(ctx context.Context, raw string) (*decoder.Token, error)

type decoderMock struct {
	calls    int
	delegate DecodeFunc
}

func (d *decoderMock) Decode(ctx context.Context, raw string) (*decoder.Token, error) {
	d.calls = d.calls + 1
	return d.delegate(ctx, raw)
}

func newMock(delegate DecodeFunc) *decoderMock {
	return &decoderMock{calls: 0, delegate: delegate}
}
