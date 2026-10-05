package dreego

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
)

// ErrSecretTooShort is returned by NewCookieStore when the secret is too short.
var ErrSecretTooShort = errors.New("dreego: Session-Secret muss mindestens 32 Bytes haben")

// ErrSessionTooLarge is returned when the encoded session exceeds the cookie
// size limit (4096 bytes).
var ErrSessionTooLarge = errors.New("dreego: Session ist zu groß für ein Cookie")

// CookiePolicy holds the cookie defaults for the session store. The zero value
// is unsafe; use DefaultCookiePolicy or SetCookiePolicy.
type CookiePolicy struct {
	Name     string
	Path     string
	HttpOnly bool
	SameSite http.SameSite
	// Secure marks the cookie Secure. It is OR-ed with the request's TLS state,
	// so a TLS request always gets a Secure cookie.
	Secure bool
}

// DefaultCookiePolicy returns secure defaults: HttpOnly, SameSite=Lax, Path=/.
func DefaultCookiePolicy() CookiePolicy {
	return CookiePolicy{
		Name:     "dreego_session",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// CookieStore keeps the session in an AES-256-GCM encrypted cookie. Encryption
// protects both confidentiality and integrity (GCM authenticates), so a tampered
// cookie is rejected on read.
//
// The store is stateless: a stolen cookie can be replayed until it expires or
// the secret rotates. Immediate invalidation needs a server-side store.
type CookieStore struct {
	encKey []byte
	policy CookiePolicy

	// trustedProxies are the addresses allowed to assert X-Forwarded-Proto.
	// Behind a TLS-terminating proxy (Traefik), r.TLS is nil in the app, so
	// Secure can only be set from the forwarded header — and only from a proxy
	// we trust, never from an arbitrary client.
	trustedProxies map[string]bool
}

// SetTrustedProxies configures the proxy addresses allowed to set Secure via
// X-Forwarded-Proto. Without this, a request is Secure only over direct TLS.
func (store *CookieStore) SetTrustedProxies(adressen []string) {
	if len(adressen) == 0 {
		store.trustedProxies = nil
		return
	}
	menge := map[string]bool{}
	for _, adresse := range adressen {
		menge[adresse] = true
	}
	store.trustedProxies = menge
}

// NewCookieStore builds the store from a secret of at least 32 bytes. It
// derives a 256-bit AES key via HMAC-SHA256, so any secret length >= 32 works.
func NewCookieStore(secret []byte) (*CookieStore, error) {
	if len(secret) < 32 {
		return nil, ErrSecretTooShort
	}

	ableitung := hmac.New(sha256.New, secret)
	ableitung.Write([]byte("dreego-session-enc"))
	schluessel := ableitung.Sum(nil)

	return &CookieStore{encKey: schluessel, policy: DefaultCookiePolicy()}, nil
}

// SetCookiePolicy merges a policy with the secure defaults, so a partial policy
// cannot silently drop HttpOnly, Path or the cookie name.
func (store *CookieStore) SetCookiePolicy(policy CookiePolicy) {
	zusammen := store.policy

	if policy.Name != "" {
		zusammen.Name = policy.Name
	}
	if policy.Path != "" {
		zusammen.Path = policy.Path
	}
	if policy.SameSite != 0 {
		zusammen.SameSite = policy.SameSite
	}
	// HttpOnly und Secure können nur ANgeschaltet werden, nie ab.
	zusammen.HttpOnly = store.policy.HttpOnly || policy.HttpOnly
	zusammen.Secure = store.policy.Secure || policy.Secure

	store.policy = zusammen
}

// Policy returns the active cookie policy.
func (store *CookieStore) Policy() CookiePolicy { return store.policy }

// Load reads the whole session into a map (empty if missing or corrupt).
func (store *CookieStore) Load(r *http.Request) map[string]string {
	daten, _ := store.load(r)
	return daten
}

// Save encrypts and writes the whole session as one cookie.
func (store *CookieStore) Save(w http.ResponseWriter, r *http.Request, daten map[string]string) error {
	return store.write(w, r, daten)
}

// Get reads a single value from the session cookie. A missing, corrupt or
// tampered cookie yields an error; callers treat that as "empty session".
func (store *CookieStore) Get(r *http.Request, key string) (string, error) {
	daten, err_lesen := store.load(r)
	if err_lesen != nil {
		return "", err_lesen
	}
	return daten[key], nil
}

// Set writes a single value into the session cookie (empty value deletes it).
func (store *CookieStore) Set(w http.ResponseWriter, r *http.Request, key string, val string) error {
	daten, _ := store.load(r)
	if daten == nil {
		daten = map[string]string{}
	}

	if val == "" {
		delete(daten, key)
	} else {
		daten[key] = val
	}

	return store.write(w, r, daten)
}

// Delete removes one key from the session.
func (store *CookieStore) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	return store.Set(w, r, key, "")
}

// Destroy clears the whole session and expires the cookie (logout).
func (store *CookieStore) Destroy(w http.ResponseWriter, r *http.Request) error {
	daten, _ := store.load(r)
	for key := range daten {
		delete(daten, key)
	}
	return store.write(w, r, daten)
}

// load decrypts the session cookie into a map. Missing/corrupt → empty map.
func (store *CookieStore) load(r *http.Request) (map[string]string, error) {
	cookie, err_cookie := r.Cookie(store.policy.Name)
	if err_cookie != nil {
		return map[string]string{}, nil
	}

	daten, err_decode := store.decode(cookie.Value)
	if err_decode != nil {
		return map[string]string{}, err_decode
	}

	return daten, nil
}

// write encrypts the map and sets the cookie.
func (store *CookieStore) write(w http.ResponseWriter, r *http.Request, daten map[string]string) error {
	wert, err_encode := store.encode(daten)
	if err_encode != nil {
		return err_encode
	}

	if len(wert) > 4096 {
		return ErrSessionTooLarge
	}

	http.SetCookie(w, &http.Cookie{
		Name:     store.policy.Name,
		Value:    wert,
		Path:     store.policy.Path,
		HttpOnly: store.policy.HttpOnly,
		SameSite: store.policy.SameSite,
		Secure:   store.isSecure(r),
	})

	return nil
}

// isSecure decides whether the cookie is marked Secure:
//
//   - true when the policy demands it, or
//   - true over a direct TLS connection (r.TLS != nil), or
//   - true when the request arrives from a TRUSTED proxy that asserts
//     X-Forwarded-Proto: https.
//
// An untrusted client cannot force Secure by sending the header.
func (store *CookieStore) isSecure(r *http.Request) bool {
	if store.policy.Secure {
		return true
	}
	if r.TLS != nil {
		return true
	}
	if store.trustedProxies[clientHost(r.RemoteAddr)] {
		return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	}
	return false
}

// clientHost strips the port from a RemoteAddr ("10.0.0.1:5432" -> "10.0.0.1").
func clientHost(remoteAddr string) string {
	host, _, err_split := net.SplitHostPort(remoteAddr)
	if err_split != nil {
		return remoteAddr
	}
	return host
}

// encode marshals the map and seals it with AES-256-GCM.
func (store *CookieStore) encode(daten map[string]string) (string, error) {
	roh, err_json := json.Marshal(daten)
	if err_json != nil {
		return "", err_json
	}

	block, err_block := aes.NewCipher(store.encKey)
	if err_block != nil {
		return "", err_block
	}

	gcm, err_gcm := cipher.NewGCM(block)
	if err_gcm != nil {
		return "", err_gcm
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err_nonce := rand.Read(nonce)
	if err_nonce != nil {
		return "", err_nonce
	}

	versiegelt := gcm.Seal(nonce, nonce, roh, nil)
	return base64.RawURLEncoding.EncodeToString(versiegelt), nil
}

// decode reverses encode; any tampering or wrong key fails the GCM check.
func (store *CookieStore) decode(wert string) (map[string]string, error) {
	versiegelt, err_b64 := base64.RawURLEncoding.DecodeString(wert)
	if err_b64 != nil {
		return nil, err_b64
	}

	block, err_block := aes.NewCipher(store.encKey)
	if err_block != nil {
		return nil, err_block
	}

	gcm, err_gcm := cipher.NewGCM(block)
	if err_gcm != nil {
		return nil, err_gcm
	}

	groesse := gcm.NonceSize()
	if len(versiegelt) < groesse {
		return nil, errors.New("dreego: Cookie zu kurz")
	}

	nonce, ciphertext := versiegelt[:groesse], versiegelt[groesse:]
	roh, err_oeffnen := gcm.Open(nil, nonce, ciphertext, nil)
	if err_oeffnen != nil {
		return nil, err_oeffnen
	}

	daten := map[string]string{}
	err_json := json.Unmarshal(roh, &daten)
	if err_json != nil {
		return nil, err_json
	}

	return daten, nil
}
