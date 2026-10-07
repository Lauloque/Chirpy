/* SPDX-License-Identifier: GPL-3.0-or-later */
package auth

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	password := "Donut8"
	hashed := "$argon2id$v=19$m=65536,t=1,p=16$CNg/+0uxNn0gQcTptmpgxA$EJc3xTHkhWgjVic7hMNHnWLgeZi1cqwpyvHR/E1To2Y"

	success, err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Errorf("CheckPasswordHash error: %v", err)
	}
	if !success {
		t.Errorf("CheckPasswordHash failed")
	}

	success, err = CheckPasswordHash("incorrect", hashed)
	if err != nil {
		t.Errorf("CheckPasswordHash error: %v", err)
	}
	if success {
		t.Errorf("CheckPasswordHash should have found unmatched hash.")
	}

	success, err = CheckPasswordHash(password, "incorrect")
	if err == nil {
		t.Errorf("CheckPasswordHash should have errored due to hash being in the incorrect format")
	}
	if success {
		t.Errorf("CheckPasswordHash should have found unmatched hash.")
	}
}
