package main

// wisp secret - the credential entry CLI (ticket 63, ruling R7).
//
// It is a thin shell over ticket 06's DPAPI SecretStore: no new storage
// mechanism, no new file format, no crypto here. The command exists so the
// owner can put an API key into the store without it ever passing through a
// chat UI, a shell history line or a log file.
//
// The three leak surfaces this ticket is responsible for, and how they are
// closed by construction:
//
//   - argv: the secret can only arrive as (a) a hidden console read or (b) the
//     bytes on stdin when --from-stdin is passed. There is deliberately no
//     string-typed flag anywhere in this file (pinned by
//     TestSecretFlagsAreBoolOnly), so no argv slot can hold a value; and
//     Get-CimInstance Win32_Process / the PEB command line of a running
//     `wisp secret set --from-stdin` therefore shows only the blob name
//     (pinned by TestSecretArgvCarriesNoSecret, which also pins the shape of
//     that argv against a planted-value control,
//     TestProcessCommandLineProbeDetectsAPlantedValue).
//   - logs / error strings: everything printed or logged names the ref, the
//     env and the field path. The only place secret material is rendered at
//     all is secret.RedactSecret (last 4), and `--show` writing plaintext to
//     stdout is the single sanctioned exception (it is never logged).
//   - intermediate files: collectSecret reads stdin into one buffer and hands
//     it straight to the store; nothing is ever written before it is encrypted
//     (pinned by TestSecretFromStdinWritesNoIntermediateFile, which walks the
//     data dir and the OS temp dir for the entered value).
//
// This command never calls proc.Boot: entering a credential must work while
// the resident instance is running, and the store needs no mutex (the blob
// files are per-user, and last-writer-wins on one ref is the store's existing
// semantics from ticket 06).

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/secret"
	"golang.org/x/term"
)

// maxSecretBytes caps a --from-stdin read. API keys are tens of bytes; the cap
// exists so a stray `cat bigfile.bin | wisp secret set` cannot turn the
// credential entry point into a file uploader (and cannot land a large
// plaintext buffer in memory).
const maxSecretBytes = 8 << 10

// configFileName is the config file of the active env (SPEC-02 §6: it sits in
// the data root, next to secrets\). `wisp secret unset` reads it to find out
// whether a blob is still referenced.
const configFileName = "config.toml"

// timeLayout is the listing/notice timestamp form: UTC RFC3339 with seconds,
// so an audit line can be correlated with the log pipeline's own stamps.
const timeLayout = "2006-01-02T15:04:05Z"

// errNoTerminal marks "there is no interactive console to type a secret at".
var errNoTerminal = errors.New("no interactive console on stdin")

// secretIO is the process surface of one `wisp secret` run, injected so the
// leak tests can drive the command in-process with captured writers and a fake
// hidden reader.
type secretIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	// readHidden returns one secret typed with console echo disabled. It is
	// the only channel a secret value arrives on besides stdin, which is why
	// argv can be proved irrelevant to it.
	readHidden func(prompt string) (string, error)
	// now supplies the clock for the listing/notice timestamps.
	now func() time.Time
}

// secretCmd is one `wisp secret` invocation bound to an environment's data
// root. env/dataDir/portable come from ticket 06's layout fork, so the three
// envs are three separate stores by construction (SPEC-03 §5.2).
type secretCmd struct {
	env      buildinfo.Env
	dataDir  string
	portable bool
	sio      secretIO
}

// sessionLayout resolves the data root a `wisp secret` run writes into: the
// per-env fork plus the portable.txt override, with the same precedence the
// resident process uses (SPEC-03 §5.2, SPEC-02 §6). configRoot/exeDir are
// parameters rather than OS reads here so a test can point dev/test/prod at
// three temp dirs and prove the isolation through this exact command surface.
// The consequence of that is stated in R-119-3's terms: a root passed in here is
// a declaration, so this function must not resolve it for the caller - the
// wrapper that reads the OS does that (resolveSecretLayout), and the placement
// floor keeps refusing anything still spelled through a link.
func sessionLayout(env buildinfo.Env, configRoot, exeDir string) (proc.Layout, error) {
	l, err := proc.LayoutFor(env, configRoot)
	if err != nil {
		return proc.Layout{}, err
	}
	if exeDir != "" {
		l, _, err = proc.ApplyPortableOverride(env, l, exeDir)
		if err != nil {
			return proc.Layout{}, err
		}
	}
	return l, nil
}

// resolveSecretLayout is the production wrapper around sessionLayout: the real
// user config dir and the real exe directory.
func resolveSecretLayout(env buildinfo.Env) (proc.Layout, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return proc.Layout{}, fmt.Errorf("wisp secret: user config dir: %w", err)
	}
	// Ticket 119's rework (R-119-1): `wisp secret` is the second production
	// reader of that OS answer, and it never passes through proc.DefaultLayout,
	// so the resolution ticket 119 installed there did not reach this line. On a
	// dotfiles-managed Linux home $HOME/.config is a symlink, and macOS spells
	// both /tmp and /var that way, so the root handed to winsec's placement floor
	// (internal/winsec/winsec_other.go, ticket 113) was unresolved and the first
	// sealing call refused it: measured on one and the same binary, `wisp doctor`
	// printed the resolved tree while `wisp secret list` returned rc=1 in three
	// dev shapes (HOME through a link, $HOME/.config itself a link,
	// XDG_CONFIG_HOME through a link). Ruling R-119-3, pinned by
	// TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119: what this process
	// learns by asking the OS is resolved, because a data root is a location
	// contract; what it is handed by name (proc.TestDataDirEnv) is returned
	// verbatim, because whoever wrote it compares strings with the answer.
	root = proc.SealableRoot(root)
	exeDir := ""
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}
	// os.Executable failing is not fatal: without an exe dir there is no
	// portable.txt to find, so the per-env fork dir stands.
	return sessionLayout(env, root, exeDir)
}

const secretUsage = `wisp secret - credential entry (keys never pass through argv, logs or chat)`

const secretSubUsage = `
Usage:
  wisp secret set   <name> [--from-stdin]
      hidden console input, typed twice to confirm; --from-stdin reads the
      value from stdin instead. The plaintext never becomes an argv value.
  wisp secret get   <name> [--show]
      masked (last 4 only) by default; --show prints the plaintext to stdout
      once, and logs nothing.
  wisp secret list
      blob ids with their creation times and refs. Never any content.
  wisp secret unset <name> [--force]
      refuses while the active env's config.toml still references the blob;
      --force deletes it and writes an audit line naming the dangling fields.

<name> is the blob id: 1-128 chars of [A-Za-z0-9._-] (no path separators).
The reference to paste into config.toml is api_key_ref = "dpapi:<name>".
Environment: the active WISP_ENV is echoed by every command (SPEC-03 §5.2),
because storing a dev key into the prod store by mistake is the failure mode.
`

// cmdSecret is the `wisp secret` entry point (main dispatch): resolve the
// env, resolve its data root, and run the subcommand against it.
func cmdSecret(args []string) int {
	attachParentConsole()
	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp secret: %v\n", err)
		return 2
	}
	layout, err := resolveSecretLayout(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp secret: %v\n", err)
		return 2
	}
	c := &secretCmd{
		env:      env,
		dataDir:  layout.DataDir,
		portable: layout.Portable,
		sio: secretIO{
			stdin:      os.Stdin,
			stdout:     os.Stdout,
			stderr:     os.Stderr,
			readHidden: termHiddenReader(os.Stdin, os.Stderr),
			now:        time.Now,
		},
	}
	return c.run(args)
}

// termHiddenReader is the production readHidden: golang.org/x/term
// ReadPassword over the stdin handle, prompt on the caller's writer. Echo
// stays off for the whole read, so nothing the user types is reflected into a
// terminal scrollback line or a redirected stderr.
func termHiddenReader(stdin *os.File, promptTo io.Writer) func(string) (string, error) {
	return func(prompt string) (string, error) {
		fd := int(stdin.Fd())
		if !term.IsTerminal(fd) {
			return "", errNoTerminal
		}
		fmt.Fprint(promptTo, prompt)
		b, err := term.ReadPassword(fd)
		fmt.Fprint(promptTo, "\n") // the newline the disabled echo owes us
		if err != nil {
			return "", fmt.Errorf("reading hidden input: %w", err)
		}
		return string(b), nil
	}
}

// run dispatches one subcommand. Exit codes follow the `wisp slo` convention:
// 2 usage/environment, 1 operational refusal or failure, 0 success.
func (c *secretCmd) run(args []string) int {
	if len(args) == 0 {
		c.usage()
		return 2
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "set":
		return c.runSet(rest)
	case "get":
		return c.runGet(rest)
	case "list":
		return c.runList(rest)
	case "unset":
		return c.runUnset(rest)
	case "help", "-h", "--help":
		c.usage()
		return 0
	default:
		fmt.Fprintf(c.sio.stderr, "wisp secret: unknown subcommand %q\n", sub)
		c.usage()
		return 2
	}
}

func (c *secretCmd) usage() {
	fmt.Fprint(c.sio.stderr, secretUsage+"\n"+secretSubUsage+c.envFooter())
}

// envFooter is the one line every usage/diagnostic surface carries: which env
// this command is about to act on, and where (SPEC-03 §5.2; the same rule the
// ball badge follows - a bare "Wisp" means prod).
func (c *secretCmd) envFooter() string {
	return fmt.Sprintf("Active environment: WISP_ENV=%s, data dir=%s, portable=%v\n",
		c.env, c.dataDir, c.portable)
}

// openStore opens ticket 06's store for this invocation's data dir, passing
// the layout's portable flag through unchanged: the CLI has no portable-mode
// logic of its own (P13 lives in the store).
func (c *secretCmd) openStore() (*secret.Store, error) {
	return secret.NewStore(c.dataDir, secret.WithPortable(c.portable))
}

func (c *secretCmd) configPath() string {
	return filepath.Join(c.dataDir, configFileName)
}

// refFor maps a user-supplied <name> onto its ref. ValidBlobID is ticket 06's
// traversal guard and is the only validation the CLI adds: the name becomes a
// file name under secrets\, so it can never contain a separator.
func refFor(name string) (string, error) {
	if !secret.ValidBlobID(name) {
		return "", fmt.Errorf("secret name %q is not usable as a blob id: want 1-128 characters of [A-Za-z0-9._-], with no path separators or spaces", name)
	}
	return secret.RefPrefixDPAPI + name, nil
}

// newFlagSet builds a subcommand FlagSet whose diagnostics go to the injected
// stderr. Only boolean flags may ever be added: see TestSecretFlagsAreBoolOnly.
func (c *secretCmd) newFlagSet(name, usage string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(c.sio.stderr)
	fs.Usage = func() { fmt.Fprint(c.sio.stderr, usage+"\n"+c.envFooter()) }
	return fs
}

// runSet stores one credential. The value comes from stdin (--from-stdin) or
// from two hidden console reads; never from an argument.
func (c *secretCmd) runSet(args []string) int {
	fs := c.newFlagSet("wisp secret set", `Usage: wisp secret set <name> [--from-stdin]
  hidden console input, typed twice to confirm; --from-stdin reads the value
  from stdin in a single read and writes no intermediate file.`)
	fromStdin := fs.Bool("from-stdin", false, "read the secret from stdin instead of the console")
	pos, code, done := c.parseFlags(fs, args, "set")
	if done {
		return code
	}
	if len(pos) != 1 {
		fmt.Fprintf(c.sio.stderr, "wisp secret set: want exactly one <name>, got %d\n", len(pos))
		fs.Usage()
		return 2
	}
	name := pos[0]
	ref, err := refFor(name)
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret set: %v\n", err)
		return 2
	}
	value, err := c.collectSecret(*fromStdin)
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret set: %v\n", err)
		return 1
	}
	st, err := c.openStore()
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret set: %v\n", err)
		return 1
	}
	existed, err := st.Exists(ref)
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret set: %v\n", err)
		return 1
	}
	if err := st.Store(ref, value); err != nil {
		// The CLI adds the ref the operation was about; Store's own errors
		// carry the ref or the blob path and never the value.
		fmt.Fprintf(c.sio.stderr, "wisp secret set: storing %s: %v\n", ref, err)
		return 1
	}
	fmt.Fprintf(c.sio.stdout, "wisp secret set: stored %s (WISP_ENV=%s, portable=%v)\n", ref, c.env, c.portable)
	if existed {
		fmt.Fprintf(c.sio.stdout, "  note: %s already existed and was overwritten\n", ref)
	}
	fmt.Fprintf(c.sio.stdout, "  masked: %s\n", secret.RedactSecret(value))
	fmt.Fprintf(c.sio.stdout, "  blob:   %s\n", filepath.Join(st.Dir(), name))
	fmt.Fprintf(c.sio.stdout, "  reference for config.toml: api_key_ref = %q\n", ref)
	fmt.Fprint(c.sio.stdout, c.envFooter())
	// Audit without secret material: the masked form is the most that may ever
	// be logged (C28), and even that is left to RedactSecret's caller in the
	// GUI path - here the ref is the fact worth recording.
	slog.Info("wisp secret: stored dpapi blob", "ref", ref, "env", string(c.env), "portable", c.portable, "overwrote", existed)
	return 0
}

// runGet resolves one credential. Masked unless --show is passed.
func (c *secretCmd) runGet(args []string) int {
	fs := c.newFlagSet("wisp secret get", `Usage: wisp secret get <name> [--show]
  prints the masked form (last 4 characters) by default; --show prints the
  plaintext to stdout exactly once and writes no log record.`)
	show := fs.Bool("show", false, "print the plaintext secret to stdout")
	pos, code, done := c.parseFlags(fs, args, "get")
	if done {
		return code
	}
	if len(pos) != 1 {
		fmt.Fprintf(c.sio.stderr, "wisp secret get: want exactly one <name>, got %d\n", len(pos))
		fs.Usage()
		return 2
	}
	ref, err := refFor(pos[0])
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret get: %v\n", err)
		return 2
	}
	st, err := c.openStore()
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret get: %v\n", err)
		return 1
	}
	value, err := st.Resolve(ref)
	if err != nil {
		// Includes P13: under portable mode an undecryptable blob surfaces as
		// ErrPortableDecrypt with env: guidance from the store. The CLI adds
		// nothing to that text (and it carries no secret material).
		fmt.Fprintf(c.sio.stderr, "wisp secret get: %v\n", err)
		return 1
	}
	if !*show {
		fmt.Fprintf(c.sio.stdout, "wisp secret get: %s (WISP_ENV=%s, portable=%v)\n", ref, c.env, c.portable)
		fmt.Fprintf(c.sio.stdout, "  masked: %s\n", secret.RedactSecret(value))
		fmt.Fprint(c.sio.stdout, "  pass --show for the plaintext\n")
		return 0
	}
	// The one sanctioned plaintext output of this ticket: stdout only, no
	// trailing newline (so `wisp secret get k --show | clip` copies the key
	// byte-exactly), and nothing logged.
	fmt.Fprint(c.sio.stderr, "wisp secret get: --show writes the plaintext to stdout; do not redirect it to a file or a log\n")
	fmt.Fprint(c.sio.stdout, value)
	return 0
}

// runList enumerates blobs by metadata only.
func (c *secretCmd) runList(args []string) int {
	fs := c.newFlagSet("wisp secret list", `Usage: wisp secret list
  lists blob ids, creation times (UTC) and refs. Never any secret content.`)
	pos, code, done := c.parseFlags(fs, args, "list")
	if done {
		return code
	}
	if len(pos) != 0 {
		fmt.Fprintf(c.sio.stderr, "wisp secret list: takes no name, got %d argument(s)\n", len(pos))
		fs.Usage()
		return 2
	}
	st, err := c.openStore()
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret list: %v\n", err)
		return 1
	}
	blobs, err := st.Blobs()
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret list: %v\n", err)
		return 1
	}
	fmt.Fprintf(c.sio.stdout, "wisp secret list: WISP_ENV=%s, portable=%v, dir=%s\n", c.env, c.portable, st.Dir())
	if len(blobs) == 0 {
		fmt.Fprint(c.sio.stdout, "  (no blobs)\n")
		return 0
	}
	for _, b := range blobs {
		fmt.Fprintf(c.sio.stdout, "  %-40s %s  %s\n", b.ID, b.Created.Format(timeLayout), b.Ref)
	}
	fmt.Fprintf(c.sio.stdout, "  %d blob(s)\n", len(blobs))
	return 0
}

// runUnset deletes a blob, refusing while the active config still references
// it unless --force is passed.
func (c *secretCmd) runUnset(args []string) int {
	fs := c.newFlagSet("wisp secret unset", `Usage: wisp secret unset <name> [--force]
  refuses when config.toml of the active environment still names this blob,
  listing the referencing fields; --force deletes anyway and writes an audit
  line.`)
	force := fs.Bool("force", false, "delete even while config.toml references the blob (writes an audit line)")
	pos, code, done := c.parseFlags(fs, args, "unset")
	if done {
		return code
	}
	if len(pos) != 1 {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: want exactly one <name>, got %d\n", len(pos))
		fs.Usage()
		return 2
	}
	name := pos[0]
	ref, err := refFor(name)
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: %v\n", err)
		return 2
	}
	configPath := c.configPath()
	// Fail closed on an unreadable config: "is anything referencing this?"
	// must be answered, not assumed.
	fields, err := secret.RefFieldNames(configPath, ref)
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: %v\n", err)
		return 1
	}
	if len(fields) > 0 && !*force {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: %s is still referenced by %s at:\n", ref, configPath)
		for _, f := range fields {
			fmt.Fprintf(c.sio.stderr, "  %s\n", f)
		}
		fmt.Fprintf(c.sio.stderr, "repoint or remove those fields first, or pass --force to delete anyway\n")
		return 1
	}
	st, err := c.openStore()
	if err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: %v\n", err)
		return 1
	}
	if err := st.Delete(ref); err != nil {
		fmt.Fprintf(c.sio.stderr, "wisp secret unset: %v\n", err)
		return 1
	}
	audit := fmt.Sprintf("wisp secret unset: deleted %s (WISP_ENV=%s, portable=%v, forced=%v, referenced_fields=%d)",
		ref, c.env, c.portable, *force, len(fields))
	fmt.Fprintf(c.sio.stdout, "audit: %s\n", audit)
	if len(fields) > 0 {
		// --force on a referenced blob is the one destructive act in this
		// ticket, so it is logged at Warn with the dangling field paths:
		// names and refs only, never values (C28).
		fmt.Fprintf(c.sio.stdout, "  dangling references left in %s:\n", configPath)
		for _, f := range fields {
			fmt.Fprintf(c.sio.stdout, "    %s\n", f)
		}
		slog.Warn(audit + " fields=" + strings.Join(fields, ","))
	} else {
		slog.Info(audit)
	}
	return 0
}

// parseFlags runs a subcommand's flag parse and returns the positional
// arguments. Flags are accepted in any position: `wisp secret set mykey
// --from-stdin` is the form the usage text and the ticket document, and the
// stdlib parser would otherwise stop at the first positional and silently
// treat --from-stdin as a name. A bare "--" ends flag parsing, so a name that
// starts with '-' stays expressible.
//
// The second return is the exit code and the third says the caller must stop
// (usage error, or -h, which prints and succeeds).
func (c *secretCmd) parseFlags(fs *flag.FlagSet, args []string, sub string) ([]string, int, bool) {
	var flagArgs, pos []string
	positionsOnly := false
	for _, a := range args {
		switch {
		case positionsOnly:
			pos = append(pos, a)
		case a == "--":
			positionsOnly = true
		case strings.HasPrefix(a, "-") && a != "-":
			flagArgs = append(flagArgs, a)
		default:
			pos = append(pos, a)
		}
	}
	if err := fs.Parse(flagArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, 0, true
		}
		// flag prints "flag provided but not defined: -<name>" (the name only,
		// never the value) plus the usage above.
		fmt.Fprintf(c.sio.stderr, "wisp secret %s: %v\n", sub, err)
		return nil, 2, true
	}
	return append(fs.Args(), pos...), 0, false
}

// collectSecret obtains the plaintext from the one channel that is not argv.
func (c *secretCmd) collectSecret(fromStdin bool) (string, error) {
	if fromStdin {
		raw, err := io.ReadAll(io.LimitReader(c.sio.stdin, int64(maxSecretBytes)+1))
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		if len(raw) > maxSecretBytes {
			return "", fmt.Errorf("stdin carries more than %d bytes: a credential is a single short line, not a file", maxSecretBytes)
		}
		return singleLine(raw)
	}
	if c.sio.readHidden == nil {
		return "", c.wrapHidden(errNoTerminal)
	}
	first, err := c.sio.readHidden("secret (input hidden): ")
	if err != nil {
		return "", c.wrapHidden(err)
	}
	second, err := c.sio.readHidden("repeat it to confirm (input hidden): ")
	if err != nil {
		return "", c.wrapHidden(err)
	}
	if first == "" {
		return "", errors.New("no secret entered")
	}
	if first != second {
		return "", errors.New("the two entries did not match; nothing was stored")
	}
	return first, nil
}

// singleLine normalizes the stdin read: trailing newline(s) only, then a hard
// refusal of any remaining line break (a pasted multi-line blob is a mistake,
// not a key, and silently storing the first line would corrupt it).
func singleLine(raw []byte) (string, error) {
	v := strings.TrimRight(string(raw), "\r\n")
	if v == "" {
		return "", errors.New("stdin carried no secret")
	}
	if strings.ContainsAny(v, "\r\n") {
		return "", errors.New("stdin carries more than one line; a credential must be a single line")
	}
	if !utf8.ValidString(v) {
		return "", errors.New("stdin is not valid UTF-8; refusing to store it")
	}
	return v, nil
}

func (c *secretCmd) wrapHidden(err error) error {
	if errors.Is(err, errNoTerminal) {
		return fmt.Errorf("%w: run this at a terminal, or pipe the value in with --from-stdin", errNoTerminal)
	}
	return err
}
