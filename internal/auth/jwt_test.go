package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	secretValid   = "validSecret"
	secretInvalid = "invalidSecret"
)

func TestJWTSignature(t *testing.T) {
	testCases := map[string]struct {
		id       uuid.UUID
		validFor time.Duration
		secret   string
		valid    bool
	}{
		"base case": {
			id:       uuid.New(),
			secret:   secretValid,
			validFor: time.Second * 10,
			valid:    true,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := MakeJWT(test.id, test.secret, test.validFor)
			if err != nil && test.valid {
				t.Fatalf("Failed with error %#v", err)
			}
			if err == nil && !test.valid {
				t.Fatalf("Expected failure for payload %#v", test)
			}
		})
	}
}

const (
	ErrorIDValidation = `ID lost in validation
expected: %v
got: %v
`
)

func TestJWTValidation(t *testing.T) {
	testCases := map[string]struct {
		secret   string
		valid    bool
		id       uuid.UUID
		validFor time.Duration
	}{
		"Validate base case": {
			id:       uuid.New(),
			secret:   secretValid,
			valid:    true,
			validFor: time.Hour,
		},
		"Invalid secret": {
			id:       uuid.New(),
			secret:   secretInvalid,
			valid:    false,
			validFor: time.Hour,
		},
		"Expired token": {
			id:       uuid.New(),
			secret:   secretValid,
			valid:    false,
			validFor: 0,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			jwtStr, err := MakeJWT(test.id, test.secret, test.validFor)
			if err != nil {
				t.Fatal("Failed due to make jwt (impropet test case setup)")
			}
			uid, err := ValidateJWT(jwtStr, secretValid)
			if err != nil {
				if test.valid {
					t.Fatalf("Validation failed for payload %#v", test)
				}
			} else {
				if !test.valid {
					t.Fatalf("Expected failure for payload %#v", test)
				} else {
					if uid.String() != test.id.String() {
						t.Fatalf(ErrorIDValidation, test.id, uid)
					}
				}

			}
		})
	}
}
