//go:build windows

package main

// Ticket 247: the capture leg wired into the resident process.
//
// What this file is. Until it existed, internal/audio had ZERO non-test
// importers in the whole repository (ticket 现量 1, re-measured before this leg
// started: the grep ruler answered 0), and the only production producer of
// Ball.SetAudioLevel was cmd/balldebug feeding a "synthetic audio envelope"
// typed on the command line (现量 2). So the whole capture stack - real WASAPI
// shared-mode capture, the D38d bounded frame path, the gate, the level scale of
// ticket 241 - was finished code nobody ran. This file runs it.
//
// Shape, all of it ruled by the orchestrator in the ticket's "AC#0" section and
// not chosen here:
//
//   -落点 (form 甲 plus the 乙-2 exit): the capture goroutine is assembled by
//     this resident leg, cmd/wisp. The one new package-level dependency edge is
//     cmd/wisp -> internal/audio. internal/audio still does not import
//     internal/ball, and no edge audio -> ball was opened.
//   - owner (P4 甲): the thread is booked through audio.SpawnCapture into the
//     registry this process already owns (rt.Registry, observe/goroutine.go:262),
//     never a bare go func (D22 ban #1). One registry per source: the injected
//     registry REPLACES the process default rather than joining it.
//   - level (P6 甲 + P8 甲): computed on the capture side with audio.FrameLevel
//     inside the existing audio-capture goroutine and handed out as ONE float32
//     through the callback below. The D38b resident roster stays at six names;
//     no second thread, no reused name (a reused name would double the Resident
//     count and dilute PLAN.md:2831's "6 resident goroutines").
//   - privacy (P1 甲): the gate is always mounted, with
//     audio.WithStartMuted(c.Audio.MicMutedDefault), so [audio]
//     mic_muted_default - a key with zero production readers before this line -
//     is what decides whether the device opens at boot. voice.enabled=false
//     means the collector is not built at all. No default value moved.
//   - degraded device (P5 甲): an occupied / denied / absent device is said out
//     loud with its observe error class and the guidance text the audio package
//     already carries, and the boot continues. It pushes NO state: D43's
//     EvAudioDeviceLost edge starts only from Listening (internal/statemachine/
//     table.go:97), so a boot-time failure has no legal edge, and inventing one
//     would be a contract change. It also does not refuse to start: ticket 128
//     settled exactly one refusal condition (no data root).
//   - stop (AC#5): the hook slot proc.StepStopAudio already exists in the
//     frozen D38(e) order at step 4; this registers into it. Ten steps, one
//     order, nothing new. WASAPIMicrophone.Stop carries its own bounded 2s join.
//
// What this file does NOT claim (AC#10): handing the level to the ball is not
// "the user sees the ball breathe". prototypeVisuals defaults to off
// (internal/ball/statevisual.go:102) and Ball.SetAudioLevel returns immediately
// while it is off (internal/ball/liquid_windows.go:52). Turning that default is
// ticket 68 AC#2's decision, not this leg's.
//
// What this file does NOT have: a PCM consumer. The C8 frame channel is created
// with the seam's own bounded window and nothing drains it, because the ASR/KWS
// side (internal/speech) is unwritten and 唤醒词／ASR／TTS are explicitly outside
// this ticket. Every frame is therefore dropped AFTER the level has been read,
// counted in the D38d meter and reported at most once a second by
// meter.push's own throttle - the honest interim shape, named in the posture
// line below so the boot report does not pretend otherwise.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync/atomic"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// residentAudio is this process's hold on the capture leg: the gate it started
// (or the reason it did not), the counter the boot and exit lines read, and the
// teardown the D38(e) step-4 hook runs. It is nil-safe everywhere, like
// residentBall.
type residentAudio struct {
	// gate and mic are nil when the collector was never built (voice disabled,
	// or the config read refused to guess).
	gate *audio.HalfDuplexGate
	// mic is behind the captureSource seam: the reads this leg makes of it are
	// Stats() (the posture switch) and Err() (the error class), so AC#6's
	// device failures can be assembled without unplugging hardware.
	mic    captureSource
	cancel context.CancelFunc

	// The two config facts this leg acted on, recorded so the boot line and any
	// later reading of this handle state the same privacy posture the assembly
	// saw. Nothing here writes a default: both values come off the file, or off
	// the compiled default table when the file could not be read.
	voiceEnabled bool
	mutedAtBoot  bool
	// started records that the gate handed the device to the inner source and
	// capture is actually running. posture() reads it so a leg that built a
	// gate but never opened a microphone cannot print "running".
	started bool

	// levels is how many scalars this leg handed to the ball seam. It is what
	// makes "the level path is live" a reading instead of a comment; the exit
	// line prints it next to the frame counters.
	levels atomic.Uint64

	verdict string
}

// levelOut is the one crossing this leg performs into the ball: a count, then
// the callback the caller built (in runResident that is
// residentBall.setAudioLevel, a nil-safe method value). AC#3 is a capability
// judgement and float32 is the only data that crosses.
func (ra *residentAudio) levelOut(out func(float32)) func(float32) {
	return func(level float32) {
		if ra != nil {
			ra.levels.Add(1)
		}
		if out != nil {
			out(level)
		}
	}
}

// setAudioLevel is the ball host's side of the one-scalar handoff. It lives
// here, not in resident_ball_windows.go, so this ticket's diff stays in one
// file; it is nil-safe for the same reason every other residentBall method is:
// a machine whose window never came up must still be handed a level without
// anybody dereferencing a nil window.
func (rb *residentBall) setAudioLevel(level float32) {
	if rb == nil || rb.b == nil {
		return
	}
	rb.b.SetAudioLevel(level)
}

// startResidentAudio assembles the capture leg at boot, reading [voice] and
// [audio] off this process's own data root. It never stops the boot: every
// branch that cannot open a microphone returns a handle with its verdict
// already written, and the caller prints that verdict.
func startResidentAudio(rt *proc.Runtime, out func(float32)) *residentAudio {
	return buildResidentAudio(rt, rt.Layout.DataDir, out)
}

// captureSource is the capture side this leg drives: the C8 seam plus the two
// telemetry reads the posture is built from. *audio.WASAPIMicrophone is the
// production answer; the seam exists so AC#6's three device failures (occupied
// / permission denied / no device) can be shown at THIS assembly level - "the
// process still boots and says the loss out loud" is a claim about the resident
// leg, not about internal/audio's own hotplug tests.
type captureSource interface {
	audio.AudioSource
	Stats() audio.Stats
	Err() error
}

// newRealCaptureSource is the production factory.
func newRealCaptureSource(opts ...audio.CaptureOption) captureSource {
	return audio.NewWASAPIMicrophone(opts...)
}

// buildResidentAudio is the same assembly with the config directory named: the
// privacy posture this leg decides (does a double click open the microphone?)
// has to be testable against a config.toml the test wrote, not against the
// developer's own data root. The runtime handle still supplies the registry and
// the shutdown-hook slot, so nothing about the production shape is bypassed.
func buildResidentAudio(rt *proc.Runtime, dataDir string, out func(float32)) *residentAudio {
	return assembleCapture(rt, dataDir, out, newRealCaptureSource)
}

func assembleCapture(rt *proc.Runtime, dataDir string, out func(float32),
	newSource func(opts ...audio.CaptureOption) captureSource,
) *residentAudio {
	ra := &residentAudio{}

	// The per-fresh-read shape this package already uses for hot-tier views
	// (hotCfg258 above, panelGeometrySource since ticket 255 AC#4): the task
	// pipeline's Manager does not exist at this point of the boot, so [voice]
	// and [audio] are read off the file rather than held.
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	cfgSource := cfgPath
	if err != nil || c == nil {
		// Not a guess at a value: config.NewDefaults is the repository's own
		// compiled default table, the same one `wisp run` writes into the first
		// config.toml (firstrun.go:92). It is named out loud below because which
		// table decided the privacy posture is a fact the operator has to be
		// able to see. Its [audio] half is the conservative one:
		// mic_muted_default = true, so this fallback cannot open a device that
		// the missing file would have kept shut.
		c = config.NewDefaults()
		cfgSource = "compiled defaults (config.NewDefaults)"
		slog.Warn("audio: config unreadable at boot; the capture leg is built from the compiled default table",
			"path", cfgPath, "err", configErrText(err, c),
			"mic_muted_default", c.Audio.MicMutedDefault, "voice_enabled", c.Voice.Enabled)
	}

	ra.mutedAtBoot = c.Audio.MicMutedDefault
	ra.voiceEnabled = c.Voice.Enabled

	// P1 甲, second half: voice off means no collector is constructed at all -
	// not "constructed and idle", so there is no device handle to hold and no
	// thread to join.
	if !c.Voice.Enabled {
		ra.verdict = "采集腿未构造：[voice] enabled=false（配置来源 " + cfgSource + "）"
		slog.Info("audio: capture leg not built", "reason", "voice.enabled=false", "config_source", cfgSource)
		fmt.Printf("wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 %s）：球不会收到任何电平，本进程其余部分照常\n", cfgSource)
		return ra
	}

	mic := newSource(
		audio.WithCaptureRegistry(rt.Registry),
		audio.WithLevelSink(ra.levelOut(out)),
	)
	gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))

	ctx, cancel := context.WithCancel(context.Background())
	buf := audio.NewBoundedFrames()
	ra.mic, ra.gate, ra.cancel = mic, gate, cancel

	// The gate's Start returns nil even when the inner open failed (the failure
	// is metered inside openInnerLocked and surfaced through Stats), so the
	// posture below is read off the gate and the counters, not off this error.
	startErr := gate.Start(ctx, buf)
	if startErr != nil {
		ra.verdict = "采集腿启动失败：" + startErr.Error()
		slog.Error("audio: capture leg refused to start", "err", startErr.Error())
		fmt.Printf("wisp: 麦克风采集腿启动失败（%v）：本进程仍带球与任务管线常驻，不因此退场\n", startErr)
		cancel()
		return ra
	}

	// P5 甲: the device posture, said loudly off measured state.
	lastErr := mic.Stats().LastError
	switch {
	case gate.Muted():
		// The default档 lands here: mic_muted_default=true means the inner
		// source was never started, so wasapi_windows.go's synchronous device
		// open never ran and no microphone is capturing.
		ra.verdict = fmt.Sprintf("采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）", cfgSource)
		slog.Info("audio: capture armed but muted at boot; device not opened",
			"config_source", cfgSource, "path", string(gate.Path()))
	case lastErr != "":
		// The class comes off the error object the capture thread posted
		// (WASAPIMicrophone.Err), not off the meter's string: the string is the
		// guidance the operator reads, the class is what D37 says it is.
		class := string(observe.ClassAudioDevice)
		if devErr := mic.Err(); devErr != nil {
			var oe *observe.Error
			if errors.As(devErr, &oe) {
				class = string(oe.Class)
			}
			lastErr = devErr.Error()
		}
		ra.verdict = fmt.Sprintf("麦克风不可用（分类 %s）：%s；球不会收到电平，本进程继续跑", class, lastErr)
		slog.Error("audio: capture device unavailable",
			"class", class, "err", lastErr, "posture", "boot continues, no state pushed")
		fmt.Printf("wisp: 麦克风不可用（错误分类 %s）：%s\n", class, lastErr)
		fmt.Printf("wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，" +
			"启动期没有；票 128 只定了「没有数据根」这一种拒绝启动\n")
	case !gate.Open():
		ra.verdict = "采集腿在跑但设备未交接（gate 未 open）：球不会收到电平"
		slog.Warn("audio: gate reports not open without a device error", "posture", ra.verdict)
	default:
		ra.started = true
		ra.verdict = fmt.Sprintf("采集腿在跑：真麦克风 PCM 进，%s 每帧一枚 float32 交给球（路径 %s）",
			audio.FrameDuration.String(), string(gate.Path()))
		slog.Info("audio: capture leg running", "path", string(gate.Path()), "frame", audio.FrameDuration.String())
	}

	// D38(e) step 4 already has a slot (proc.StepStopAudio); this is its first
	// production registrant. The ten steps and their order are untouched and no
	// eleventh step was invented. WASAPIMicrophone.Stop carries its own bounded
	// 2s join, so the step is bounded by the module that owns the thread.
	// A leg that built nothing does NOT register: the audit trail keeps saying
	// "skipped" for a step with no owner instead of claiming a teardown that
	// would do nothing.
	if err := rt.RegisterShutdownHook(proc.StepStopAudio, func(context.Context) error {
		return ra.stop()
	}); err != nil {
		slog.Error("proc: 第 4 步（停采集线程）未能注册为钩子", "step", int(proc.StepStopAudio), "err", err)
		fmt.Printf("wisp: 退出序列第 4 步未注册（%v）：采集线程不会被有序关闭\n", err)
	}
	return ra
}

// stop is D38(e) step 4's body: cancel the consumer ctx, then hand the gate
// the stop (which stops the inner source, and whose own join is bounded at 2s
// inside WASAPIMicrophone.Stop). Returns the device error this leg ends with,
// if any, so the shutdown record carries it instead of a clean-looking nil.
func (ra *residentAudio) stop() error {
	if ra == nil || ra.gate == nil {
		return nil
	}
	if ra.cancel != nil {
		ra.cancel()
	}
	err := ra.gate.Stop()
	st := ra.mic.Stats()
	line := fmt.Sprintf("audio capture leg stopped: levels_delivered=%d frames_sent=%d frames_dropped=%d reopens=%d",
		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.Reopens)
	if st.LastError != "" {
		line += " last_error=" + st.LastError
	}
	slog.Info(line)
	fmt.Printf("wisp: %s\n", line)
	if err != nil {
		slog.Error("audio: capture leg stop returned an error", "err", err.Error())
		return err
	}
	return nil
}

// posture is the boot-report sentence. It is built from the state the assembly
// recorded, never from a hope, and it names the missing PCM consumer because
// that is the half of the voice chain this ticket does not own.
func (ra *residentAudio) posture() string {
	if ra == nil {
		return "采集腿未装配：本进程没有任何麦克风路径"
	}
	if !ra.started {
		return ra.unrunningClaim()
	}
	return ra.verdict + "；帧消费侧（ASR/KWS，internal/speech 一块没写）不属于本票，" +
		"帧按 D38d 计数丢弃并由 meter 限流告警"
}

// unrunningClaim states the no-collector shapes in the words that are true for
// them, so the boot line cannot be read as "the mic is merely muted".
func (ra *residentAudio) unrunningClaim() string {
	if ra != nil && ra.verdict != "" {
		return ra.verdict
	}
	return "采集腿未装配：本进程没有任何麦克风路径"
}

// configErrText keeps the warn record honest about a read that produced no
// value at all (err non-nil) versus a read that returned nothing usable.
func configErrText(err error, c *config.Config) string {
	if err != nil {
		return err.Error()
	}
	if c == nil {
		return "config.LoadFile returned a nil config with a nil error"
	}
	return ""
}
