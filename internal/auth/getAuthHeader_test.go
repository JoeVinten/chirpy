package auth

import (
	"net/http"
	"testing"
)

func TestGetBearer(t *testing.T) {

	testCases := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Missing Authorization header",
			headers: http.Header{
				"Auth": {"something"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed toke (no 'Bearer' prefix)",
			headers: http.Header{
				"Authorization": {"token"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed token (Bearer prefix but no token)",
			headers: http.Header{
				"Authorization": {"Bearer"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Valid token with extra spacing",
			headers: http.Header{
				"Authorization": {"Bearer   token  "},
			},
			wantToken: "token",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tc.headers)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && gotToken != tc.wantToken {
				t.Errorf("GetBearerToken() = %q, want %q", gotToken, tc.wantToken)
			}
		})
	}

}
