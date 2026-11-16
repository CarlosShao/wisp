package config

// Ticket 215a's write half: plugins.<id>.enabled, the Tier-2 plugin switch,
// landing on the settings route ticket 248 opened and ticket 195's protected
// section-write entry beneath it.
//
// THE WRITE PATH IS NOT NEW. Manager.SetPluginEnabled is the 24th setter in
// the same family as SetBallStyle (settings_65.go) and every other settings
// field: one owned dotted key path, the write handed to Manager.WriteSection
// (sectionwrite.go), which validates before a byte moves, diffs the write to
// exactly one key, and takes the D36 direction verdict FIRST. No second write
// route is minted here, and no caller can reach this key except through the
// same verdict the reload path uses (lockedDirectionOf dispatches
// [plugins] to the SAME pluginsDirection the reload's planLocked runs, so a
// key planLocked calls a loosening can never be a tightening here).
//
// THE DIRECTION IS THE SECURITY ARCHITECTURE, NOT A CONVENIENCE.
// pluginsDirection (manager.go) classifies enabled false->true as a LOOSENING
// and true->false as a tighten. WriteSection then does what it always does:
//
//   - tighten (switch OFF)   - written immediately, memory and file, with the
//     tighten audit line ("locked section tightened" applies quietly);
//   - loosen from the panel  - refused outright (ErrPanelLoosenRefused): the
//     panel is a permission input (R20), and "由面板侧来源的 L2『允许』" is
//     on AGENTS.md sec 1.2's forbidden list - a page can never enable a
//     plugin, only ask the native side for one;
//   - loosen from native     - queued to the L2 approval queue (LoosenQueue;
//     nil refuses fail-closed with ErrLoosenQueueUnwired, the same
//     nil-means-deny rule Manager.ConfirmLocked carries).
//
// NO ROW IS INVENTED. The house rule for settings writes is "设置写不能新建
// 行": an id the live config does not declare is refused by name, and the
// Check below re-checks the FILE base (the write is judged against what the
// file holds, ticket 226's rule), so a typo or a hand-edited file can never
// plant a [plugins.<id>] row that was never there.
//
// THE TOGGLE ITSELF IS AUDITED IN EVERY OUTCOME (ticket 215 AC#3's "开关动
// 作本身要进审计（三档都写）"). EVERY attempt - including the ones refused
// in this file before WriteSection is even reached, and the neutral no-op
// WriteSection says nothing about - lands on exactly one PLUGIN-TOGGLE line
// with its outcome word: written / unchanged (both under the written
// umbrella), queued, or refused. WriteSection's own SECTION-WRITE lines
// carry the direction detail on top. The line goes through the same
// m.auditf sink every adjudicated write uses; a nil Audit sink falls through
// to the process log - the line is never silenced, only re-homed.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// SetPluginEnabled writes plugins.<id>.enabled through the protected
// section-write entry. id is the [plugins.<id>] table the config declares
// (verbatim, no case folding - the id is also the plugins\<id> directory
// name, C14's charset); on is the switch's target value; src names who asked
// (SourcePanel through the settings route, SourceNative through the CLI) and
// origin names the route, for the audit line. The returned slice is the
// receipt: the key path(s) the file actually changed by, empty when nothing
// was written (queued, refused, or unchanged).
//
// The body lives in setPluginEnabled215 so the audit line below fires on
// every path back, including the refusals this file hands out itself.
func (m *Manager) SetPluginEnabled(id string, on bool, src WriteSource, origin string) ([]string, error) {
	written, queued, err := m.setPluginEnabled215(id, on, src, origin)

	// The toggle's own audit line, every outcome. A refusal in this file has
	// no WriteSection receipt to lean on, so the outcome word here is the
	// only record the attempt happened at all. The three outcome words the
	// ticket's AC#3 names: written (the write landed, plus unchanged - the
	// attempt that wanted the value the file already carries), queued (the
	// loosening is now an L2 card), and refused (every refusal, in this file
	// or in WriteSection).
	outcome := "refused"
	switch {
	case err == nil && queued:
		outcome = "queued"
	case err == nil && written != nil:
		outcome = "written"
	case err == nil:
		outcome = "unchanged"
	}
	auditKey := "plugins.<none>.enabled"
	if trimmed := strings.TrimSpace(id); trimmed != "" {
		auditKey = "plugins." + trimmed + ".enabled"
	}
	m.auditf("config: PLUGIN-TOGGLE key=%s on=%v source=%s origin=%s outcome=%s",
		auditKey, on, src, origin, outcome)

	return written, err
}

// setPluginEnabled215 is the write body. See SetPluginEnabled's header. The
// queued return distinguishes a handed-off loosening (nothing written, one
// pending card) from an unchanged no-op.
func (m *Manager) setPluginEnabled215(id string, on bool, src WriteSource, origin string) (written []string, queued bool, err error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, false, observe.New(observe.ClassConfig,
			"config: plugins.<id>.enabled: 写入必须点名插件 id（[plugins.<id>] 的 <id> 本身）；不能新建行，也不能盲写")
	}
	if !src.valid() {
		return nil, false, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: plugins.%s.enabled: 写入来源 %q 不是已定义的来源（panel|native）；谁在开关插件必须落在审计里", id, src))
	}
	cur := m.Config()
	if cur == nil {
		return nil, false, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: plugins.%s.enabled: Manager 里没有已加载配置，无处落写", id))
	}
	if _, ok := cur.Plugins.Entries[id]; !ok {
		return nil, false, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: plugins.%s.enabled: 配置里没有这个插件 id（现声明的有 %d 个）；设置写不新建行，要先在 config.toml 声明 [plugins.<id>]",
			id, len(cur.Plugins.Entries)))
	}

	key := "plugins." + id + ".enabled"
	res, err := m.WriteSection(SectionWrite{
		Section: "plugins",
		Key:     key,
		Source:  src,
		Origin:  origin,
		Check: func(c *Config) error {
			// The FILE base must still declare the row (ticket 226's rule: the
			// verdict is judged against the file, and the house rule is that a
			// settings write cannot create one). The memory-side check above
			// caught the typo case; this one catches a hand-edited file that
			// dropped the row between the two reads.
			if _, ok := c.Plugins.Entries[id]; !ok {
				return observe.New(observe.ClassConfig, fmt.Sprintf(
					"config: %s: 盘上的 config.toml 已不再声明插件 %q；设置写不新建行", key, id))
			}
			return nil
		},
		Apply: func(c *Config) error {
			entry := c.Plugins.Entries[id]
			entry.Enabled = on
			c.Plugins.Entries[id] = entry
			return nil
		},
	})
	if err != nil {
		return nil, false, err
	}
	if res.PendingID != "" {
		// The receipt names what the FILE changed by. A queued loosening
		// changed nothing - WriteSection's Keys on a queued loosening carry
		// the would-be loosen key paths for the card, and handing those to a
		// settings receipt would read as "written".
		return nil, true, nil
	}
	if !res.Written {
		return nil, false, nil
	}
	return res.Keys, false, nil
}
