package decoder_test

import (
	"testing"

	dt "github.com/jinsoo-youn/traefik-jwt-decode/decodertest"

	"github.com/jinsoo-youn/traefik-jwt-decode/decoder"
)

func TestInvalidJwksURLGivesWarning(t *testing.T) {
	claimMapping := make(map[string]string)
	dec, err := decoder.NewJwsDecoder("https://this.com/does/not/exist", claimMapping)
	dt.Report(t, dec == nil && err == nil, "not able to create jws decoder with incorrect jwks url, and no warning given: %s", err)
}

func TestValidJwksURL(t *testing.T) {
	claimMapping := make(map[string]string)
	dec, err := decoder.NewJwsDecoder("https://www.googleapis.com/oauth2/v3/certs", claimMapping)
	dt.Report(t, dec == nil || err != nil, "not able to create jws decoder with correct jwks url")
}

func TestTokenWithInvalidClaims(t *testing.T) {
	invalidTokens := map[string]interface{}{
		"int":    123,
		"double": 123.321,
		"struct": struct{ key string }{key: "123"},
	}
	tc := dt.NewTest()
	claimKey := "claimKey"
	claimMapping := map[string]string{
		claimKey: "claim-email",
	}
	dec, _ := decoder.NewJwsDecoder(tc.JwksURL, claimMapping)
	for k, v := range invalidTokens {
		token := tc.NewValidToken(map[string]interface{}{claimKey: v})
		resp, err := dec.Decode(dt.Ctx(), string(token))
		dt.Report(t, err != nil, "not able to decode token with unusual JSON type: (%s : %+v) into %+v", k, v, resp)
	}
}
