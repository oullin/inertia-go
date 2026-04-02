package cryptox_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/oullin/inertia-go/core/cryptox"
)

func testKey(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	return key
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	key := testKey(t)
	plaintext := "hello-csrf-token-value"

	encrypted, err := cryptox.Encrypt(plaintext, key)

	if err != nil {
		t.Fatal(err)
	}

	if encrypted == plaintext {
		t.Error("encrypted should differ from plaintext")
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypt_ProducesDifferentOutputs(t *testing.T) {
	key := testKey(t)

	a, _ := cryptox.Encrypt("same", key)
	b, _ := cryptox.Encrypt("same", key)

	if a == b {
		t.Error("two encryptions of the same plaintext should differ (random IV)")
	}
}

func TestDecrypt_TamperedMAC(t *testing.T) {
	key := testKey(t)

	encrypted, err := cryptox.Encrypt("token", key)

	if err != nil {
		t.Fatal(err)
	}

	// Decode, tamper with MAC, re-encode.
	js, _ := base64.StdEncoding.DecodeString(encrypted)

	var p struct {
		IV    string `json:"iv"`
		Value string `json:"value"`
		MAC   string `json:"mac"`
	}

	json.Unmarshal(js, &p)

	p.MAC = strings.Repeat("ff", 32)

	tampered, _ := json.Marshal(p)
	tamperedB64 := base64.StdEncoding.EncodeToString(tampered)

	_, err = cryptox.Decrypt(tamperedB64, key)

	if err == nil {
		t.Error("expected error for tampered MAC")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := testKey(t)
	key2 := testKey(t)

	encrypted, _ := cryptox.Encrypt("token", key1)

	_, err := cryptox.Decrypt(encrypted, key2)

	if err == nil {
		t.Error("expected error for wrong key")
	}
}

func TestDecrypt_MalformedBase64(t *testing.T) {
	key := testKey(t)

	_, err := cryptox.Decrypt("not-valid-base64!!!", key)

	if err == nil {
		t.Error("expected error for malformed base64")
	}
}

func TestDecrypt_MalformedJSON(t *testing.T) {
	key := testKey(t)
	bad := base64.StdEncoding.EncodeToString([]byte("{invalid json"))

	_, err := cryptox.Decrypt(bad, key)

	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestDecrypt_EmptyPayload(t *testing.T) {
	key := testKey(t)

	_, err := cryptox.Decrypt("", key)

	if err == nil {
		t.Error("expected error for empty payload")
	}
}

func TestEncryptDecrypt_LongPlaintext(t *testing.T) {
	key := testKey(t)
	plaintext := strings.Repeat("a", 1024)

	encrypted, err := cryptox.Encrypt(plaintext, key)

	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Error("round-trip failed for long plaintext")
	}
}

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	key := testKey(t)

	encrypted, err := cryptox.Encrypt("", key)

	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := cryptox.Decrypt(encrypted, key)

	if err != nil {
		t.Fatal(err)
	}

	if decrypted != "" {
		t.Errorf("decrypted = %q, want empty", decrypted)
	}
}
