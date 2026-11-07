package main

// Ticket 65's CLI half: `wisp config get|set|list`, the operator's own surface
// for the ball style switch the owner asked for (默认玻璃球，切 Jarvis 才换球
// 的风格). The Go-side four segments already existed (schema key, SetBallStyle,
// the panel write gate, the applier seam, the boot seed), but every one of them
// hangs off the resident process or the panel page - and the settings page has
// no ball-style row yet (frontend/** is the external agent's, Q-51). This file
// is the path the owner can run TODAY, from any console, while wisp is or is
// not running.
//
// THE SHAPE IS THE MAINSTREAM ONE, ON PURPOSE (铁律 ⓪⓪⓪): `git config
// user.name` / `kubectl config set` - get <key>, set <key> <value>, list.
// No flags to remember, no new vocabulary: the key path is the dotted
// `ball.style` the guarded setter already owns, and the values are the
// schema's own glass|jarvis words. The subcommand dispatch, the exit codes
// (2 usage/environment, 1 operational refusal, 0 success) and the layout
// resolution follow the `wisp secret` leg line for line.
//
// THE GUARD IS NOT DUPLICATED HERE. SetBallStyle already refuses anything
// outside the schema's ballStyles list by name, before a byte moves
// (internal/config/settings_65.go); the CLI calls it and passes its refusal
// through. The CLI adds exactly one guard of its own: the KEY roster - this
// command surface owns ball.style today, and an unknown key is refused by
// name rather than guessed at, the same closed-roster shape the panel's
// knownWritableField union uses.
//
// THE TIER IS SAID OUT LOUD, NOT GUESSED. A CLI write lands in config.toml;
// the running resident's reload watchdog will take the new value into memory
// (hot tier), but the family swap itself hangs on exactly two hooks - the
// panel write path's applyBallStyle65 (in-process) and the boot seed
// (resident_windows.go) - so a CLI write cannot switch the live ball of
// ANOTHER process. The receipt says so: 下次启动生效；正在运行的 wisp 重启后才
// 换球. Claiming "立即生效" here would be the one lie a receipt must never
// tell.
//
// THE DEFAULT IS NOT WRITTEN HERE EITHER. A missing config.toml answers `get`
// with config.NewDefaults().Ball.Style - the schema's own default:"glass" tag
// through its exported table (internal/config/defaults.go), the same table
// first-run creation writes from - never a local literal. A file that exists
// answers with the effective value the loader produces (defaults applied),
// which is what this process would act on.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
)

// configCliKey65 is the one dotted key path this CLI surface owns today. The
// roster is closed on purpose: get/set name it, everything else is refused
// with the roster in the refusal, so a typo'd key can never write a second
// path around the guarded setters.
const configCliKey65 = "ball.style"

// configCliUsage65 is the subcommand block printed by `wisp config` and by
// every usage refusal this leg answers with.
const configCliUsage65 = `wisp config - read and write config.toml keys from the console

Usage:
  wisp config get <key>     print the effective value (schema default applied
                            when the key is absent from config.toml)
  wisp config set <key> <value>
                            write through the same guarded setter the settings
                            page uses; the receipt names the file, the key and
                            the tier the write really lands at
  wisp config list          the keys this CLI owns, each with its current value
                            and its schema vocabulary

Keys owned today: ball.style (取值 glass|jarvis，默认 glass).
Environment: the active WISP_ENV is resolved like ` + "`wisp secret`" + ` (SPEC-03 §5).
`

// configIO is the process surface of this leg: the same injection seam
// commandsIO/modelsIO use, so the tests drive the real entry point.
type configIO struct {
	stdout io.Writer
	stderr io.Writer
}

// cmdConfig is the `wisp config` entry point (main dispatch): resolve the env,
// resolve its data root the way cmdSecret does, then dispatch the subcommand
// against <dataDir>\config.toml - the same join run.go:513 and firstrun.go
// make over the resolved root.
func cmdConfig(args []string, cio configIO) int {
	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: %v\n", err)
		return 2
	}
	layout, err := resolveSecretLayout(env)
	if err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: %v\n", err)
		return 2
	}
	cfgPath := filepath.Join(layout.DataDir, configFileName)
	if len(args) == 0 {
		fmt.Fprint(cio.stderr, configCliUsage65)
		return 2
	}
	switch args[0] {
	case "get":
		key := ""
		if len(args) > 1 {
			key = args[1]
		}
		return runConfigGet65(cio, cfgPath, key)
	case "set":
		key, rest := "", []string(nil)
		if len(args) > 1 {
			key = args[1]
			rest = args[2:]
		}
		return runConfigSet65(cio, layout.DataDir, cfgPath, key, rest)
	case "list":
		return runConfigList65(cio, cfgPath)
	default:
		fmt.Fprintf(cio.stderr, "wisp config: unknown subcommand %q\n\n", args[0])
		fmt.Fprint(cio.stderr, configCliUsage65)
		return 2
	}
}

// runConfigGet65 answers with the effective value of one key. A missing
// config.toml is not an error: the schema's default table IS the effective
// value before the file exists, and the note on stderr keeps "the file said
// this" and "nobody has written a file yet" two different facts.
func runConfigGet65(cio configIO, cfgPath, key string) int {
	if key == "" {
		fmt.Fprint(cio.stderr, configCliUsage65)
		return 2
	}
	if key != configCliKey65 {
		fmt.Fprintf(cio.stderr, "wisp config: %q 不在本 CLI 的键名册上（今天只开放 %s）\n", key, configCliKey65)
		return 2
	}
	style, ok, err := loadBallStyle65(cfgPath)
	if err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: 读取 %s 失败：%v\n", cfgPath, err)
		return 1
	}
	if !ok {
		fmt.Fprintln(cio.stdout, config.NewDefaults().Ball.Style)
		fmt.Fprintf(cio.stderr, "wisp config: %s 尚不存在；以上是 schema 默认值（default 标签），不是文件里的值\n", cfgPath)
		return 0
	}
	fmt.Fprintln(cio.stdout, style)
	return 0
}

// runConfigSet65 writes one key through the guarded setter. A missing
// config.toml is created first through ensureFirstRunConfig - the ONE first-run
// entry (firstrun.go) - because an explicit `set` IS the ask that entry exists
// for; a `get` never creates. The setter's enum refusal passes through
// unchanged: the value is named, the vocabulary is named, and the file is
// byte-identical (refusal before a byte moves, settings_65.go).
func runConfigSet65(cio configIO, dataDir, cfgPath, key string, rest []string) int {
	if key != configCliKey65 || len(rest) != 1 {
		if key != configCliKey65 {
			fmt.Fprintf(cio.stderr, "wisp config: %q 不在本 CLI 的键名册上（今天只开放 %s）\n", key, configCliKey65)
		}
		fmt.Fprint(cio.stderr, configCliUsage65)
		return 2
	}
	value := rest[0]
	if _, err := ensureFirstRunConfig(dataDir, cio.stderr); err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: %v\n", err)
		return 1
	}
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: %s 读不进来（先修好它，CLI 不替你改坏文件）：%v\n", cfgPath, err)
		return 1
	}
	if _, err := mgr.SetBallStyle(value); err != nil {
		fmt.Fprintf(cio.stderr, "wisp config: 拒绝写入：%v\n", err)
		return 1
	}
	fmt.Fprintf(cio.stdout, "已写入 %s：%s = %s\n", cfgPath, configCliKey65, value)
	fmt.Fprintln(cio.stdout, "生效档位：重启档——下次 wisp 启动时换球；正在运行的 wisp 进程不会立刻换球"+
		"（换球挂在面板写路径与本进程开机种子上，CLI 写的是文件）。要马上看到：重启 wisp。")
	return 0
}

// runConfigList65 prints the closed roster: one row per owned key, the current
// effective value and the schema's own vocabulary (BallStyleOptions asks the
// schema's ballStyles list, never a copy).
func runConfigList65(cio configIO, cfgPath string) int {
	var view config.BallStyleView
	if _, statErr := os.Stat(cfgPath); statErr == nil {
		mgr, err := config.NewManager(cfgPath, nil)
		if err != nil {
			fmt.Fprintf(cio.stderr, "wisp config: %s 读不进来：%v\n", cfgPath, err)
			return 1
		}
		view = mgr.BallStyleOptions()
	} else if errors.Is(statErr, fs.ErrNotExist) {
		// No file yet: the schema vocabulary stands and the value is the
		// default table's, both said through the schema's own faces.
		view = (*config.Manager)(nil).BallStyleOptions()
		view.Style = config.NewDefaults().Ball.Style
		fmt.Fprintf(cio.stderr, "wisp config: %s 尚不存在；以下值是 schema 默认值\n", cfgPath)
	} else {
		fmt.Fprintf(cio.stderr, "wisp config: 读取 %s 失败：%v\n", cfgPath, statErr)
		return 1
	}
	fmt.Fprintf(cio.stdout, "%s = %s    取值: %s\n",
		configCliKey65, view.Style, strings.Join(view.Vocabulary, "|"))
	fmt.Fprintln(cio.stdout, "本 CLI 今天只开放这一枚键。")
	return 0
}

// loadBallStyle65 reads the effective ball.style off the disk the way
// ballStyleSeed65 reads it at boot: config.LoadFile over the data dir, defaults
// applied by the loader. (style, true) is the file's own answer; (style, false)
// says the file is absent - distinct from unreadable, which is the error.
func loadBallStyle65(cfgPath string) (string, bool, error) {
	if _, err := os.Stat(cfgPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", false, nil
		}
		return "", false, err
	}
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil {
		return "", false, err
	}
	return c.Ball.Style, true, nil
}
