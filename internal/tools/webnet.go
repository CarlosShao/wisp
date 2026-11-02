package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// webnet.go is the D30 layer-3 outbound guard (ticket 22): the net layer
// internal/risk/rules_network.go hands its resolution-time checks to ("DNS
// names are NOT resolved here — pure logic only; resolution-time checks are
// the net layer's job"). R5 judges the verdict BEFORE execution; this file
// enforces the same frozen list at connect time, which is the only place a
// DNS-rebound name can be caught:
//
//   - scheme must be http/https and the URL at most risk.MaxURLLength (2048);
//   - every IP form — literal v4/v6, decimal, octal, hex, mixed dotted
//     inet_aton spellings (2130706433, 0x7f000001, 017700000001, 0x7f.0.0.1,
//     127.1) — is DECODED and range-checked before any socket is opened,
//     because url.Parse hands those to DNS as plain names and a resolver
//     would happily reconnect them back into the private space;
//   - DNS results are checked per address and the dial is pinned to one
//     CHECKED address (net.Dialer is never handed the hostname), so a second
//     lookup flipping to 169.254.169.254 cannot race in (rebinding);
//   - redirects: never followed to a non-http(s) scheme, never downgraded
//     https->http, each hop re-passes the whole guard (SPEC-06 §6.3).
//
// The blocked-range predicate mirrors risk.ruleNetwork's isPrivateHost over
// the same netip standard-library verdicts — IsPrivate (10/8, 172.16/12,
// 192.168/16, fc00::/7) plus !IsGlobalUnicast (127/8, 169.254/16, ::1,
// fe80::/10, 0.0.0.0, reserved/multicast). A single-source-of-truth refactor
// of the predicate itself is risk's; the equality of the two verdicts is
// pinned by web_guard_test.go reading both sides.
//
// Response bodies are NEVER logged by anything below: fetched bodies travel
// only inside fetchedDoc.Body (ticket 22 key constraint "body never logged").

const (
	// defaultFetchMaxBytes reuses the spill layer's D15(3) raw hard ceiling
	// at reference scale (internal/agent/budgets.go refRawOutputBytes =
	// 1<<20). A fetch larger than this is cut and marked Truncated; the loop's
	// Spiller then owns spill-vs-truncate as for any tool output.
	defaultFetchMaxBytes = 1 << 20
	defaultDialTimeout   = 10 * time.Second
	defaultFetchTimeout  = 30 * time.Second
	defaultMaxRedirects  = 10
	webUserAgent         = "wisp/1.0 (builtin web.fetch; +https://github.com/CarlosShao/wisp)"
)

// ssrfBlockedAddr reports whether ip sits inside the frozen SSRF list
// (SPEC-06 §6.3): 10/8, 172.16/12, 192.168/16, 127/8, 169.254/16 (incl.
// 169.254.169.254), ::1, fc00::/7, fe80::/10, plus the unspecified address
// and every non-unicast shape. IPv4-mapped v6 (::ffff:x) is unmapped first,
// the only embedded-IPv4 form Go's netip predicates normalize (6to4/NAT64
// spellings are a documented residual — see the ticket-22 evidence file).
func ssrfBlockedAddr(ip netip.Addr) bool {
	if !ip.IsValid() {
		return true
	}
	if ip.Is4In6() {
		ip = ip.Unmap()
	}
	return ip.IsPrivate() || !ip.IsGlobalUnicast()
}

// ssrfBlockedHost judges a URL host string: IP literals and every numeric
// encoding form are decoded and range-checked; plain DNS names return false
// here and are checked at resolution time by the guard dialer.
func ssrfBlockedHost(host string) bool {
	ip, encoded := decodeNumericHost(host)
	if !encoded {
		return false
	}
	return ssrfBlockedAddr(ip)
}

// decodeNumericHost decodes the inet_aton family of host spellings. The bool
// reports "this host IS an address" (literal or encoded); (zero,false) means
// it is a DNS name and the caller must check at resolution time instead.
// Go's own url/dial layer treats all but the plain literals as DNS names,
// which is exactly the bypass class this decode closes before the socket.
func decodeNumericHost(host string) (netip.Addr, bool) {
	host = strings.TrimSpace(host)
	if host == "" {
		return netip.Addr{}, false
	}
	// Plain literals first: v6 keeps its brackets inside URL parsing, and
	// url.Hostname() has already stripped them; accept both spellings here.
	if a, err := netip.ParseAddr(strings.Trim(host, "[]")); err == nil {
		return a, true
	}
	// netip refuses BARE v4-suffix v6 spellings ("::ffff:127.0.0.1"
	// without brackets, which is exactly the shape an unbracketed URL
	// host carries): decode the tail as v4 and lift it into the v6 form.
	// Only the two real spellings qualify (::x.y.z.w and ::ffff:x.y.z.w);
	// anything else stays a malformed host for the generic refusal below.
	if i := strings.LastIndexByte(host, ':'); i >= 0 && !strings.HasPrefix(host, "[") {
		tail := host[i+1:]
		prefix := host[:i]
		if strings.Contains(tail, ".") && (prefix == "::" || strings.EqualFold(prefix, "::ffff")) {
			if a, err := netip.ParseAddr(tail); err == nil && a.Is4() {
				return netip.AddrFrom16(a.As16()), true
			}
		}
	}
	if !strings.Contains(host, ".") {
		// Single-scalar IPv4: decimal (2130706433), hex (0x7f000001),
		// octal (017700000001).
		if v, ok := parseHostNumber(host); ok {
			return v4fromUint(v), true
		}
		return netip.Addr{}, false
	}
	parts := strings.Split(host, ".")
	if len(parts) > 4 {
		return netip.Addr{}, false
	}
	vals := make([]uint32, 0, len(parts))
	for _, p := range parts {
		v, ok := parseHostNumber(p)
		if !ok {
			return netip.Addr{}, false // a letter-bearing label: a DNS name
		}
		vals = append(vals, v)
	}
	// inet_aton tail rule: with fewer than 4 labels the last one absorbs the
	// remaining host bytes (127.1 -> 127.0.0.1, 10.2.5152 -> 10.2.20.0).
	var u uint32
	switch len(vals) {
	case 1:
		u = vals[0]
	case 2:
		if vals[0] > 255 || vals[1] > 0xffffff {
			return netip.Addr{}, false
		}
		u = vals[0]<<24 | vals[1]
	case 3:
		if vals[0] > 255 || vals[1] > 255 || vals[2] > 0xffff {
			return netip.Addr{}, false
		}
		u = vals[0]<<24 | vals[1]<<16 | vals[2]
	case 4:
		if vals[0] > 255 || vals[1] > 255 || vals[2] > 255 || vals[3] > 255 {
			return netip.Addr{}, false
		}
		u = vals[0]<<24 | vals[1]<<16 | vals[2]<<8 | vals[3]
	}
	return v4fromUint(u), true
}

// v4fromUint builds the netip.Addr of a big-endian uint32 IPv4.
func v4fromUint(u uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)})
}

// parseHostNumber parses one inet_aton label: decimal, 0x hex, or leading-0
// octal (base-0 semantics WITHOUT strconv's underscore/plus freedoms). "0"
// alone is decimal zero. A label nothing fits is not numeric.
func parseHostNumber(s string) (uint32, bool) {
	if s == "" {
		return 0, false
	}
	lower := strings.ToLower(s)
	base := 10
	switch {
	case strings.HasPrefix(lower, "0x"):
		base, s = 16, s[2:]
	case len(s) > 1 && s[0] == '0':
		base, s = 8, s[1:]
	}
	if s == "" {
		return 0, false
	}
	if strings.IndexFunc(s, func(r rune) bool { return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') }) >= 0 {
		return 0, false
	}
	v, err := strconv.ParseUint(s, base, 32)
	if err != nil {
		return 0, false
	}
	return uint32(v), true
}

// normalizeGuardedURL applies the shape checks every outbound web tool runs
// BEFORE anything touches DNS or a socket: length cap, scheme matrix, host
// presence, backslash smuggles. It returns the parsed URL; private-range and
// encoded-address verdicts are the caller's guard (blockedHost / dial).
func normalizeGuardedURL(raw string) (*url.URL, error) {
	if len(raw) > risk.MaxURLLength {
		return nil, fmt.Errorf("URL 超过 %d 字符上限（D30③ 出站约束，防把文件内容塞进 query）", risk.MaxURLLength)
	}
	if strings.ContainsAny(raw, `\`) {
		return nil, errors.New("URL 含反斜杠，拒绝（跳转解析混淆防御）")
	}
	u, err := url.Parse(raw)
	if err != nil {
		// url.Parse refuses bare (unbracketed) IPv6 — the SAME shape the
		// old SSRF bypass list targets (http://::1/, http://[::1] is the
		// legal spelling). Recover the authority and let the decoder judge
		// it, so a malformed v6 literal lands on an honest SSRF refusal
		// instead of a generic parse error that hides the intent.
		if host := rawAuthorityHost(raw); host != "" {
			if ip, ok := decodeNumericHost(host); ok && ssrfBlockedAddr(ip) {
				return nil, fmt.Errorf("目标地址属于私有/回环/链路本地网段，连接前即拒绝（SSRF 防御）：%s", host)
			}
		}
		return nil, fmt.Errorf("URL 无法解析：%w", err)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return nil, fmt.Errorf("协议 %q 不在 http/https 之内（D30③；本地协议走 file.open 白名单）", u.Scheme)
	}
	if u.Hostname() == "" {
		return nil, errors.New("URL 缺少主机名")
	}
	if ip, encoded := decodeNumericHost(u.Hostname()); encoded && ssrfBlockedAddr(ip) {
		return nil, fmt.Errorf("目标地址属于私有/回环/链路本地网段，连接前即拒绝（SSRF 防御）：%s", u.Hostname())
	}
	return u, nil
}

// rawAuthorityHost extracts the host part of "scheme://authority" without
// url.Parse's grammar (used only on the parse-error recovery path): cut at
// the first /?#, drop brackets, drop a trailing :port when the port half is
// all-digits (so "::1" keeps its colons, "[::1]:80" loses only the port).
func rawAuthorityHost(raw string) string {
	i := strings.Index(raw, "://")
	if i < 0 {
		return ""
	}
	rest := raw[i+3:]
	if j := strings.IndexAny(rest, "/?#"); j >= 0 {
		rest = rest[:j]
	}
	if strings.HasPrefix(rest, "[") {
		if j := strings.Index(rest, "]"); j >= 0 {
			return rest[1:j]
		}
	}
	if k := strings.LastIndex(rest, ":"); k > 0 {
		if tail := rest[k+1:]; tail != "" && strings.IndexFunc(tail, func(r rune) bool { return r < '0' || r > '9' }) < 0 {
			rest = rest[:k]
		}
	}
	return rest
}

// netGuard owns the client-side enforcement: the dial that resolves, checks
// every answer and then pins to ONE checked IP; the redirect policy that
// re-runs the whole shape check on every hop; and the capped read.
//
// Test seams (lookup/rawDial/blockedAddr) are unexported on purpose: the
// production constructor (newWebShared → newNetGuard) NEVER sets rawDial or
// blockedAddr, and only webGuard looks are config-fed. rawDial is a veth
// that runs BELOW the judge — it receives only the numeric address the
// judge already passed, so it can relocate a connect but can never make one
// happen for an address the judge refused. blockedAddr replacement exists
// for the conversion suites' shape claims; the SSRF claims themselves run
// with the production predicate (web_guard_test.go header). A composition
// root cannot switch the SSRF check off — there is no such switch in the
// contract (unwired.go: net.block_private_ranges "needs a contract decision
// (D22), not a consumer").
type netGuard struct {
	client       *http.Client
	lookup       func(ctx context.Context, host string) ([]netip.Addr, error)
	rawDial      func(ctx context.Context, network, numericAddr string) (net.Conn, error)
	blockedAddr  func(netip.Addr) bool // nil = ssrfBlockedAddr
	dialTimeout  time.Duration
	maxRedirects int
}

func newNetGuard(o guardOptions) *netGuard {
	dialTimeout := o.dialTimeout
	if dialTimeout <= 0 {
		dialTimeout = defaultDialTimeout
	}
	g := &netGuard{
		lookup:       o.lookup,
		rawDial:      o.rawDial,
		blockedAddr:  o.blockedAddr,
		dialTimeout:  dialTimeout,
		maxRedirects: orDefaultInt(o.maxRedirects, defaultMaxRedirects),
	}
	g.client = &http.Client{
		Timeout: o.fetchTimeout,
		// Proxy nil on purpose: an environment HTTP_PROXY would move the
		// SSRF decision into a process the guard cannot see. [net].proxy is
		// today consumed only by the LLM transport (config unwired.go ledger);
		// if web tools ever honor it, the guard runs against the PROXY as
		// next-hop address and that is a D30 contract conversation, not a
		// default.
		Transport: &http.Transport{
			DialContext:       g.dialContext,
			Proxy:             nil,
			ForceAttemptHTTP2: false, // no ALPN-over-a-pinned-IP surprises in v1
		},
		CheckRedirect: g.checkRedirect,
	}
	return g
}

type guardOptions struct {
	lookup       func(ctx context.Context, host string) ([]netip.Addr, error)
	rawDial      func(ctx context.Context, network, numericAddr string) (net.Conn, error)
	blockedAddr  func(netip.Addr) bool
	dialTimeout  time.Duration
	fetchTimeout time.Duration
	maxRedirects int
}

func (g *netGuard) blocked(ip netip.Addr) bool {
	if g.blockedAddr != nil {
		return g.blockedAddr(ip)
	}
	return ssrfBlockedAddr(ip)
}

// resolveAnswers turns a host into the candidate addresses, running the
// encoded-literal decode first (the DNS rebinding defense works on names;
// numeric spellings must NEVER reach a resolver).
func (g *netGuard) resolveAnswers(ctx context.Context, host string) ([]netip.Addr, error) {
	if ip, encoded := decodeNumericHost(host); encoded {
		return []netip.Addr{ip}, nil
	}
	if g.lookup != nil {
		return g.lookup(ctx, host)
	}
	ctxd, cancel := context.WithTimeout(ctx, g.dialTimeout)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIPAddr(ctxd, host)
	if err != nil {
		return nil, err
	}
	out := make([]netip.Addr, 0, len(ips))
	for _, a := range ips {
		ip, ok := netip.AddrFromSlice(a.IP)
		if !ok {
			continue
		}
		out = append(out, ip)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("解析 %s 没有可用地址", host)
	}
	return out, nil
}

// dialContext is the guard's only route to a socket: resolve, drop every
// blocked answer, then dial the numeric form of ONE checked address. The
// hostname never reaches net.Dialer, so a second DNS answer cannot change
// what connects (rebinding closes at the pin).
func (g *netGuard) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("SSRF 防御：无法解析拨号地址：%w", err)
	}
	if ip, encoded := decodeNumericHost(host); encoded && g.blocked(ip) {
		return nil, fmt.Errorf("SSRF 防御：目标在私有/回环/链路本地网段，连接前拒绝：%s", host)
	}
	answers, err := g.resolveAnswers(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("SSRF 防御：主机解析失败，未建立连接：%s: %w", host, err)
	}
	blockedAny := false
	var dialErr error
	for _, ip := range answers {
		if g.blocked(ip) {
			blockedAny = true
			continue
		}
		pinned := net.JoinHostPort(ip.String(), port)
		var c net.Conn
		var err error
		if g.rawDial != nil {
			// test veth BELOW the judge: the address it sees is one the
			// judge passed; refused answers never reach this line.
			c, err = g.rawDial(ctx, network, pinned)
		} else {
			d := &net.Dialer{Timeout: g.dialTimeout}
			c, err = d.DialContext(ctx, network, pinned)
		}
		if err == nil {
			return c, nil
		}
		dialErr = err
	}
	if blockedAny && dialErr == nil {
		return nil, fmt.Errorf("SSRF 防御：%s 解析结果全部落在私有/回环/链路本地网段，连接前拒绝", host)
	}
	if dialErr != nil {
		return nil, dialErr
	}
	return nil, fmt.Errorf("SSRF 防御：%s 无可用地址（未建立连接）", host)
}

// checkRedirect re-runs the guard on every hop: only http/https, never a
// https->http downgrade, every encoded/private literal rejected BEFORE the
// next dial (which itself re-resolves through the guard).
func (g *netGuard) checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= g.maxRedirects {
		return fmt.Errorf("重定向链超过 %d 跳，已中止", g.maxRedirects)
	}
	target := req.URL.String()
	if len(target) > risk.MaxURLLength {
		return fmt.Errorf("重定向目标超过 %d 字符上限", risk.MaxURLLength)
	}
	switch strings.ToLower(req.URL.Scheme) {
	case "http", "https":
	default:
		return fmt.Errorf("禁止跟随跳转到非 http(s) 协议：%s", req.URL.Scheme)
	}
	prev := via[0].URL
	if strings.EqualFold(prev.Scheme, "https") && strings.EqualFold(req.URL.Scheme, "http") {
		return errors.New("禁止 https 降级跳转到 http")
	}
	if host := req.URL.Hostname(); ssrfBlockedHost(host) {
		return fmt.Errorf("重定向目标在私有/回环网段，跳转前即拒绝：%s", host)
	}
	return nil
}

// fetchedDoc is one successful guarded fetch. Body is the capped raw bytes;
// nothing in this package logs or stores it beyond the tool's own markdown
// conversion result.
type fetchedDoc struct {
	FinalURL    string
	Status      int
	ContentType string
	Body        []byte
	Truncated   bool
}

// fetch performs GET through the guard with the size cap and the C22 context
// deadline (already on ctx from the bridge; the client Timeout is the second
// belt for host-internal callers). Error texts name status and URL, NEVER a
// byte of the body.
func (g *netGuard) fetch(ctx context.Context, rawURL string, maxBytes int) (fetchedDoc, error) {
	u, err := normalizeGuardedURL(rawURL)
	if err != nil {
		return fetchedDoc{}, err
	}
	if maxBytes <= 0 {
		maxBytes = defaultFetchMaxBytes
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fetchedDoc{}, fmt.Errorf("构造请求失败：%w", err)
	}
	req.Header.Set("User-Agent", webUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/json;q=0.9,text/plain;q=0.8,*/*;q=0.5")
	resp, err := g.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fetchedDoc{}, ctx.Err()
		}
		return fetchedDoc{}, fmt.Errorf("请求失败：%v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxBytes)+1))
	if err != nil {
		if ctx.Err() != nil {
			return fetchedDoc{}, ctx.Err()
		}
		return fetchedDoc{}, fmt.Errorf("读取响应失败：%v", err)
	}
	truncated := len(body) > maxBytes
	if truncated {
		body = body[:maxBytes]
	}
	final := ""
	if resp.Request != nil && resp.Request.URL != nil {
		final = resp.Request.URL.String()
	}
	return fetchedDoc{
		FinalURL:    final,
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        body,
		Truncated:   truncated,
	}, nil
}

func orDefaultInt(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}
