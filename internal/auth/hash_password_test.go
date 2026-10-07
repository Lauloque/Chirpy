/* SPDX-License-Identifier: GPL-3.0-or-later */
package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "Donut8"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword error: %v", err)
	}
	if len(hashed) != 98 {
		t.Errorf("HashPassword expected a 98 characters hash, got '%d', with error: %v", len(hashed), err)
	}

	success, err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Errorf("CheckPasswordHash error: %v", err)
	}
	if !success {
		t.Errorf("CheckPasswordHash failed")
	}
}
