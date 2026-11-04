package auth // Groups authentication helpers into a reusable package.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time" // Provides timestamps and token lifetimes.

	"github.com/alexedwards/argon2id" // Hashes and verifies passwords using Argon2id.
	"github.com/golang-jwt/jwt/v5"    // Creates, signs, parses, and validates JWTs.
	"github.com/google/uuid"          // Represents and parses user IDs.
)

// HashPassword converts a password into a salted, encoded hash.
// WHY: Store password hashes, never plaintext passwords.
// PRODUCTION: These defaults are a starting point; benchmark memory and CPU
// costs for your infrastructure. Never log passwords.
func HashPassword(password string) (string, error) {
	// CreateHash generates a random salt and includes the algorithm parameters,
	// salt, and hash in its returned string. Propagate any hashing error.
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

// CheckPasswordHash checks a supplied password against a stored hash.
// WHY: Salted hashes cannot be compared by simply hashing again and comparing
// strings; the verifier reads the salt and parameters from the stored hash.
func CheckPasswordHash(password, hash string) (bool, error) {
	// Returns whether the password matches, plus any verification error.
	// PRODUCTION: Callers should distinguish malformed hashes/errors from
	// ordinary mismatches, without exposing sensitive details to clients.
	return argon2id.ComparePasswordAndHash(password, hash)
}

// MakeJWT creates a signed access token identifying a user.
// BOOT.DEV: "chirpy-access" is the issuer required by this assignment.
// PRODUCTION: Choose an issuer appropriate to your system and validate it.
// HS256, UUID subjects, and this function's signature are design choices,
// not requirements for every production authentication system.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Capture the time once so both timestamps share the same starting point.
	// UTC makes time handling consistent; JWT timestamps encode an instant,
	// not a timezone.
	now := time.Now().UTC()

	// Create an unsigned token with an explicit signing algorithm and claims.
	// WHY: Claims describe the identity and validity period of the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// Identifies the system that issued the token.
		// Merely setting this does NOT make validators enforce it.
		Issuer: "chirpy-access",

		// Records when the token was created.
		// NewNumericDate converts time.Time into the library's JWT date type.
		// IssuedAt is not checked by default by this library.
		IssuedAt: jwt.NewNumericDate(now),

		// Defines when the token expires, using the caller's requested lifetime.
		// PRODUCTION: Keep access-token lifetimes limited by server policy.
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),

		// Stores the user's identity in the standard "sub" claim.
		// JWT subjects are strings, so convert the UUID to its string form.
		Subject: userID.String(),
	})

	// Sign and encode the token, returning either its string or an error.
	// HS256 requires a []byte key, hence the conversion from string.
	// PRODUCTION: Use a strong, randomly generated secret stored securely,
	// not a hardcoded password. HS256 uses the same secret to sign and verify.
	//
	// IMPORTANT: Signing prevents undetected modification; it does not encrypt
	// the claims. Anyone holding the token can read its contents.
	return token.SignedString([]byte(tokenSecret))
}

// ValidateJWT verifies a token and returns the user ID from its subject.
// BOOT.DEV: This implements the lesson's signature, expiration, and UUID checks.
// PRODUCTION: Also enforce the issuer, audience, and required claims that
// your application relies on; see the notes below.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Allocate a destination for the decoded registered claims.
	// WHY: ParseWithClaims populates this struct while parsing the token.
	thisTokenClaims := jwt.RegisteredClaims{}

	// Define the callback required by the JWT library to obtain a verification key.
	// WHY: The library supports applications that select different keys based
	// on token headers. This application always uses one shared secret.
	getSecretAsBytes := func(token *jwt.Token) (interface{}, error) {
		// The token parameter is unused here; its signature is not verified yet.
		// Do not treat information in it as authenticated at this stage.

		// Return the HS256 key in the required byte-slice form.
		// interface{} allows the library to support multiple key types.
		// nil means key lookup succeeded; this callback does no external lookup.
		return []byte(tokenSecret), nil
	}

	// Parse the encoded token, verify its signature, and validate its claims.
	token, err := jwt.ParseWithClaims(
		tokenString,                             // The encoded JWT supplied by the caller.
		&thisTokenClaims,                        // Pointer lets the parser populate claims.
		getSecretAsBytes,                        // Supplies the key for signature verification.
		jwt.WithValidMethods([]string{"HS256"}), // Reject other signing algorithms.
	)

	// Reject malformed tokens, invalid signatures, expired tokens, and other
	// parsing/validation failures before trusting any claims.
	// Expiration is checked when present; this call does not REQUIRE an exp claim.
	if err != nil {
		// uuid.Nil is a placeholder; callers must handle the non-nil error.
		return uuid.Nil, err
	}

	// Read the standard subject claim from the successfully validated token.
	// WHY: MakeJWT stored the user's UUID here.
	subj, err := token.Claims.GetSubject()

	// Follow the Claims interface's error contract when extracting the subject.
	if err != nil {
		return uuid.Nil, err
	}

	// Convert the subject string back into the UUID expected by this application.
	// WHY: A valid signature does not guarantee that a claim has the right format.
	thisUUID, err := uuid.Parse(subj)

	// Reject missing or malformed UUID subjects.
	if err != nil {
		return uuid.Nil, err
	}

	// Return the authenticated user ID and indicate success.
	// PRODUCTION: Authentication is not authorization; callers must still check
	// whether this user may access the requested resource.
	return thisUUID, nil
}

// PRODUCTION HARDENING NOTES:
//
// Depending on your system, add parser options such as:
//   jwt.WithIssuer("chirpy-access") // Require the expected issuer.
//   jwt.WithExpirationRequired()   // Reject tokens without an expiration.
//   jwt.WithAudience("your-api")   // Require an audience you also issue.
//
// Also use HTTPS, protect tokens from leakage, plan key rotation, and decide
// how logout/revocation works. A signed JWT alone does not provide those.
//
// None of the cryptographic operations above are "just for Boot.dev."
// The project-specific names and chosen policies are the lesson-specific parts.

// ------ end auto-generated comments -----

func GetBearerToken(headers http.Header) (string, error) {
	// Bearer TOKEN_STRING
	// return TOKEN_STRING
	token := headers.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("Header Authorization not present")
	}

	tokenSlice := strings.Split(token, " ")

	if len(tokenSlice) != 2 {
		return "", fmt.Errorf("Malformed Bearer TOKEN header")
	}
	if tokenSlice[0] != "Bearer" {
		return "", fmt.Errorf("Invalid Authorization header typ")
	}

	return tokenSlice[1], nil
}

func getAuthTokenFromHeader(headers http.Header) (string, error) {
	// ApiKey TOKEN_STRING
	// return TOKEN_STRING
	token := headers.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("Header Authorization not present")
	}

	tokenSlice := strings.Split(token, " ")

	if len(tokenSlice) != 2 {
		return "", fmt.Errorf("Malformed Auth TOKEN header")
	}
	if tokenSlice[0] != "ApiKey" {
		return "", fmt.Errorf("Invalid Authorization header typ")
	}

	return tokenSlice[1], nil
}

func MakeRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GetAPIKey(headers http.Header) (string, error) {
	return getAuthTokenFromHeader(headers)
}
