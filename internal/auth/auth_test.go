/* SPDX-License-Identifier: GPL-3.0-or-later */
package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "Donut8"
	password2 := "Baguette123"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct Password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect Password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty Password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid Password",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	tokenSecret1 := "Donut8"
	tokenSecret2 := "Baguette123"
	uuid1 := uuid.New()
	durationPos := time.Hour
	durationNeg := -time.Hour

	tests := []struct {
		name        string
		userID      uuid.UUID
		signSecret  string
		checkSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Same secrets, positive duration, no error wanted",
			userID:      uuid1,
			signSecret:  tokenSecret1,
			checkSecret: tokenSecret1,
			expiresIn:   durationPos,
			wantErr:     false,
		},
		{
			name:        "Different secrets, positive duration, error wanted",
			userID:      uuid1,
			signSecret:  tokenSecret1,
			checkSecret: tokenSecret2,
			expiresIn:   durationPos,
			wantErr:     true,
		},
		{
			name:        "Same secrets, negative duration, error wanted",
			userID:      uuid1,
			signSecret:  tokenSecret1,
			checkSecret: tokenSecret1,
			expiresIn:   durationNeg,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tt.signSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() error: %v", err)
			}

			newID, err := ValidateJWT(tokenString, tt.checkSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != true {
				if newID != tt.userID {
					t.Errorf("ValidateJWT() generated uuid %v, expected %v", newID, tt.userID)
				}
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {

	tests := []struct {
		name          string
		inputToken    string
		expectedToken string
		genHeader     bool
		wantErr       bool
	}{
		{
			name:          "Matching tokens, no error wanted",
			inputToken:    "Bearer Donut",
			expectedToken: "Donut",
			genHeader:     true,
			wantErr:       false,
		},
		{
			name:          "Token with space, no error wanted",
			inputToken:    "Bearer Sunflower 6",
			expectedToken: "Sunflower 6",
			genHeader:     true,
			wantErr:       false,
		},
		{
			name:          "No header gen, error wanted",
			inputToken:    "owo",
			expectedToken: "",
			genHeader:     false,
			wantErr:       true,
		},
		{
			name:          "No Bearer token",
			inputToken:    "owo",
			expectedToken: "",
			genHeader:     true,
			wantErr:       true,
		},
		{
			name:          "No space separator",
			inputToken:    "BearerOwo",
			expectedToken: "",
			genHeader:     true,
			wantErr:       true,
		},
		{
			name:          "Empty token",
			inputToken:    "",
			expectedToken: "",
			genHeader:     true,
			wantErr:       true,
		},
		{
			name:          "Only spaces token",
			inputToken:    "    ",
			expectedToken: "",
			genHeader:     true,
			wantErr:       true,
		},
		{
			name:          "Trailing spaces separator",
			inputToken:    "Bearer Owo ",
			expectedToken: "Owo",
			genHeader:     true,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.genHeader {
				h.Set("Authorization", tt.inputToken)
			}

			token, err := GetBearerToken(h)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if token != tt.expectedToken {
				t.Errorf("GetBearerToken() expected token %v, got %v", tt.expectedToken, token)
			}
		})
	}
}
