package model

import (
	"encoding/hex"
	"fmt"
)

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User encapsulates all identity-related information about a user. That includes the public GUID, the internal
// database ID, the ORCiD identifier and the hashed secret.
type User struct {
	// Non-optional attributes
	Id           int64
	GUID         string
	HashedSecret string
	// Optional attributes
	OrcId string
}

// Authorize checks wether the provided secret matches the hashed secret of a user.
func (u *User) Authorize(secret string) error {

	// convert the HEX-encoded hash from string to []byte
	bHash, err := hex.DecodeString(u.HashedSecret)
	if err != nil {
		return fmt.Errorf("Internal error. Could not convert secret hash from string to binary format. %s", err)
	}

	err = bcrypt.CompareHashAndPassword(bHash, []byte(secret))
	if err != nil {
		return fmt.Errorf("Authorization failure. Password does not match hashed secret for GUID '%s'. %s", u.GUID, err)
	}

	return nil
}

// Equals compares two user data structures for equality. This function is used to compare user as part of unit tests.
func (a *User) Equals(b *User) bool {
	return a.Id == b.Id &&
		a.GUID == b.GUID &&
		a.HashedSecret == b.HashedSecret &&
		a.OrcId == b.OrcId
}

// NewUserWithSecret takes a user-defined secret and creates a new user object with a fresh GUID and the hashed
// version of the secret.
func NewUserWithSecret(secret string) (u User, err error) {

	// Compute hash for salted secret.
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("Could not compute salted password hash. %s", err)
		return
	}

	// Create fresh, unused GUID, assuming that there is no collision.
	freshGUID, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("Could not create fresh user GUID. %s", err)
		return
	}

	// Encode the GUID and the hash as strings.
	u.GUID = freshGUID.String()
	u.HashedSecret = hex.EncodeToString(hash)

	return u, nil
}
