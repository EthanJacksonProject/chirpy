package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func getUUID() uuid.UUID {
	uuid, _ := uuid.NewUUID()
	return uuid
}

func minToDuration(minToSurvive int) time.Duration {
	return time.Minute * time.Duration(minToSurvive)
}
func TestRoundTripJWT(t *testing.T) {
	// Arrange: create a user UUID, a secret, and a positive lifetime.
	uuid := getUUID()
	expiresAt := minToDuration(5)
	secret := "My_Secret_123"

	// Act: call MakeJWT with those values.
	// If it returns an error, stop with t.Fatalf.
	str, err := MakeJWT(uuid, secret, expiresAt)
	if err != nil {
		t.Fatalf("failed to get jwt: %v", err)
	}

	// Act: pass the resulting token and the same secret to ValidateJWT.
	// If it returns an error, stop with t.Fatalf.
	newUuid, err := ValidateJWT(str, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	// Assert: compare the returned UUID with the original UUID.
	// If they differ, report the mismatch with t.Errorf.

	if uuid != newUuid {
		t.Errorf("UUID does not match validation")
		return
	}
}

func TestRoundTripPWHash(t *testing.T) {
	testPassword := "___^&%abraCadabra1234!@#$"

	hashed, err := HashPassword(testPassword)

	if err != nil {
		t.Fatalf("failed to hash pw: %v", err)
	}

	pwMatched, err := CheckPasswordHash(testPassword, hashed)

	if err != nil {
		t.Fatalf("failed to check hashed pw: %v", err)
	}

	if !pwMatched {
		t.Errorf("Hashed password does not match validation")
		return
	}
}

func TestGetBearerToken(t *testing.T) {

	goodHdr := http.Header{}
	goodHdr.Set("Authorization", "Bearer 12345abcdeGood")

	// wrong token
	badHdr1 := http.Header{}
	badHdr1.Set("Authorization", "BooBaa 12345abcde")

	// wrong key
	badHdr2 := http.Header{}
	badHdr2.Set("Authorization", "BooBaa 12345abcdeGood")

	// wrong length
	badHdr3 := http.Header{}
	badHdr3.Set("Authorization", "12345abcde")

	goodStr, goodErr := GetBearerToken(goodHdr)
	str1, err1 := GetBearerToken(badHdr1)
	str2, err2 := GetBearerToken(badHdr2)
	str3, err3 := GetBearerToken(badHdr3)

	if goodErr != nil {
		t.Fatalf("Failed to get good bearer token: %v", goodErr)
	}
	if str1 != "" || err1 == nil {
		t.Fatalf("[1] Failed to get good bearer token: %v", err1)
	}
	if str2 != "" || err2 == nil {
		t.Fatalf("[2] Failed to get good bearer token: %v", err2)
	}
	if str3 != "" || err3 == nil {
		t.Fatalf("[3] Failed to get good bearer token: %v", err3)
	}

	if goodStr != "12345abcdeGood" {
		t.Errorf("Token does not match expected value")
		return
	}
}
