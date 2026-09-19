// Command signmodels is the C29 manifest signing tool (ticket 14). It generates
// minisign-format keypairs, signs models/manifest.json into
// models/manifest.minisig, and cross-verifies - the runtime verifier in
// internal/models accepts exactly this format (real minisign CLI compatible:
// the verifier also accepts pre-hashed "ED" signatures; this tool emits the
// legacy "Ed" form, which minisign verifies by default).
//
// Key custody (D33/F3, C29):
//   - the SECRET key never enters the repository. It lives out-of-repo
//     (default E:\work\base\wisp-minisign\), injected via -key or the
//     WISP_MINISIGN_KEY environment variable (scripts/sign-models.ps1 wraps
//     this). The dev key is an UNENCRYPTED dev-only envelope; production keys
//     are a real minisign ceremony at S8 and MUST rotate (docs/reports).
//   - the public key is hardcoded into internal/buildinfo (MinisignPublicKey).
//
// Usage:
//
//	signmodels keygen  -dir <dir>                       # wisp-models.key + wisp-models.pub
//	signmodels sign    -key <key> [-in manifest.json] [-out manifest.minisig]
//	signmodels verify  -pub <pub|buildinfo> -in f -sig s
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/models"
)

const devKeyComment = "wisp dev minisign secret key (PLAINTEXT ENVELOPE - dev only; S8 rotates to an encrypted production key)"

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "signmodels: FATAL: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "keygen":
		keygen(os.Args[2:])
	case "sign":
		sign(os.Args[2:])
	case "verify":
		verifyCmd(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage:
  signmodels keygen -dir <out-dir>
  signmodels sign   -key <secret-key> [-in models/manifest.json] [-out models/manifest.minisig]
  signmodels verify -pub <public-key-file> -in <file> -sig <signature>
`)
	os.Exit(2)
}

// keygen writes wisp-models.key (secret, 0600) and wisp-models.pub
// (minisign two-line format) into dir.
func keygen(args []string) {
	fs := flag.NewFlagSet("keygen", flag.ExitOnError)
	dir := fs.String("dir", `E:\work\base\wisp-minisign`, "output directory (out of repo)")
	_ = fs.Parse(args)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fatal("keygen: %v", err)
	}
	keyID := keyIDFromSeed(pub)
	pkgPub := models.MinisignPublicKey{KeyID: keyID, Key: pub}

	if err := os.MkdirAll(*dir, 0o700); err != nil {
		fatal("keygen: %v", err)
	}
	secret := struct {
		keyID [8]byte
		seed  []byte
		pub   []byte
	}{keyID, priv.Seed(), pub}
	blob := make([]byte, 0, 2+8+len(priv.Seed())+len(pub))
	blob = append(blob, "Ed"...)
	blob = append(blob, secret.keyID[:]...)
	blob = append(blob, secret.seed...)
	blob = append(blob, secret.pub...)

	keyPath := filepath.Join(*dir, "wisp-models.key")
	pubPath := filepath.Join(*dir, "wisp-models.pub")
	if err := os.WriteFile(keyPath, []byte(
		"untrusted comment: "+devKeyComment+"\n"+
			base64.StdEncoding.EncodeToString(blob)+"\n"), 0o600); err != nil {
		fatal("keygen: %v", err)
	}
	if err := os.WriteFile(pubPath, []byte(pkgPub.String("wisp models signing key (dev)")), 0o644); err != nil {
		fatal("keygen: %v", err)
	}
	fmt.Printf("signmodels: secret key  %s (NEVER commit; 0600)\n", keyPath)
	fmt.Printf("signmodels: public key  %s\n", pubPath)
	fmt.Printf("signmodels: keyid       %s\n", pkgPub.KeyIDHex())
	fmt.Printf("signmodels: paste the SECOND line of the .pub into internal/buildinfo.MinisignPublicKey\n")
}

// sign loads the dev secret envelope, signs the manifest file, writes the
// 4-line minisign signature, then verifies it back through the runtime path.
func sign(args []string) {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)
	keyPath := fs.String("key", os.Getenv("WISP_MINISIGN_KEY"), "secret key path (env WISP_MINISIGN_KEY)")
	in := fs.String("in", filepath.Join("models", "manifest.json"), "manifest to sign")
	out := fs.String("out", *in+".minisig", "signature output")
	_ = fs.Parse(args)
	if *keyPath == "" {
		fatal("no -key and no WISP_MINISIGN_KEY")
	}

	priv, pub, err := loadSecretKey(*keyPath)
	if err != nil {
		fatal("%v", err)
	}
	payload, err := os.ReadFile(*in)
	if err != nil {
		fatal("%v", err)
	}
	sig, err := models.SignPayload(priv, pub, payload,
		"signature from wisp models key "+pub.KeyIDHex(),
		fmt.Sprintf("timestamp:%d\tfile:%s", time.Now().Unix(), filepath.Base(*in)))
	if err != nil {
		fatal("%v", err)
	}
	if err := os.WriteFile(*out, sig, 0o644); err != nil {
		fatal("%v", err)
	}
	// Self-check through the exact runtime verifier.
	if err := models.VerifyMinisignSignature(pub, payload, sig); err != nil {
		fatal("post-sign verification failed: %v", err)
	}
	fmt.Printf("signmodels: signed %s -> %s (keyid %s, runtime-verified)\n", *in, *out, pub.KeyIDHex())
}

// verifyCmd checks a signature against a public key file - the script's
// cross-check command (no signing rights needed).
func verifyCmd(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	pubPath := fs.String("pub", "", "public key file")
	in := fs.String("in", "", "signed file")
	sigPath := fs.String("sig", "", "signature file")
	_ = fs.Parse(args)
	if *pubPath == "" || *in == "" || *sigPath == "" {
		usage()
	}
	pubRaw, err := os.ReadFile(*pubPath)
	if err != nil {
		fatal("%v", err)
	}
	pub, err := models.ParseMinisignPublicKey(string(pubRaw))
	if err != nil {
		fatal("%v", err)
	}
	payload, err := os.ReadFile(*in)
	if err != nil {
		fatal("%v", err)
	}
	sig, err := os.ReadFile(*sigPath)
	if err != nil {
		fatal("%v", err)
	}
	if err := models.VerifyMinisignSignature(pub, payload, sig); err != nil {
		fatal("VERIFICATION FAILED: %v", err)
	}
	fmt.Printf("signmodels: OK %s (keyid %s)\n", *in, pub.KeyIDHex())
}

// loadSecretKey parses the dev plaintext envelope ("Ed" || keyid || seed || pub).
func loadSecretKey(path string) (ed25519.PrivateKey, models.MinisignPublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, models.MinisignPublicKey{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 2 || !strings.HasPrefix(lines[0], "untrusted comment: ") {
		return nil, models.MinisignPublicKey{}, fmt.Errorf("%s: not a wisp dev secret envelope", path)
	}
	blob, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[1]))
	if err != nil {
		return nil, models.MinisignPublicKey{}, fmt.Errorf("%s: base64: %w", path, err)
	}
	if len(blob) != 2+8+ed25519.SeedSize+ed25519.PublicKeySize || string(blob[:2]) != "Ed" {
		return nil, models.MinisignPublicKey{}, fmt.Errorf("%s: malformed envelope", path)
	}
	seed := blob[10 : 10+ed25519.SeedSize]
	pubBytes := blob[10+ed25519.SeedSize:]
	priv := ed25519.NewKeyFromSeed(seed)
	if string(priv.Public().(ed25519.PublicKey)) != string(pubBytes) {
		return nil, models.MinisignPublicKey{}, fmt.Errorf("%s: envelope pubkey does not match seed", path)
	}
	var keyID [8]byte
	copy(keyID[:], blob[2:10])
	return priv, models.MinisignPublicKey{KeyID: keyID, Key: pubBytes}, nil
}

// keyIDFromSeed derives the key id the same way internal/models does for
// parsed keys: first 8 bytes of SHA-256 over the public key.
func keyIDFromSeed(pub ed25519.PublicKey) [8]byte {
	return models.DeriveKeyID(pub)
}
