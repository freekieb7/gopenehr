package oauth

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/freekieb7/gopenehr/internal/telemetry"
)

// ---------------- Errors ----------------
var (
	ErrUnsupported      = errors.New("unsupported token type")
	ErrInvalidFormat    = errors.New("invalid token format")
	ErrInvalidHeader    = errors.New("invalid header")
	ErrMissingKid       = errors.New("missing kid in header")
	ErrUnsupportedAlg   = errors.New("unsupported alg (only RS256)")
	ErrInvalidClaims    = errors.New("invalid claims")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrKidNotFound      = errors.New("signing key not found")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrNoKeyMaterial    = errors.New("no RSA key material found in JWK")
)

// ---------------- Config / Types ----------------

const defaultJWKSCache = time.Hour
const clockSkewLeeway = int64(60) // seconds

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type OpenIDConfig struct {
	JWKSURI string `json:"jwks_uri"`
}

type CachedJWKS struct {
	JWKS   JWKS
	Expiry time.Time
}

type Service struct {
	Logger *telemetry.Logger

	TrustedIssuers []string
	Audience       string
	HTTPClient     *http.Client

	cache map[string]CachedJWKS // keyed by jwks_uri
	mu    sync.Mutex
}

// NewService returns a pointer (so receiver methods work on it)
func NewService(logger *telemetry.Logger, trustedIssuers []string, audience string) Service {
	return Service{
		Logger:         logger,
		TrustedIssuers: trustedIssuers,
		Audience:       audience,
		HTTPClient:     &http.Client{Timeout: 10 * time.Second},
		cache:          make(map[string]CachedJWKS),
	}
}

func (s *Service) WarmupCache(ctx context.Context) error {
	for _, issuer := range s.TrustedIssuers {
		iss := strings.TrimRight(issuer, "/")
		jwksURI, err := s.discoverJWKSURI(ctx, iss)
		if err != nil {
			return fmt.Errorf("discover jwks_uri for issuer %s: %w", issuer, err)
		}
		if _, _, err := s.fetchAndCacheJWKS(ctx, jwksURI); err != nil {
			return fmt.Errorf("fetch and cache jwks for issuer %s: %w", issuer, err)
		}
	}
	return nil
}

// ---------------- Public API ----------------

// ValidateToken accepts either a JWT (three dot-separated parts) or an opaque token.
// Opaque token support not implemented here (returns ErrUnsupported).
func (s *Service) ValidateToken(ctx context.Context, token string) (map[string]any, error) {
	if strings.Count(token, ".") == 2 {
		return s.ParseJWT(ctx, token)
	}
	return nil, ErrUnsupported
}

// ParseJWT validates a RS256 JWT using JWKS discovered from issuer.
// Returns the claims map if valid.
func (s *Service) ParseJWT(ctx context.Context, token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	// decode header
	headerB, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: header decode: %v", ErrInvalidHeader, err)
	}
	var header map[string]any
	if err := json.Unmarshal(headerB, &header); err != nil {
		return nil, fmt.Errorf("%w: header json: %v", ErrInvalidHeader, err)
	}

	// alg check
	algI, ok := header["alg"]
	if !ok {
		return nil, fmt.Errorf("%w: alg missing", ErrUnsupportedAlg)
	}
	alg, ok := algI.(string)
	if !ok || alg != "RS256" {
		return nil, ErrUnsupportedAlg
	}

	// kid required
	kidI, ok := header["kid"]
	if !ok {
		return nil, ErrMissingKid
	}
	kid, ok := kidI.(string)
	if !ok || kid == "" {
		return nil, ErrMissingKid
	}

	// decode payload (use json.Number)
	payloadB, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: payload decode: %v", ErrInvalidClaims, err)
	}
	claims, err := decodeJSONToMap(payloadB)
	if err != nil {
		return nil, fmt.Errorf("%w: payload json: %v", ErrInvalidClaims, err)
	}

	// issuer must be present
	issI, ok := claims["iss"]
	if !ok {
		return nil, errors.New("missing iss claim")
	}
	iss, ok := issI.(string)
	if !ok || iss == "" {
		return nil, errors.New("invalid iss claim")
	}
	iss = strings.TrimRight(iss, "/")

	// issuer whitelist check (normalize)
	if !s.isTrustedIssuer(iss) {
		return nil, ErrInvalidIssuer
	}

	// aud check (if configured)
	if s.Audience != "" {
		if err := verifyAudience(claims["aud"], s.Audience); err != nil {
			return nil, ErrInvalidAudience
		}
	}

	// time claims: exp, nbf (with leeway)
	now := time.Now().Unix()
	if err := validateTimeClaimWithLeeway(claims, "exp", now, clockSkewLeeway); err != nil {
		return nil, err
	}
	if err := validateTimeClaimWithLeeway(claims, "nbf", now, clockSkewLeeway); err != nil {
		return nil, err
	}

	// discover jwks_uri and get JWKS (cache keyed by jwks_uri)
	jwksURI, err := s.discoverJWKSURI(ctx, iss)
	if err != nil {
		return nil, fmt.Errorf("discover jwks_uri: %w", err)
	}

	jwks, err := s.getJWKS(ctx, jwksURI)
	if err != nil {
		return nil, fmt.Errorf("get jwks: %w", err)
	}

	// find key by kid
	jwkPtr := findJWKByKid(jwks, kid)
	if jwkPtr == nil {
		// retry fetch once in case of rotation
		if _, _, err := s.fetchAndCacheJWKS(ctx, jwksURI); err == nil {
			jwks, _ = s.getJWKS(ctx, jwksURI)
			jwkPtr = findJWKByKid(jwks, kid)
		}
		if jwkPtr == nil {
			return nil, ErrKidNotFound
		}
	}

	// convert to RSA public key (n/e or x5c)
	pub, err := jwkToPublicRSAWithCertChecks(*jwkPtr)
	if err != nil {
		return nil, fmt.Errorf("jwk->rsa: %w", err)
	}

	// signature verify
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("signature decode: %w", err)
	}
	signingInput := []byte(parts[0] + "." + parts[1])
	hash := sha256.Sum256(signingInput)
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig); err != nil {
		return nil, ErrInvalidSignature
	}

	return claims, nil
}

// ---------------- Helpers: issuer/trusted ----------------

func (s *Service) isTrustedIssuer(iss string) bool {
	// normalize trusted list too
	for _, t := range s.TrustedIssuers {
		if strings.TrimRight(t, "/") == iss {
			return true
		}
	}
	return false
}

// ---------------- JWKS fetch & cache (keyed by jwks_uri) ----------------

func (s *Service) getJWKS(ctx context.Context, jwksURI string) (JWKS, error) {
	s.mu.Lock()
	c, ok := s.cache[jwksURI]
	s.mu.Unlock()

	if ok && time.Now().Before(c.Expiry) {
		return c.JWKS, nil
	}
	// fetch and cache
	jwks, _, err := s.fetchAndCacheJWKS(ctx, jwksURI)
	return jwks, err
}

func (s *Service) fetchAndCacheJWKS(ctx context.Context, jwksURI string) (JWKS, time.Time, error) {
	jwks, expiry, err := s.FetchJWKS(ctx, jwksURI)
	if err != nil {
		return JWKS{}, time.Time{}, err
	}
	s.mu.Lock()
	s.cache[jwksURI] = CachedJWKS{JWKS: jwks, Expiry: expiry}
	s.mu.Unlock()
	return jwks, expiry, nil
}

func (s *Service) discoverJWKSURI(ctx context.Context, issuer string) (string, error) {
	issuer = strings.TrimRight(issuer, "/")

	// ensure issuer is a valid URL
	if _, err := url.ParseRequestURI(issuer); err != nil {
		return "", fmt.Errorf("invalid issuer url: %w", err)
	}

	confURL := issuer + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, confURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			s.Logger.Error("Failed to close response body", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openid configuration fetch status %d", resp.StatusCode)
	}
	var cfg OpenIDConfig
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&cfg); err != nil {
		return "", err
	}
	if cfg.JWKSURI == "" {
		return "", errors.New("jwks_uri not found in openid-configuration")
	}
	return cfg.JWKSURI, nil
}

func (s *Service) FetchJWKS(ctx context.Context, jwksURI string) (JWKS, time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURI, nil)
	if err != nil {
		return JWKS{}, time.Time{}, err
	}
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return JWKS{}, time.Time{}, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			s.Logger.Error("Failed to close response body", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return JWKS{}, time.Time{}, fmt.Errorf("jwks fetch status %d", resp.StatusCode)
	}
	var jwks JWKS
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&jwks); err != nil {
		return JWKS{}, time.Time{}, err
	}
	expiry := calcJWKSExpiry(resp.Header)
	return jwks, expiry, nil
}

func calcJWKSExpiry(h http.Header) time.Time {
	now := time.Now()
	if cc := h.Get("Cache-Control"); cc != "" {
		// look for max-age=...
		if idx := strings.Index(cc, "max-age="); idx != -1 {
			rest := cc[idx+len("max-age="):]
			secs := 0
			for _, ch := range rest {
				if ch < '0' || ch > '9' {
					break
				}
				secs = secs*10 + int(ch-'0')
			}
			if secs > 0 {
				return now.Add(time.Duration(secs) * time.Second)
			}
		}
	}
	if exp := h.Get("Expires"); exp != "" {
		if t, err := http.ParseTime(exp); err == nil {
			return t
		}
	}
	return now.Add(defaultJWKSCache)
}

// ---------------- JWK utilities ----------------

func findJWKByKid(jwks JWKS, kid string) *JWK {
	for i := range jwks.Keys {
		if jwks.Keys[i].Kid == kid {
			return &jwks.Keys[i]
		}
	}
	return nil
}

// jwkToPublicRSAWithCertChecks converts JWK to *rsa.PublicKey.
// Preference: n/e fields. Fallback: x5c[0] cert with expiry checks.
func jwkToPublicRSAWithCertChecks(jwk JWK) (*rsa.PublicKey, error) {
	// prefer n/e
	if jwk.N != "" && jwk.E != "" {
		nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			return nil, fmt.Errorf("invalid jwk n: %w", err)
		}
		eb, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			return nil, fmt.Errorf("invalid jwk e: %w", err)
		}
		n := new(big.Int).SetBytes(nb)

		// e usually small; convert big-endian bytes to int
		eInt := 0
		for _, b := range eb {
			eInt = eInt<<8 + int(b)
		}
		if eInt == 0 {
			return nil, errors.New("invalid exponent in jwk")
		}
		return &rsa.PublicKey{N: n, E: eInt}, nil
	}

	// fallback to x5c (first cert)
	if len(jwk.X5c) > 0 {
		der, err := base64.StdEncoding.DecodeString(jwk.X5c[0])
		if err != nil {
			return nil, fmt.Errorf("invalid x5c cert encoding: %w", err)
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return nil, fmt.Errorf("parse x5c certificate: %w", err)
		}
		// basic validity check (expiry)
		now := time.Now()
		if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
			return nil, errors.New("x5c certificate is not currently valid")
		}
		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("x5c certificate does not contain RSA public key")
		}
		return pub, nil
	}

	return nil, ErrNoKeyMaterial
}

// ---------------- Claims / time / aud helpers ----------------

func decodeJSONToMap(b []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func toInt64FromAny(v any) (int64, bool) {
	switch t := v.(type) {
	case float64:
		return int64(t), true
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			// try parse as float then cast
			if f, ferr := t.Float64(); ferr == nil {
				return int64(f), true
			}
			return 0, false
		}
		return i, true
	case int64:
		return t, true
	case int:
		return int64(t), true
	default:
		return 0, false
	}
}

func validateTimeClaimWithLeeway(claims map[string]any, name string, now int64, leeway int64) error {
	v, ok := claims[name]
	if !ok {
		return nil
	}
	val, ok := toInt64FromAny(v)
	if !ok {
		return fmt.Errorf("invalid %s claim type", name)
	}
	switch name {
	case "exp":
		if now > val+leeway {
			return fmt.Errorf("token %s validation failed", name)
		}
	case "nbf":
		if now < val-leeway {
			return fmt.Errorf("token %s validation failed", name)
		}
	default:
		// default numeric check
		if now > val+leeway {
			return fmt.Errorf("token %s validation failed", name)
		}
	}
	return nil
}

func verifyAudience(aud any, want string) error {
	if aud == nil {
		return ErrInvalidAudience
	}
	switch t := aud.(type) {
	case string:
		if t != want {
			return ErrInvalidAudience
		}
	case []any:
		for _, e := range t {
			if s, ok := e.(string); ok && s == want {
				return nil
			}
		}
		return ErrInvalidAudience
	case []string:
		for _, s := range t {
			if s == want {
				return nil
			}
		}
		return ErrInvalidAudience
	default:
		return ErrInvalidAudience
	}

	return nil
}

// ---------------- Utility ----------------

/*
If you want to prune cache you could add a background goroutine that evicts expired items,
or expose a method to clear entries. For most uses current lazy-cache is fine.
*/

// Optional helper: clear cache entry (by jwks_uri)
func (s *Service) ClearJWKS(jwksURI string) {
	s.mu.Lock()
	delete(s.cache, jwksURI)
	s.mu.Unlock()
}
