/* SPDX-License-Identifier: GPL-3.0-or-later */
package auth

import "github.com/alexedwards/argon2id"

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
