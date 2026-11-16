//go:build windows

package main

// Ticket 313's cloud-voice-chain wiring leg: voice.cloud_asr_chain and
// voice.cloud_tts_chain stop being "a value the settings page can save" and
// become "the engines the resident voice loop actually loads". The gap this
// closes was measured before the leg moved anything: speech.NewCloudASRChain /
// NewCloudTTSChain had ZERO call sites outside internal/speech's own tests, so
// the two chains the UI writes were read by no runtime path at all - the
// honest-placeholder copy on the voice settings page said exactly that.
//
// WHERE THE WIRING SITS. The composition root (buildResidentVoice) resolves
// the two chains off the same boot config.LoadFile feed that has carried every
// voicePlan field since ticket 28, and hands the resolved endpoints to the
// residency; the residency's loader (resident_voice_engines_windows.go) builds
// the chain engines when a session scope asks for the ASR or TTS family. The
// chains are TierReload config (internal/config/tiers.go), so "save then
// restart" is the contract the UI copy now states - a boot feed is exactly
// that contract, not a shortcut around it.
//
// WHAT HAPPENS ON A BAD CHAIN. A chain element the catalog cannot resolve
// (unknown pair, disabled model, missing capability bit, unresolvable
// api_key_ref) refuses the WHOLE chain here, and the refusal is a named
// problem + a boot line that does not claim the cloud chain - never a silent
// skip and never a boot failure. A voice machine whose cloud chain is misfed
// must still listen and speak locally, which is the terminal-fallback design
// the cascade itself is built around; what it must never do is say "云端听写
// stepfun/..." while holding a local-only loader. So: resolve errors become
// notes, resolved endpoints become the verdict line, and the two lists can
// never disagree with what the loader will do.
//
// KEYS. Endpoints carry the plaintext key ONLY inside speech.CloudEndpoint.APIKey
// (resolved here through the process's own DPAPI store, secret.Store - the one
// this process's settings route uses). It never reaches a log line, an error
// detail, or the verdict: every string this file prints is ChainKey/baseURL
// shaped, which CloudEndpoint.String already guarantees for the endpoint.

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/speech"
)

// resolveVoiceCloudChains turns the boot config's two voice cloud chains into
// ordered, fully resolved endpoints. An empty chain list on both sides is the
// "local cascade only" configuration and returns nothing at all - no store is
// opened, no note is written, the boot line stays exactly what it always said.
// A configured chain that refuses to resolve is a NAMED note (the element and
// the reason, key-free) plus an empty endpoint list for that direction: the
// boot says "云端听写链未装配", the local cascade keeps the direction, and
// nothing silently pretends the cloud head exists.
func resolveVoiceCloudChains(dataDir string, c *config.Config) (asrEps, ttsEps []speech.CloudEndpoint, notes []string) {
	if len(c.Voice.CloudASRChain) == 0 && len(c.Voice.CloudTTSChain) == 0 {
		return nil, nil, nil
	}
	keys, err := secret.NewStore(dataDir)
	if err != nil {
		return nil, nil, []string{fmt.Sprintf(
			"云端语音链未装配：凭据存储打不开（%v），两条链都退回本地级联", err)}
	}
	if len(c.Voice.CloudASRChain) > 0 {
		eps, err := speech.ResolveCloudChain(c, speech.ChainASR, keys)
		if err != nil {
			notes = append(notes, "云端听写链未装配（"+err.Error()+"）：听写走本地级联")
		} else {
			asrEps = eps
		}
	}
	if len(c.Voice.CloudTTSChain) > 0 {
		eps, err := speech.ResolveCloudChain(c, speech.ChainTTS, keys)
		if err != nil {
			notes = append(notes, "云端朗读链未装配（"+err.Error()+"）：朗读走本地级联")
		} else {
			ttsEps = eps
		}
	}
	return asrEps, ttsEps, notes
}

// voiceCloudChainNames renders the resolved chains the way the boot verdict and
// the evidence table quote them: "provider/model" per element, first element
// first (it is the head the cascade tries before anything else). Keys have no
// representation in this spelling by construction.
func voiceCloudChainNames(eps []speech.CloudEndpoint) string {
	out := ""
	for i, ep := range eps {
		if i > 0 {
			out += ", "
		}
		out += speech.ChainKey(ep.Provider, ep.Model)
	}
	return out
}

// voiceCloudChainOptions is the observer group the resident chains report
// through. Both callbacks run on the voice chain's owner goroutine (D38b) and
// are log lines only: failover hops and the privacy posture (is audio leaving
// the machine) are the two things a reader of the log must be able to answer,
// and both carry the same log-safe vocabulary the events already enforce
// (VoiceFailoverEvent.Detail is contractually never the key, never the audio).
// Quota and usage stay unwired here on purpose: the cascade treats a nil gate
// as "no quota configured" and a nil sink as "no metering", which is this
// leg's honest posture until the C23 voice rows get their own ticket.
func voiceCloudChainOptions(registry *observe.Registry) speech.ChainOptions {
	return speech.ChainOptions{
		// The per-utterance transport goroutines spawn through the process's
		// own registry (observe.Registry.Spawn), so the watchdog sees one
		// picture; nil falls back to observe.Default inside the cascade.
		Registry: registry,
		Callbacks: speech.ChainCallbacks{
			OnFailover: func(ev speech.VoiceFailoverEvent) {
				slog.Info("voice: 云端语音链切换",
					"chain", string(ev.Chain), "from", ev.From, "to", ev.To,
					"cause", ev.Cause, "detail", ev.Detail,
					"attempt", ev.Attempt, "total", ev.Total)
			},
			OnProviderState: func(st speech.VoiceProviderState) {
				slog.Info("voice: 云端语音链活跃引擎",
					"chain", string(st.Chain), "active", st.Active, "cloud", st.Cloud)
			},
		},
	}
}

// voiceCloudAsrFamily is what the slot actually holds for an ASR epoch with a
// cloud chain: the cascade PLUS the local sherpa session it falls back to.
// The cascade deliberately does not own its fallback (cloudchain.go's ownership
// rule: the fallback belongs to whoever loaded it) - but here the composition
// root is the loader, so "whoever loaded it" is the same family, and the slot
// expects ONE engine whose Close releases the family. This wrapper joins the
// two halves so the family's memory is max(ASR, TTS) at slot level and not
// max-plus-a-leaked-sherpa-session; C11 idempotence on both halves makes the
// joined Close safe to call as many times as the slot's bookkeeping asks.
type voiceCloudAsrFamily struct {
	*speech.CloudASRChain
	local speech.AsrEngine
}

func (e *voiceCloudAsrFamily) Close() error {
	err := e.CloudASRChain.Close()
	if lerr := e.local.Close(); lerr != nil {
		err = errors.Join(err, lerr)
	}
	return err
}

// voiceCloudTtsFamily is the TTS twin of voiceCloudAsrFamily. SampleRate
// follows the cascade's own rule (the active head's rate, the fallback's on a
// local hop), so the render side needs no family awareness at all.
type voiceCloudTtsFamily struct {
	*speech.CloudTTSChain
	local speech.TtsEngine
}

func (e *voiceCloudTtsFamily) Close() error {
	err := e.CloudTTSChain.Close()
	if lerr := e.local.Close(); lerr != nil {
		err = errors.Join(err, lerr)
	}
	return err
}

// buildCloudAsrEngine is the residency's ASR factory when a cloud chain is
// configured: the local sherpa session is built as the terminal fallback when
// its model set exists (a local load failure is a named problem and a
// cloud-only chain, never a failed utterance path), and the cascade wraps it.
// NewCloudASRChain refuses a no-endpoints-no-fallback build outright - the
// "would swallow every utterance" shape - and that refusal is the loader's
// error, which the acquire path reports with its own name.
func (r *voiceResidency) buildCloudAsrEngine() (speech.AsrEngine, error) {
	var fallback speech.AsrEngine
	if r.plan.AsrDir != "" {
		fb, ferr := speech.NewStreamingAsr(speech.AsrConfig{
			EncoderPath: filepath.Join(r.plan.AsrDir, voiceAsrEncoder),
			DecoderPath: filepath.Join(r.plan.AsrDir, voiceAsrDecoder),
			TokensPath:  filepath.Join(r.plan.AsrDir, voiceAsrTokens),
		})
		if ferr != nil {
			r.noteProblem("云端听写链的本地回退没有装上（" + ferr.Error() + "）：链尽即停，无本地兜底")
		} else {
			fallback = fb
		}
	}
	chain, err := speech.NewCloudASRChain(r.plan.CloudASREps, fallback, voiceCloudChainOptions(r.cloudRegistry))
	if err != nil {
		if fallback != nil {
			_ = fallback.Close()
		}
		return nil, err
	}
	if fallback == nil {
		return chain, nil
	}
	return &voiceCloudAsrFamily{CloudASRChain: chain, local: fallback}, nil
}

// buildCloudTtsEngine is the TTS twin: the local matcha checkpoint is the
// terminal fallback when the store carries a complete one, the cascade wraps
// it, and an incomplete checkpoint is the cloud-only chain (the same degraded
// posture the local-only loader reports, now with the cloud head in front).
func (r *voiceResidency) buildCloudTtsEngine() (speech.TtsEngine, error) {
	var fallback speech.TtsEngine
	if r.plan.Tts.complete() {
		files := r.plan.Tts
		fb, ferr := speech.NewMatchaTts(speech.TtsConfig{
			AcousticModel: files.Acoustic,
			Vocoder:       files.Vocoder,
			Lexicon:       files.Lexicon,
			Tokens:        files.Tokens,
			DictDir:       files.DictDir,
			RuleFsts:      files.RuleFsts,
		})
		if ferr != nil {
			r.noteProblem("云端朗读链的本地回退没有装上（" + ferr.Error() + "）：链尽即停，无本地兜底")
		} else {
			fallback = fb
		}
	}
	chain, err := speech.NewCloudTTSChain(r.plan.CloudTTSEps, fallback, voiceCloudChainOptions(r.cloudRegistry))
	if err != nil {
		if fallback != nil {
			_ = fallback.Close()
		}
		return nil, err
	}
	if fallback == nil {
		return chain, nil
	}
	return &voiceCloudTtsFamily{CloudTTSChain: chain, local: fallback}, nil
}
