package config

// Ticket 65's settings-write additions: [ball].style, the orb style family
// (appearance switch the owner asked for: 默认玻璃球，切 Jarvis 才换球的风格).
//
// THE KEY ALREADY EXISTED. The [ball] section, its style field, the glass|jarvis
// vocabulary (schema.go's ballStyles), the enum refusal in validate.go and the
// glass default all landed with the ticket-65 impl leg (commit 30d8f2e38, pinned
// by ball_style_65_test.go). What was missing was everything the settings page
// needs to TOUCH the key: a setter on the guarded path and a read face for a
// picker. Both live here, and neither invents a vocabulary of its own - both
// ask the schema's ballStyles list, the same list validate.go judges a load by.
//
// SAME SHAPE AS THE SETTERS BEFORE IT (settings_320.go's header says it best).
// One owned key path, writeOneKey, validate() on the base the write would
// produce, refusal before a byte moves, receipt names the key path the FILE
// changed by. The enum is checked HERE and not only at load, the same ruling
// ticket 320 made for provider_billing: a third spelling must come back as a
// sentence naming the field, not as a whole-file validation failure the
// operator has to localize.
//
// THE TIER IS NOT DECIDED HERE. Unlike every [llm] field, ball.style can be
// applied to the RUNNING process (ball.ApplyStyle swaps the process-wide family
// and re-navigates the orb page), but only when this process actually hosts the
// ball - the applier half lives in cmd/wisp (ball_style_65_windows.go), where
// the ball exists, and the settings write leg answers with the measured tier
// (now when the applier ran, restart when there is nothing to apply it to).
// internal/config stays out of that question on purpose: it must not import
// internal/ball, and a tier guessed here would be a lie one of the two hosts
// would have to repeat.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// keyPathBallStyle is the one dotted key path SetBallStyle owns.
const keyPathBallStyle = "ball.style"

// SetBallStyle writes ball.style through the same guarded key-by-key path every
// settings setter uses. The value is checked against the schema's own vocabulary
// BEFORE the write; anything else is refused with the field's name and the
// vocabulary, leaving config.toml byte-identical.
func (m *Manager) SetBallStyle(style string) ([]string, error) {
	v := strings.TrimSpace(style)
	if !isBallStyle(v) {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: ball.style: %q 不在枚举内（%s）；这一枚只接受这几种写法",
			v, vocabList(ballStyles)))
	}
	return m.writeOneKey(keyPathBallStyle, func(base *Config) bool {
		base.Ball.Style = v
		return true
	}, nil)
}

// isBallStyle asks the schema's own list, never a copy of the two words.
func isBallStyle(v string) bool {
	for _, s := range ballStyles {
		if s == v {
			return true
		}
	}
	return false
}

// BallStyleView is the read face a style picker needs: the family the live
// config names right now, and the vocabulary it may offer. Reads the in-memory
// config (Manager.Config()) the same way ChatSwitchOptions does - the file and
// memory are kept equal by the guarded writer, and a reader must answer with
// what this process would act on.
type BallStyleView struct {
	// Style is the current ball.style value. Empty only when no Manager backs
	// the read (the caller answers "not configured", never a made-up family).
	Style string
	// Vocabulary is the schema's own glass|jarvis list, in schema order, so a
	// picker offers exactly what load accepts.
	Vocabulary []string
}

// BallStyleOptions is the production reader for the panel's appearance picker
// (the reader-gated shape composer.go documents). A nil Manager answers an
// empty Style with the schema vocabulary still intact.
func (m *Manager) BallStyleOptions() BallStyleView {
	v := BallStyleView{Vocabulary: append([]string(nil), ballStyles...)}
	if m == nil {
		return v
	}
	c := m.Config()
	if c == nil {
		return v
	}
	v.Style = c.Ball.Style
	return v
}
