#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""make_asr_baseline.py -- regenerate testdata/asr-baseline/ deterministically.

Upstream (both Apache-2.0, both fetched from OpenSLR):
  * AISHELL-1  https://openslr.org/resources/33/data_aishell.tgz
               15,582,913,665 B  ->  transcript/aishell_transcript_v0.8.txt  +  wav/test/**
  * THCHS-30   https://openslr.org/resources/18/resource.tgz
               24,813,708 B      ->  resource/noise/{cafe,car,white}/*.wav

What it builds
  near_clean/  60 utterances, AISHELL-1 *test* split, byte-identical copies of the
               upstream wav (upstream is already 16 kHz / mono / 16-bit PCM, so
               "convert" is a no-op and each file's sha256 equals the upstream one).
  noisy_far/   40 utterances, AISHELL-1 *test* split, speakers DISJOINT from
               near_clean; real THCHS-30 noise mixed in at a target 10 dB SNR,
               then a 6 kHz FIR low-pass and a -6 dB gain to fake far-field.
               This is a SIMULATION, not a real far-field recording.

Determinism: no RNG state leaks between files, no wall-clock, no dict-iteration
order. Every choice is a pure function of SEED plus sorted upstream identifiers.
Only numpy is required (stdlib `wave`/`hashlib` do the rest).

Usage:
  python make_asr_baseline.py --aishell <dir with transcript/ and wav/> \
                              --noise <dir with cafe/ car/ white/> \
                              --out <testdata/asr-baseline>
"""

from __future__ import annotations

import argparse
import hashlib
import os
import random
import sys
import wave

import numpy as np

# ---------------------------------------------------------------- parameters
SEED = 20261026

N_NEAR = 60          # near_clean utterance count
N_FAR = 40           # noisy_far  utterance count
N_SPEAKERS = 8       # speakers used per condition; the two sets are disjoint

DUR_MIN = 2.0        # seconds, inclusive
DUR_MAX = 5.0        # seconds, inclusive
MIN_CHARS = 8        # drop very short transcripts (unstable CER)

TARGET_SNR_DB = 10.0  # mixing target for noisy_far
LOWPASS_HZ = 6000.0   # "far-field" low-pass corner
FAR_GAIN = 0.501187   # -6 dB overall gain ("far-field" level drop)
FIR_TAPS = 129        # odd length, symmetric

NOISES = ("cafe", "car", "white")  # cycled in this fixed order

SENTINEL = "<UNK>"


# ------------------------------------------------------------------- helpers
def cjk_only(text: str) -> bool:
    """True iff every char is a BMP CJK ideograph (no ASCII, no punctuation)."""
    if not text:
        return False
    for ch in text:
        if not ("\u4e00" <= ch <= "\u9fff"):
            return False
    return True


def normalize(raw: str) -> str:
    """AISHELL raw is space-separated tokens; the reference text is space-free."""
    return "".join(raw.split())


def lossless(raw: str) -> bool:
    """True iff normalisation keeps every character and every token is pure CJK.

    A handful of transcript lines contain full-width Latin (e.g. a token built on
    U+FF41). Those must be skipped, not silently shortened -- a reference text
    that quietly loses a character would inflate CER against every hypothesis.
    """
    tokens = raw.split()
    if not tokens:
        return False
    if not all(cjk_only(t) for t in tokens):
        return False
    return "".join(tokens) == normalize(raw)


def read_wav(path: str) -> tuple[np.ndarray, int]:
    with wave.open(path, "rb") as w:
        if w.getnchannels() != 1 or w.getsampwidth() != 2:
            raise SystemExit(f"{path}: expected mono 16-bit PCM, got "
                             f"{w.getnchannels()}ch/{w.getsampwidth() * 8}bit")
        sr = w.getframerate()
        frames = w.readframes(w.getnframes())
    return np.frombuffer(frames, dtype="<i2"), sr


def write_wav(path: str, samples: np.ndarray, sr: int) -> None:
    if samples.dtype != np.int16:
        raise SystemExit("write_wav expects int16")
    with wave.open(path, "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(sr)
        w.writeframes(samples.astype("<i2").tobytes())


def sha256_file(path: str) -> str:
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def lowpass_kernel(sr: int, cutoff_hz: float, taps: int) -> np.ndarray:
    """Blackman-windowed sinc, unity DC gain. taps must be odd."""
    if taps % 2 == 0:
        raise SystemExit("taps must be odd")
    m = (taps - 1) / 2.0
    n = np.arange(taps) - m
    fc = cutoff_hz / sr  # normalised to sample rate
    h = 2.0 * fc * np.sinc(2.0 * fc * n)
    h *= np.blackman(taps)
    return h / h.sum()


def measure_snr_db(signal: np.ndarray, noise: np.ndarray) -> float:
    ps = float(np.mean(signal.astype(np.float64) ** 2))
    pn = float(np.mean(noise.astype(np.float64) ** 2))
    if ps <= 0.0 or pn <= 0.0:
        raise SystemExit("degenerate power; refusing to report an SNR")
    return 10.0 * np.log10(ps / pn)


# --------------------------------------------------------------- transcript
def load_transcript(path: str) -> dict[str, str]:
    table: dict[str, str] = {}
    with open(path, "r", encoding="utf-8") as f:
        for lineno, line in enumerate(f, 1):
            line = line.strip()
            if not line:
                continue
            parts = line.split(None, 1)  # utt id, then the char sequence
            if len(parts) != 2:
                raise SystemExit(f"{path}:{lineno}: cannot split id from text: {line!r}")
            table[parts[0]] = parts[1]
    if not table:
        raise SystemExit(f"{path}: transcript parsed to zero rows")
    return table


def candidate_pool(wav_root: str, transcript: dict[str, str]) -> dict[str, list[tuple[str, float]]]:
    """speaker -> [(utt_id, duration_s)] for every eligible utterance, sorted."""
    pool: dict[str, list[tuple[str, float]]] = {}
    for utt, raw in sorted(transcript.items()):
        if SENTINEL in raw:
            continue
        if not lossless(raw):
            continue
        if len(normalize(raw)) < MIN_CHARS:
            continue
        speaker = utt[6:11]  # BAC009S0002W0122 -> S0002
        path = os.path.join(wav_root, speaker, utt + ".wav")
        if not os.path.isfile(path):
            continue
        nframes, sr = _nframes(path)
        if sr != 16000:
            raise SystemExit(f"{path}: expected 16 kHz, got {sr}")
        dur = nframes / float(sr)
        if not (DUR_MIN <= dur <= DUR_MAX):
            continue
        pool.setdefault(speaker, []).append((utt, dur))
    for speaker in pool:
        pool[speaker].sort()
    return pool


def _nframes(path: str) -> tuple[int, int]:
    with wave.open(path, "rb") as w:
        return w.getnframes(), w.getframerate()


def take(pool: dict[str, list[tuple[str, float]]], speakers: list[str], total: int) -> list[tuple[str, str, float]]:
    """Round-robin over speakers, so both spoken-length and speaker coverage are flat."""
    picked: list[tuple[str, str, float]] = []
    cursors = {s: 0 for s in speakers}
    while len(picked) < total:
        progressed = False
        for s in speakers:
            if len(picked) >= total:
                break
            idx = cursors[s]
            if idx < len(pool[s]):
                utt, dur = pool[s][idx]
                picked.append((utt, s, dur))
                cursors[s] = idx + 1
                progressed = True
        if not progressed:
            raise SystemExit(f"not enough eligible utterances to reach {total}")
    return picked


# --------------------------------------------------------------------- main
def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--aishell", required=True, help="dir containing transcript/ and wav/")
    ap.add_argument("--transcript", default="transcript/aishell_transcript_v0.8.txt")
    ap.add_argument("--noise", required=True, help="dir containing cafe/ car/ white/")
    ap.add_argument("--out", required=True)
    ap.add_argument("--n-near", type=int, default=N_NEAR, help="override near_clean count (trial runs)")
    ap.add_argument("--n-far", type=int, default=N_FAR, help="override noisy_far count (trial runs)")
    args = ap.parse_args()

    wav_root = os.path.join(args.aishell, "wav", "test")
    if not os.path.isdir(wav_root):
        raise SystemExit(f"missing test split wav dir: {wav_root}")
    transcript = load_transcript(os.path.join(args.aishell, args.transcript))
    print(f"transcript rows: {len(transcript)}")

    pool = candidate_pool(wav_root, transcript)
    speakers = sorted(s for s in pool if len(pool[s]) >= 5)
    print(f"test speakers with >=5 eligible utts: {len(speakers)} -> {speakers}")
    if len(speakers) < 2 * N_SPEAKERS:
        raise SystemExit(f"need {2 * N_SPEAKERS} usable test speakers, found {len(speakers)}")

    near_spk = speakers[:N_SPEAKERS]
    far_spk = speakers[N_SPEAKERS:2 * N_SPEAKERS]
    assert not set(near_spk) & set(far_spk)

    near = take(pool, near_spk, args.n_near)
    far = take(pool, far_spk, args.n_far)
    print(f"near_clean: {len(near)} utts over {len(near_spk)} speakers {near_spk}")
    print(f"noisy_far : {len(far)} utts over {len(far_spk)} speakers {far_spk}")

    for cond in ("near_clean", "noisy_far"):
        os.makedirs(os.path.join(args.out, cond), exist_ok=True)

    fir = lowpass_kernel(16000, LOWPASS_HZ, FIR_TAPS)
    noise_cache: dict[str, np.ndarray] = {}
    rows: list[dict] = []

    # ---- near_clean: byte-identical copy, no transcoding, no resampling
    for utt, speaker, dur in near:
        src = os.path.join(wav_root, speaker, utt + ".wav")
        dst = os.path.join(args.out, "near_clean", utt + ".wav")
        with open(src, "rb") as fi, open(dst, "wb") as fo:
            fo.write(fi.read())
        nframes, sr = _nframes(dst)
        assert sr == 16000
        text = normalize(transcript[utt])
        with open(os.path.join(args.out, "near_clean", utt + ".txt"), "w",
                  encoding="utf-8", newline="\n") as f:
            f.write(text + "\n")
        rows.append(dict(cond="near_clean", utt=utt, speaker=speaker,
                         dur=nframes / 16000.0, sha=sha256_file(dst),
                         snr="", noise="", text=text, file=f"near_clean/{utt}.wav"))

    # ---- noisy_far: deterministic mix
    for i, (utt, speaker, dur) in enumerate(far):
        src = os.path.join(wav_root, speaker, utt + ".wav")
        clean_i16, sr = read_wav(src)
        assert sr == 16000
        clean = clean_i16.astype(np.float64) / 32768.0

        name = NOISES[i % len(NOISES)]
        if name not in noise_cache:
            noise_path = os.path.join(args.noise, name, name + ".wav")
            noise_i16, nsr = read_wav(noise_path)
            if nsr != 16000:
                raise SystemExit(f"{noise_path}: expected 16 kHz, got {nsr}")
            noise_cache[name] = noise_i16.astype(np.float64) / 32768.0
        pool_noise = noise_cache[name]
        if len(pool_noise) < len(clean):
            raise SystemExit(f"noise {name} shorter than utterance {utt}")

        rng = random.Random(SEED + i)  # per-file, order-independent
        off = rng.randrange(0, len(pool_noise) - len(clean) + 1)
        noise = pool_noise[off:off + len(clean)].copy()

        ps = float(np.mean(clean ** 2))
        pn = float(np.mean(noise ** 2))
        gain = np.sqrt((ps / (10.0 ** (TARGET_SNR_DB / 10.0))) / pn)
        noise *= gain
        snr_mix = measure_snr_db(clean, noise)  # == TARGET_SNR_DB by construction

        mix = clean + noise
        mix_lp = np.convolve(mix, fir, mode="same")
        # linear + time-invariant => lp(s+n) == lp(s) + lp(n); measure both sides
        # through the identical chain to report the SNR actually present in-file
        s_lp = np.convolve(clean, fir, mode="same")
        n_lp = np.convolve(noise, fir, mode="same")
        snr_chain = measure_snr_db(s_lp, n_lp)

        out = mix_lp * FAR_GAIN
        peak = float(np.max(np.abs(out)))
        if peak >= 1.0:
            raise SystemExit(f"{utt}: clipping (peak {peak:.4f})")
        out_i16 = np.clip(np.round(out * 32768.0), -32768, 32767).astype(np.int16)

        dst = os.path.join(args.out, "noisy_far", utt + ".wav")
        write_wav(dst, out_i16, 16000)
        text = normalize(transcript[utt])
        with open(os.path.join(args.out, "noisy_far", utt + ".txt"), "w",
                  encoding="utf-8", newline="\n") as f:
            f.write(text + "\n")
        rows.append(dict(cond="noisy_far", utt=utt, speaker=speaker,
                         dur=len(out_i16) / 16000.0, sha=sha256_file(dst),
                         snr=f"{snr_chain:.4f}", noise=name, text=text,
                         file=f"noisy_far/{utt}.wav"))
        print(f"  {utt} {name:5s} off={off:7d} snr_mix={snr_mix:6.2f} "
              f"snr_chain={snr_chain:6.2f} peak={peak:.3f}")

    # ---- sources.tsv
    cols = ["file", "sha256", "source_utt_id", "duration_s",
            "speaker", "condition", "snr_db", "noise_source", "text_raw", "text_norm"]
    tsv = "\t".join(cols) + "\n"
    for r in rows:
        tsv += "\t".join([
            r["file"], r["sha"], r["utt"], f"{r['dur']:.3f}",
            r["speaker"], r["cond"], r["snr"], r["noise"],
            " ".join(r["text"]), r["text"],
        ]) + "\n"
    with open(os.path.join(args.out, "sources.tsv"), "w", encoding="utf-8", newline="\n") as f:
        f.write(tsv)

    total_bytes = 0
    for r in rows:
        total_bytes += os.path.getsize(os.path.join(args.out, r["file"]))
    for cond in ("near_clean", "noisy_far"):
        sub = [r for r in rows if r["cond"] == cond]
        nb = sum(os.path.getsize(os.path.join(args.out, r["file"])) for r in sub)
        secs = sum(r["dur"] for r in sub)
        print(f"\n{cond}: n={len(sub)} speakers={len({r['speaker'] for r in sub})} "
              f"audio={secs:.2f}s bytes={nb} ({nb / 1048576:.2f} MiB) "
              f"wav_avg={nb / len(sub):.0f}B")
        if cond == "noisy_far":
            snrs = [float(r["snr"]) for r in sub]
            print(f"  snr_chain_db: min={min(snrs):.3f} mean={np.mean(snrs):.3f} max={max(snrs):.3f}")
        durs = [r["dur"] for r in sub]
        print(f"  dur_s: min={min(durs):.3f} mean={np.mean(durs):.3f} max={max(durs):.3f}")
    print(f"\nTOTAL wav bytes={total_bytes} = {total_bytes / 1048576:.2f} MiB "
          f"= {total_bytes / 1000000:.2f} MB")
    print(f"rows in sources.tsv: {len(rows)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
