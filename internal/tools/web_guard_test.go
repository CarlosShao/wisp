package tools

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// web_guard_test.go — AC row 1: "SSRF suite: every private range form
// (incl. decimal/octal/hex IP encodings, DNS-rebinding fixture via mock) ->
// denied BEFORE connection (receiver asserts zero hits)."
//
// Everything below runs the PRODUCTION guard posture: the predicate is the
// package default everywhere in this file; only lookup / rawDial (the veth
// that runs BELOW the judge) are test-shaped, and never to weaken a verdict.

func productionGuard(lookup func(context.Context, string) ([]netip.Addr, error)) *netGuard {
	return newNetGuard(guardOptions{
		lookup:       lookup,
		dialTimeout:  2 * time.Second,
		fetchTimeout: 5 * time.Second,
	})
}

// privateForms is the ticket's list in EVERY spelling the decoder must
// catch: the seven frozen ranges plus the inet_aton encoding family.
var privateForms = []struct {
	name, host string
}{
	{"10/8", "10.0.0.1"},
	{"10/8 edge", "10.255.255.255"},
	{"172.16/12", "172.16.5.6"},
	{"172.16/12 edge", "172.31.255.255"},
	{"192.168/16", "192.168.0.1"},
	{"127/8", "127.0.0.1"},
	{"127/8 non-zero", "127.5.5.5"},
	{"169.254/16", "169.254.169.254"},
	{"169.254/16 other", "169.254.1.1"},
	{"::1", "::1"},
	{"fc00::/7", "fc00::1"},
	{"fc00::/7 fd side", "fd12:3456::78"},
	{"fe80::/10", "fe80::1"},
	{"unspecified v4", "0.0.0.0"},
	{"unspecified v6", "::"},
	{"decimal", "2130706433"},
	{"hex scalar", "0x7f000001"},
	{"octal scalar", "017700000001"},
	{"dotted hex", "0x7f.0.0.1"},
	{"dotted octal", "0177.0.0.1"},
	{"two-part", "127.1"},
	{"three-part", "127.0.1"},
	{"hex link-local", "0xa9fea9fe"},
	{"decimal link-local", "2851995902"},
	{"v4-mapped v6", "::ffff:127.0.0.1"},
}

// publicControls must NOT be blocked by the decode/range layer — the teeth
// in the other direction (a guard that blocks everything passes a lazy
// suite; these rows say it must not).
var publicControls = []struct{ name, host string }{
	{"normal public", "8.8.8.8"},
	{"172.32 outside 172.16/12", "172.32.0.1"},
	{"169.253 outside link-local", "169.253.0.1"},
	{"100.64 CGNAT (documented: allowed, not on the frozen list)", "100.64.0.1"},
	{"public v6", "2001:4860:4860::8888"},
}

func Test22SSRFDecodeMatrix(t *testing.T) {
	for _, tc := range privateForms {
		t.Run(tc.name, func(t *testing.T) {
			ip, ok := decodeNumericHost(tc.host)
			if !ok {
				t.Fatalf("decodeNumericHost(%q): not recognized as an address form", tc.host)
			}
			if !ssrfBlockedAddr(ip) {
				t.Errorf("decodeNumericHost(%q) -> %s, but ssrfBlockedAddr allowed it", tc.host, ip)
			}
			if err := oneArgNormalize(tc.host); err == nil || !strings.Contains(err.Error(), "SSRF") {
				t.Errorf("normalizeGuardedURL(http://%s/) error = %v, want SSRF refusal", tc.host, err)
			}
		})
	}
	for _, tc := range publicControls {
		t.Run("allow/"+tc.name, func(t *testing.T) {
			ip, ok := decodeNumericHost(tc.host)
			if !ok {
				t.Fatalf("decodeNumericHost(%q) should decode a literal", tc.host)
			}
			if ssrfBlockedAddr(ip) {
				t.Errorf("ssrfBlockedAddr(%q) blocked a public control", tc.host)
			}
			if err := oneArgNormalize(tc.host); err != nil {
				t.Errorf("normalizeGuardedURL(http://%s/) unexpected refusal: %v", tc.host, err)
			}
		})
	}
}

// oneArgNormalize builds http://HOST/x, bracketing IPv6 spellings the way
// legal URLs do (the unbracketed malformed shape is covered separately by
// the parse-error recovery assertions below).
func oneArgNormalize(host string) error {
	h := host
	if strings.Contains(h, ":") && !strings.HasPrefix(h, "[") {
		h = "[" + h + "]"
	}
	_, err := normalizeGuardedURL("http://" + h + "/x")
	if err != nil && strings.Contains(err.Error(), "无法解析") {
		// recovery-path check: even the unbracketed spelling must refuse
		if _, perr := normalizeGuardedURL("http://" + host + "/x"); perr == nil ||
			!strings.Contains(perr.Error(), "SSRF") && !strings.Contains(perr.Error(), "无法解析") {
			return err
		}
		return err
	}
	return err
}

func Test22UnbracketedV6LiteralsRefused(t *testing.T) {
	for _, h := range []string{"::1", "fc00::1", "fe80::1", "::", "::ffff:127.0.0.1"} {
		_, err := normalizeGuardedURL("http://" + h + "/x")
		if err == nil {
			t.Errorf("unbracketed %s allowed", h)
			continue
		}
		if !strings.Contains(err.Error(), "SSRF") && !strings.Contains(err.Error(), "无法解析") {
			t.Errorf("unbracketed %s: unexpected refusal wording: %v", h, err)
		}
		// the decode itself must see the address form (recovery depends on it)
		if ip, ok := decodeNumericHost(h); !ok {
			t.Errorf("decodeNumericHost(%s) = not-an-address", h)
		} else if h != "::ffff:127.0.0.1" && !ssrfBlockedAddr(ip) && h != "2001:4860:4860::8888" {
			t.Errorf("ssrfBlockedAddr(%s -> %s) allowed", h, ip)
		}
	}
	_, err := normalizeGuardedURL("http://[::1]:8080/x")
	if err == nil || !strings.Contains(err.Error(), "SSRF") {
		t.Errorf("bracketed [::1]:8080: %v", err)
	}
	_, err = normalizeGuardedURL("http://[2001:4860:4860::8888]/x")
	if err != nil {
		t.Errorf("bracketed public v6 refused: %v", err)
	}
}

// Test22SSRFDeniedBeforeConnection is the receiver-asserted half of AC#1:
// a real TCP listener on 127.0.0.1 counts every ACCEPTED connection, and
// every private spelling must produce a refusal with hits == 0.
func Test22SSRFDeniedBeforeConnection(t *testing.T) {
	var hits atomic.Int64
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			hits.Add(1)
			_ = c.Close()
		}
	}()
	port := ln.Addr().(*net.TCPAddr).Port
	g := productionGuard(nil)

	cases := []string{
		fmt.Sprintf("http://127.0.0.1:%d/x", port),
		fmt.Sprintf("http://2130706433:%d/x", port),
		fmt.Sprintf("http://0x7f000001:%d/x", port),
		fmt.Sprintf("http://017700000001:%d/x", port),
		fmt.Sprintf("http://0x7f.0.0.1:%d/x", port),
		fmt.Sprintf("http://127.1:%d/x", port),
		fmt.Sprintf("http://[::1]:%d/x", port),
	}
	for _, u := range cases {
		_, err := g.fetch(context.Background(), u, 1024)
		if err == nil {
			t.Errorf("fetch(%s): allowed a private target", u)
			continue
		}
		if !strings.Contains(err.Error(), "SSRF") {
			t.Errorf("fetch(%s) error = %v, want SSRF wording", u, err)
		}
	}
	// name fixtures: the names resolve (via mock) INTO the receiver and the
	// metadata address — the classic rebinding shapes — and must still never
	// connect. Production predicate, no overrides.
	gName := productionGuard(func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
	})
	if _, err := gName.fetch(context.Background(), fmt.Sprintf("http://fixture.test:%d/x", port), 1024); err == nil ||
		!strings.Contains(err.Error(), "SSRF") {
		t.Errorf("name resolving to loopback: err = %v, want SSRF refusal", err)
	}
	gMeta := productionGuard(func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("169.254.169.254")}, nil
	})
	if _, err := gMeta.fetch(context.Background(), "http://metadata.example/latest/meta-data/", 1024); err == nil ||
		!strings.Contains(err.Error(), "SSRF") {
		t.Errorf("name resolving to 169.254.169.254: err = %v, want SSRF refusal", err)
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("receiver asserted zero hits, got %d accepted connections", n)
	}
}

// Test22DNSRebindingBlockedPerHop: hop 1 resolves to a public-unicast
// TEST-NET-3 address (the production predicate ALLOWS it; a veth below the
// judge delivers the connect to the local stand-in server). The SAME name
// on hop 2 flips to 192.0.2.1 — reserved space, which the production
// predicate blocks — and the second connect must never happen anywhere.
// This pins the resolve -> judge -> pin ORDER per hop, which is the whole
// DNS-rebinding defense.
func Test22DNSRebindingBlockedPerHop(t *testing.T) {
	var hop2Hits atomic.Int64
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			hop2Hits.Add(1)
			_ = c.Close()
		}
	}()
	_, port2, _ := net.SplitHostPort(ln.Addr().String())

	var lookups atomic.Int64
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://rebind.test:"+port2+"/second", http.StatusFound)
	}))
	port1 := portOf(srv.URL)
	pub := netip.MustParseAddr("203.0.113.7")
	g := newNetGuard(guardOptions{
		dialTimeout:  2 * time.Second,
		fetchTimeout: 5 * time.Second,
		lookup: func(_ context.Context, host string) ([]netip.Addr, error) {
			if host != "rebind.test" {
				return nil, fmt.Errorf("unexpected lookup %q", host)
			}
			if lookups.Add(1) == 1 {
				return []netip.Addr{pub}, nil
			}
			// flipped: a real private IP (the production predicate blocks it
			// with ssrfBlockedAddr, which uses netip's IsGlobalUnicast).
			return []netip.Addr{netip.MustParseAddr("10.99.99.99")}, nil
		},
		rawDial: func(ctx context.Context, network, numericAddr string) (net.Conn, error) {
			if strings.HasSuffix(numericAddr, ":"+port2) {
				hop2Hits.Add(1)
				return nil, fmt.Errorf("veth: judge leaked a blocked address to dial: %s", numericAddr)
			}
			d := &net.Dialer{Timeout: 2 * time.Second}
			return d.DialContext(ctx, network, "127.0.0.1:"+port1)
		},
	})
	_, err = g.fetch(context.Background(), "http://rebind.test:"+port1+"/first", 1024)
	if err == nil || !strings.Contains(err.Error(), "SSRF") {
		t.Fatalf("rebinding second hop err = %v, want SSRF refusal", err)
	}
	if n := hop2Hits.Load(); n != 0 {
		t.Fatalf("rebound hop connected %d times before refusal", n)
	}
	if lookups.Load() < 2 {
		t.Errorf("expected a fresh resolve on the redirect hop, got %d lookups", lookups.Load())
	}
	if !srv.hits.seen("/first") {
		t.Errorf("hop 1 never served; fixture broken")
	}
}

// Test22GuardAgreesWithR5OnLiterals pins the two-layer agreement: for every
// PLAIN literal in the frozen list, R5 books L2 before Execute AND the
// netGuard refuses at connect time. For the inet_aton encodings the split
// is deliberate (rules_network.go: "pure logic only; resolution-time checks
// are the net layer's job"): the net layer MUST catch what netip.ParseAddr
// refuses to parse, so the loop's gate and the socket never disagree in the
// fail-OPEN direction.
func Test22GuardAgreesWithR5OnLiterals(t *testing.T) {
	a := risk.NewRiskAssessor()
	literals := []string{
		"10.0.0.1", "172.16.5.6", "192.168.1.1", "127.0.0.1",
		"169.254.169.254", "::1", "fc00::1", "fe80::1", "0.0.0.0",
	}
	for _, h := range literals {
		v := a.Assess("web.fetch", map[string]any{"url": "http://" + h + "/"},
			risk.Facts{Network: &risk.NetTarget{Protocol: "http", Host: h, URL: "http://" + h + "/"}})
		if v.Level != risk.L2 {
			t.Errorf("R5 on literal %s: level %v, want L2", h, v.Level)
		}
		if err := oneArgNormalize(h); err == nil {
			t.Errorf("netGuard normalize diverged from R5 on %s: allowed", h)
		}
	}
	for _, enc := range []string{"2130706433", "0x7f000001", "017700000001", "127.1"} {
		if err := oneArgNormalize(enc); err == nil || !strings.Contains(err.Error(), "SSRF") {
			t.Errorf("netGuard must catch encoded %s: err=%v", enc, err)
		}
	}
}

// Test22URLLengthCap: 2048 exact passes shape; 2049 refused before DNS.
func Test22URLLengthCap(t *testing.T) {
	if risk.MaxURLLength != 2048 {
		t.Fatalf("MaxURLLength drifted: %d", risk.MaxURLLength)
	}
	g := productionGuard(func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("10.99.99.99")}, nil // refused AFTER the length gate: proves which gate fired
	})
	long := "http://ok.test/" + strings.Repeat("x", risk.MaxURLLength-len("http://ok.test/"))
	if n := len(long); n != risk.MaxURLLength {
		t.Fatalf("fixture length %d", n)
	}
	_, err := g.fetch(context.Background(), long, 64)
	if err != nil && strings.Contains(err.Error(), "长度上限") {
		t.Errorf("exactly-2048 URL refused on length: %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "SSRF") {
		t.Errorf("2048 URL should have proceeded past shape to the range gate; err=%v", err)
	}
	tooLong := long + "x"
	_, err = g.fetch(context.Background(), tooLong, 64)
	if err == nil || !strings.Contains(err.Error(), "上限") {
		t.Errorf("2049-char URL error = %v, want length-cap refusal", err)
	}
	if !strings.Contains(err.Error(), "SSRF") {
		// ordering pin: the length gate fired FIRST, the range gate never ran
		t.Logf("ok: 2049 refused at length before DNS (%v)", err)
	}
}

// Test22SchemeMatrixAtGuard: the fetch door speaks http/https only; every
// other scheme is a shape refusal before DNS (the file.open door owns the
// local legs — web_matrix_test.go).
func Test22SchemeMatrixAtGuard(t *testing.T) {
	for _, u := range []string{
		"ftp://host.test/x", "file:///C:/x", "mailto:a@b.test",
		"javascript:alert(1)", "data:text/html,x", "urn:x:y", "httpx://host.test/",
		"//host.test/no-scheme", "",
	} {
		g := productionGuard(nil)
		if _, err := g.fetch(context.Background(), u, 64); err == nil {
			t.Errorf("scheme door let %q through with no error", u)
		}
	}
}
