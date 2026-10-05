package main

// Ticket 257, shape (c), the cmd/wisp half: the first-run receipt is the ONLY
// place this build tells an operator how a provider row comes to exist. The
// internal/config half (257-r1) made the refusals say which of the three
// reasons they are; this file's leg is the sentence printed once, on a machine
// that has never had a config.toml, and these cases are the teeth on it.
//
// What is deliberately NOT reused from the sibling package: its tag constants.
// refusalFileMissing / refusalRowMissing / refusalInvalid are unexported, and
// even if they were exported, quoting a same-package constant is the blind
// spot 257-r1b measured (its M1: rename the constant and the assertion renames
// itself, so nothing reddens). Every expectation below is a literal typed into
// this file, which is what makes the two halves' wording a claim that can fail
// rather than a shared identifier that cannot drift.
//
// AC#3 boundary honored here: nothing in this file stores, resolves, echoes or
// invents a credential VALUE. The roster walk addresses the config-side leg of
// the credential field (the reference), exactly as 257-r1's walk does, and the
// reference used is a placeholder name (an env var that is never set, a canary
// blob id that never touches a store).

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
)

// The provider and model the receipt tells the operator to add by hand. Both
// names are the ones the guidance's placeholders stand for; "deepseek" is an
// internal preset so the taught shape ("对上内置预设的，protocol 与 base_url
// 可以留空") is exercised rather than assumed.
const (
	t257Provider  = "deepseek"
	t257Model     = "deepseek-chat"
	t257ModelBis  = "deepseek-reasoner"
	t257EnvRef    = "env:T257_R2_ENV_NAME_NEVER_SET"
	t257CanaryRef = "dpapi:t257r2canaryrefplaceholder"
)

// The three refusal reasons, copied as literals from internal/config's tags and
// from firstrun.go's receipt. 票 257 AC#2: three different things an operator
// fixes three different ways; 票 223 AC#4's discipline says one shared line is
// a lie by omission.
const (
	t257ReasonFileMissing = "第 1 种拒因：文件没建"
	t257ReasonRowMissing  = "第 2 种拒因：行不存在"
	t257ReasonInvalid     = "第 3 种拒因：校验不过"
)

// The strings AC#1 demands the receipt itself carry: the path spelled the way
// the schema reads it, the model row with its enabled key, and the in-place
// fill of the role section that already exists.
const (
	t257TeachProviderRow  = "第一样＝一节 [llm.providers.<名>]"
	t257TeachModelRow     = "[llm.providers.<名>.models.<模型 id>]，里面写 enabled = true"
	t257TeachRolePair     = "第三样＝就地填已有的 [llm.roles.chat] 那一节"
	t257TeachSevenOfSeven = "三样齐了这七枚才全部写得进"
	t257TeachThreeOfSeven = "只加第一样只解锁服务商那三枚"
)

// t257CleanMachine runs the user-started entry ONCE over an empty temporary
// data root and returns what the operator would see. This is the ticket's
// "干净机器那一发": the premise is a missing config.toml, asserted before the
// run, and the file that exists afterwards was written by production
// (runTextTask -> ensureFirstRunConfig), never hand-placed by this test.
func t257CleanMachine(t *testing.T) (dir, cfgPath, receipt string) {
	t.Helper()
	dir = t.TempDir()
	cfgPath = filepath.Join(dir, configFileName)
	if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("AC#1 premise broke: the clean machine already holds %q (stat err %v)", cfgPath, err)
	}
	code, _, stderrText := run198(t, dir)
	// Exit 2 is ticket 198's pin, restated because it is also this ticket's
	// premise: creating the file never softens "no model, no key".
	if code != 2 {
		t.Fatalf("first run over a fresh root: exit %d, want 2 (a created config is still an unconfigured one), stderr:\n%s", code, stderrText)
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("first run created no config.toml, so there is no clean machine to walk: %v", err)
	}
	if !strings.Contains(stderrText, "新建默认配置") {
		t.Fatalf("the run created the file without saying so, so there is no receipt to read:\n%s", stderrText)
	}
	return dir, cfgPath, stderrText
}

// t257Field is one row of the panel's writable roster plus the leg
// cmd/wisp/panel_config_store.go's switch routes it through. Seven rows, six
// distinct legs: provider_credential's value leg is StoreCredential, which
// this file must not exercise (AC#3), so the row is addressed by the same
// config-side reference leg the panel itself falls back to - the shape
// 257-r1's walk uses, restated here rather than imported.
type t257Field struct {
	field string
	leg   string
	run   func(m *config.Manager) ([]string, error)
}

func t257RosterWalk() []t257Field {
	return []t257Field{
		{"provider_base_url", "Manager.SetProviderBaseURL", func(m *config.Manager) ([]string, error) {
			return m.SetProviderBaseURL(t257Provider, "https://clean-machine.example.invalid/v1")
		}},
		{"provider_api_key_ref", "Manager.SetProviderAPIKeyRef", func(m *config.Manager) ([]string, error) {
			return m.SetProviderAPIKeyRef(t257Provider, t257EnvRef)
		}},
		{"model_context_window", "Manager.SetModelContextWindow", func(m *config.Manager) ([]string, error) {
			return m.SetModelContextWindow(t257Provider, t257Model, 128000)
		}},
		{"model_price_in", "Manager.SetModelPriceIn", func(m *config.Manager) ([]string, error) {
			return m.SetModelPriceIn(t257Provider, t257Model, 3500000)
		}},
		{"model_price_out", "Manager.SetModelPriceOut", func(m *config.Manager) ([]string, error) {
			return m.SetModelPriceOut(t257Provider, t257Model, 12000000)
		}},
		{"role_chat_model", "Manager.SetRoleChatModel", func(m *config.Manager) ([]string, error) {
			return m.SetRoleChatModel(t257ModelBis)
		}},
		{
			"provider_credential", "Manager.SetProviderAPIKeyRef (StoreCredential's config-side leg)",
			func(m *config.Manager) ([]string, error) {
				return m.SetProviderAPIKeyRef(t257Provider, t257CanaryRef)
			},
		},
	}
}

// TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain is AC#1's first
// half: on the machine that never had a config.toml, the receipt the run prints
// names the row path the schema actually reads, and every one of the seven
// refused writes says the same thing back.
func TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain(t *testing.T) {
	_, cfgPath, receipt := t257CleanMachine(t)

	for _, mustSay := range []string{
		"[llm.providers.<名>",
		t257TeachProviderRow,
		t257TeachModelRow,
		t257TeachRolePair,
		t257TeachSevenOfSeven,
		t257TeachThreeOfSeven,
	} {
		if !strings.Contains(receipt, mustSay) {
			t.Errorf("AC#1 RED: the first-run receipt never says %q, so shape (c) is not delivered - the operator is left with a refusal and no path:\n%s",
				mustSay, receipt)
		}
	}
	// A543's anti-control, as a literal: the ticket's ORIGINAL spelling is the
	// one that kills the panel chain, because a top-level [providers.x] is an
	// unknown key and decodeStrict refuses the whole file
	// (internal/config/parse.go's DisallowUnknownFields).
	if strings.Contains(receipt, "[providers.") {
		t.Errorf("AC#1 RED: the receipt tells the operator to add [providers.x], which parse.go's DisallowUnknownFields rejects outright (ledger A543):\n%s", receipt)
	}
	// The guidance lives in the receipt and nowhere else: the created file must
	// still carry no provider row (ticket 198's pin, same direction, re-measured
	// from this leg because it is the boundary between "tells" and "builds").
	created, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(created), "[llm.providers") {
		t.Errorf("AC#1 RED: the guidance leaked into the generated config.toml, which is invented-default territory and ticket 198's forbidden shape:\n%s", created)
	}

	m, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("NewManager over the created first config (production shape, res=nil per A560): %v", err)
	}
	if m.Config().LLM.Providers != nil {
		t.Fatalf("the ticket's premise is gone: the fresh registry is %#v, not nil", m.Config().LLM.Providers)
	}

	byReason := map[string]int{}
	for _, f := range t257RosterWalk() {
		written, err := f.run(m)
		if err == nil {
			t.Errorf("%s: the clean machine accepted the write (%v) - nothing was left to explain", f.field, written)
			continue
		}
		if written != nil {
			t.Errorf("%s: a refused write still reported key paths %v", f.field, written)
		}
		msg := err.Error()
		seen := 0
		for _, reason := range []string{t257ReasonFileMissing, t257ReasonRowMissing, t257ReasonInvalid} {
			if strings.Contains(msg, reason) {
				seen++
				byReason[reason]++
			}
		}
		if seen != 1 {
			t.Errorf("%s: refusal names %d of the three reasons, want exactly 1 (AC#2 forbids both folding and ambiguity): %s", f.field, seen, msg)
		} else if !strings.Contains(msg, t257ReasonRowMissing) {
			t.Errorf("%s: the file exists here, so the answer must be %q, got: %s", f.field, t257ReasonRowMissing, msg)
		}
		if !strings.Contains(msg, "[llm.providers.") {
			t.Errorf("%s: refusal does not name the section to hand-add, so it is a dead end and not shape (c): %s", f.field, msg)
		}
		for _, folded := range []string{"配置未生效：", "写入失败，配置未生效"} {
			if strings.Contains(msg, folded) {
				t.Errorf("%s: refusal folded into the shared line %q: %s", f.field, folded, msg)
			}
		}
	}
	if len(t257RosterWalk()) != 7 {
		t.Fatalf("the roster walk is %d fields, want 7 - AC#1 counts the panel's seven", len(t257RosterWalk()))
	}
	t.Logf("AC#1 clean machine: 7 refused writes by reason: %v", byReason)
}

// TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields is AC#1's second
// half and the honest one: the shape the receipt teaches is applied to the
// file the production run created, by string surgery that mirrors the printed
// sentence, and all seven fields then accept. If the guidance teaches a shape
// that does not load or does not unlock, this case says so.
func TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields(t *testing.T) {
	_, cfgPath, receipt := t257CleanMachine(t)
	body, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	// The receipt's third item promises the role section is ALREADY there and
	// gets filled in place. Re-assert it against the real file: if first-run
	// ever stops emitting it, the taught move is wrong and this case must say
	// so instead of quietly appending a section the way 257-r1b's M3 did.
	if !strings.Contains(receipt, "就地填已有的 [llm.roles.chat]") {
		t.Fatalf("the receipt stopped teaching the in-place fill: %s", receipt)
	}
	const emptyRolePair = "[llm.roles.chat]\nprovider = ''\nmodel = ''"
	if !strings.Contains(text, emptyRolePair) {
		t.Fatalf("premise drifted: the created file has no empty role pair to fill in place:\n%s", text)
	}

	hand := t257HandAdd(text, true, true, true)
	if err := os.WriteFile(cfgPath, []byte(hand), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("AC#1 RED: the hand-add shape the receipt teaches does not even load: %v", err)
	}
	if got := len(m.Config().LLM.Providers[t257Provider].Models); got != 2 {
		t.Fatalf("catalog rows after the taught hand add = %d, want 2", got)
	}

	accepted := 0
	for _, f := range t257RosterWalk() {
		written, err := f.run(m)
		if err != nil {
			t.Errorf("%s: still refused after the operator followed the receipt: %v", f.field, err)
			continue
		}
		if len(written) == 0 {
			t.Errorf("%s: accepted but named no key path, so the receipt can say nothing: %v", f.field, written)
			continue
		}
		accepted++
	}
	if accepted != 7 {
		t.Errorf("AC#1 RED: the taught chain unlocked %d of 7 fields, want 7 - the receipt claims %q", accepted, t257TeachSevenOfSeven)
	}

	after := t257ReadFile(t, cfgPath)
	if strings.Contains(after, "\napi_key =") {
		t.Errorf("AC#3 RED: a settings write put a plaintext key into config.toml:\n%s", after)
	}
}

// TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven pins the receipt's
// other number, the one that keeps "三样缺一不可" honest: add ONLY the provider
// row and exactly the three provider-level fields open up, no more.
func TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven(t *testing.T) {
	_, cfgPath, _ := t257CleanMachine(t)
	hand := t257HandAdd(t257ReadFile(t, cfgPath), true, false, false)
	if err := os.WriteFile(cfgPath, []byte(hand), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("a provider row alone must still load (preset name, protocol filled by applyPresets): %v", err)
	}
	accepted := map[string]bool{}
	for _, f := range t257RosterWalk() {
		if _, err := f.run(m); err == nil {
			accepted[f.field] = true
		}
	}
	if len(accepted) != 3 {
		t.Errorf("AC#1 RED: a provider row alone unlocked %v (%d of 7), want 3 - the receipt promises %q",
			accepted, len(accepted), t257TeachThreeOfSeven)
	}
	for field := range accepted {
		if field != "provider_base_url" && field != "provider_api_key_ref" && field != "provider_credential" {
			t.Errorf("AC#1 RED: %s unlocked with no model row and no role pair, which contradicts the refusal gates", field)
		}
	}
}

// TestTicket257R2AC2ThreeRefusalsStayThreeSentences is AC#2 on this leg: the
// receipt pre-states the three reasons a settings write says no, each with its
// own remedy, and never as one shared verdict.
func TestTicket257R2AC2ThreeRefusalsStayThreeSentences(t *testing.T) {
	_, cfgPath, receipt := t257CleanMachine(t)

	lines := strings.Split(receipt, "\n")
	remedy := map[string][]string{
		t257ReasonFileMissing: {"替你办完", "没有任何旧配置可言"},
		t257ReasonRowMissing:  {"改不了服务商与模型的存在性", "手加上面那三样"},
		t257ReasonInvalid:     {"过不了这份 schema 的校验", "文件一个字节都没动"},
	}
	for _, reason := range []string{t257ReasonFileMissing, t257ReasonRowMissing, t257ReasonInvalid} {
		var hits []string
		for _, ln := range lines {
			if strings.Contains(ln, reason) {
				hits = append(hits, ln)
			}
		}
		if len(hits) != 1 {
			t.Errorf("AC#2 RED: %q appears on %d receipt lines, want exactly 1 (three reasons, three sentences):\n%s",
				reason, len(hits), receipt)
			continue
		}
		for _, mustSay := range remedy[reason] {
			if !strings.Contains(hits[0], mustSay) {
				t.Errorf("AC#2 RED: the sentence for %q does not carry its own remedy %q, so two of the three causes are saying the same thing:\n%s",
					reason, mustSay, hits[0])
			}
		}
	}

	// Folding check: the shared verdict may only exist in the sentence that
	// forbids it, and exactly once.
	for _, folded := range []string{"配置未生效", "重启就好"} {
		n := strings.Count(receipt, folded)
		if n != 1 {
			t.Errorf("AC#2 RED: the folded verdict %q appears %d times, want exactly 1 (its own prohibition): %s", folded, n, receipt)
		}
		for _, ln := range strings.Split(receipt, "\n") {
			if strings.Contains(ln, folded) && !strings.Contains(ln, "不会合成一句") && !strings.Contains(ln, "不是一句") {
				t.Errorf("AC#2 RED: %q is used as a verdict rather than named as the forbidden collapse:\n%s", folded, ln)
			}
		}
	}
	// 票 223 AC#4's siblings: the first-run receipt must not borrow the
	// hot-reload missing-file sentence, which promises an old config this boot
	// never had.
	for _, borrowed := range []string{"cause=missing", "本次运行继续用内存里的旧配置"} {
		if strings.Contains(receipt, borrowed) {
			t.Errorf("AC#2 RED: the receipt borrows hot-reload wording %q:\n%s", borrowed, receipt)
		}
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("the receipt's first sentence says the file is there; the disk disagrees: %v", err)
	}
}

// TestTicket257R2AC2EffectTimingSaysThreeProcessShapes is 票 257 §8 boundary 1
// as amended by ledger A560: "手改＝热加载认／面板写＝要重启" is two shapes
// speaking for three, and the third (a resident boot whose task leg never
// passes the console gate) has nothing re-reading the disk at all.
func TestTicket257R2AC2EffectTimingSaysThreeProcessShapes(t *testing.T) {
	_, _, receipt := t257CleanMachine(t)

	shapes := []struct {
		name   string
		mustIn string
	}{
		{"控制台里的 wisp run", "每 1s 重读一次 config.toml 的看门狗"},
		{"没有可答卡入口的常驻形状", "根本没有会重读盘的东西"},
		{"设置页那一页", "自己明说不带轮询"},
	}
	for _, s := range shapes {
		if !strings.Contains(receipt, s.mustIn) {
			t.Errorf("AC#2 / boundary 1 RED: the receipt never states the %s shape (%q), so one boot shape is folded into another:\n%s",
				s.name, s.mustIn, receipt)
		}
	}
	// The two halves of the timing claim must not be merged into each other:
	// the hot line is only allowed to speak about THIS process's memory, and the
	// restart line only about the assembled model path.
	if !strings.Contains(receipt, "热加载不会替它换脑") {
		t.Errorf("AC#2 RED: the receipt lost the sentence separating memory from the once-built model path:\n%s", receipt)
	}
	// A560: the config-layer dereference channel is NOT wired in production, so
	// the receipt must not offer it as a way to pick up a changed key.
	if !strings.Contains(receipt, "把那份引用再解一次是新建端点时才做的事") {
		t.Errorf("AC#2 / A560 RED: the receipt no longer states when a changed api_key_ref is re-resolved, which is the exact place a confident sentence goes wrong:\n%s", receipt)
	}
}

// TestTicket257R2AC3CredentialSurfaceUntouched is AC#3 read off the two
// artifacts this leg can actually observe: the printed receipt (names and
// placeholders only, never a value) and the one function this leg may change.
func TestTicket257R2AC3CredentialSurfaceUntouched(t *testing.T) {
	_, cfgPath, receipt := t257CleanMachine(t)

	for _, form := range []struct{ prefix, placeholder string }{
		{"dpapi:", "dpapi:<blob 名>"},
		{"env:", "env:<环境变量名>"},
	} {
		n := strings.Count(receipt, form.prefix)
		if n == 0 {
			t.Errorf("AC#3: the receipt stopped naming %s references at all, which is the one way a key changes; re-check the count below against a text that moved elsewhere", form.prefix)
		}
		if got := strings.Count(receipt, form.placeholder); got != n {
			t.Errorf("AC#3 RED: the receipt uses %d %s references but only %d of them the placeholder %q; a quoted reference that is neither is a value written into operator text",
				n, form.prefix, got, form.placeholder)
		}
	}
	if strings.Contains(receipt, "api_key =") {
		t.Errorf("AC#3 RED: the receipt shows a plaintext key assignment shape:\n%s", receipt)
	}
	if strings.Contains(receipt, "sk-") {
		t.Errorf("AC#3 RED: the receipt carries a vendor key-shaped literal:\n%s", receipt)
	}
	created := t257ReadFile(t, cfgPath)
	if strings.Contains(created, "\napi_key =") || strings.Contains(created, "api_key =\"") {
		t.Errorf("AC#3 RED: the first config carries a plaintext key field the schema forbids:\n%s", created)
	}

	source, err := os.ReadFile(filepath.Join(".", "firstrun.go"))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(source), "\nfunc "); got != 1 {
		t.Errorf("AC#3 RED: firstrun.go now declares %d functions, want 1 (ensureFirstRunConfig). A new helper here is a new surface; the credential surface in particular must not grow one", got)
	}
	if !strings.Contains(string(source), "func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error) {") {
		t.Errorf("AC#3 RED: ensureFirstRunConfig's signature drifted; it takes a writer and answers (created, error) and nothing else - no read side, no value return")
	}
}

// --- helpers ---------------------------------------------------------------

// t257HandAdd performs the edits the receipt teaches, as literally as string
// surgery gets: fill the existing role pair in place, and append the provider
// row / its model rows as separate items so the three-of-seven and
// seven-of-seven readings differ only by which items are switched on.
func t257HandAdd(text string, providerRow, modelRows, rolePair bool) string {
	if rolePair {
		text = strings.Replace(text,
			"[llm.roles.chat]\nprovider = ''\nmodel = ''",
			"[llm.roles.chat]\nprovider = \""+t257Provider+"\"\nmodel = \""+t257Model+"\"",
			1)
	}
	if providerRow {
		text += "\n[llm.providers." + t257Provider + "]\napi_key_ref = \"env:T257_R2_HANDADD_ENV_NAME\"\n"
	}
	if modelRows {
		text += "\n[llm.providers." + t257Provider + ".models." + t257Model + "]\nenabled = true\ncontext_window = 64000\n"
		text += "\n[llm.providers." + t257Provider + ".models." + t257ModelBis + "]\nenabled = true\n"
	}
	return text
}

func t257ReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}
