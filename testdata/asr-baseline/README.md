# ASR baseline subset

A small, offline, byte-stable CER baseline for the speech pipeline. It ships two
acoustic conditions cut from **AISHELL-1** (OpenSLR SLR33, Apache-2.0), plus
additive noise from **THCHS-30** (OpenSLR SLR18, Apache-2.0).

The directory names are frozen by `SPEC-01:86-87` and repeated verbatim in
`SPEC-04:125`. They are cut **by acoustic condition**, not by train/dev/test
split. Both sets happen to come from the AISHELL-1 *test* split, and their
speaker sets are disjoint from each other.

```
near_clean/   60 utterances (.wav + .txt)   clean, close microphone
noisy_far/    40 utterances (.wav + .txt)   same corpus, disjoint speakers,
                                            + noise at 10 dB SNR, 6 kHz
                                            low-pass, -6 dB gain (simulated)
sources.tsv   100 data rows, one per .wav
LICENSE       attribution + the full Apache-2.0 text
tools/        the generator that produced everything above
```

Every `.wav` is 16 kHz, mono, 16-bit PCM. Every `<utt>.txt` holds one line: the
normalized reference text for that utterance.

## Readouts

| set | utts | speakers | wav bytes | MiB | duration s (min / mean / max) |
|---|---|---|---|---|---|
| `near_clean` | 60 | 8 | 7,654,144 | 7.30 | 2.987 / 3.985 / 4.927 |
| `noisy_far`  | 40 | 8 | 5,317,700 | 5.07 | 3.056 / 4.153 / 4.971 |
| **total**    | **100** | 16 | **12,971,844** | **12.37** | 405.2 s of audio |

Speakers: `near_clean` = S0764, S0765, S0766, S0767, S0768, S0769, S0770, S0901
(8/8/8/8/7/7/7/7 utts). `noisy_far` = S0902, S0903, S0904, S0905, S0906, S0907,
S0908, S0912 (5 utts each). Overlap between the two sets: none.

## `sources.tsv`

Tab separated, header row, then exactly 100 data rows in a fixed order
(`near_clean` first, then `noisy_far`). The header is frozen at five columns, in
this order, by the CER reader (`internal/speech/cer`, `manifestHeader`), which
refuses any other first line; the row order is the report order. Columns:

| column | meaning |
|---|---|
| `file` | path relative to this directory, e.g. `near_clean/BAC009S0764W0121.wav`; its first segment is the condition |
| `sha256` | sha256 of the shipped wav. For `near_clean` this equals the sha256 of the upstream AISHELL-1 file, because those files are byte-identical copies |
| `speaker` | AISHELL-1 speaker id |
| `condition` | `near_clean` or `noisy_far` |
| `origin` | provenance of the audio: `aishell-1` for `near_clean`; `aishell-1+thchs30-noise` for `noisy_far`, whose utterances are AISHELL-1 with additive THCHS-30 noise |

`file` and `condition` are what the reader uses; `sha256` is the integrity
record; `speaker` and `origin` keep the rows attributable without
re-downloading 15.6 GB.

## `sources-extra.tsv`

The derivation columns live here, same 100 rows in the same order, joined to
`sources.tsv` by `file`. They are not part of the frozen manifest, so the reader
never sees them:

| column | meaning |
|---|---|
| `file` | joins to `sources.tsv` |
| `source_utt_id` | AISHELL-1 utterance id, e.g. `BAC009S0764W0121` |
| `duration_s` | duration in seconds, 3 decimals, measured from the shipped wav |
| `snr_db` | measured in-file SNR for `noisy_far`; empty for `near_clean` |
| `noise_source` | `cafe`, `car` or `white`; empty for `near_clean` |
| `text_raw` | upstream transcript, space separated as published |

`text_norm` appears in neither file: it is byte-identical to the sibling `.txt`
for all 100 rows, so the `.txt` is the reference and there is one source of
truth.

## Reference text

The reference text is the manual transcription shipped inside AISHELL-1
(`transcript/aishell_transcript_v0.8.txt`, 141,600 lines). Nothing here was
transcribed or corrected by hand.

Normalization, applied identically to `text_norm` and to the `.txt` files:

1. upstream lines are `<utt_id><space><token><space><token>...`, where tokens are
   words and/or single characters;
2. tokens are concatenated with no separator, reproducing the original sentence;
3. selection requires that this is **lossless**: `text_norm` must equal the raw
   line with its spaces removed, and every token must be pure CJK ideographs
   (`U+4E00`-`U+9FFF`). A line that would lose a character is excluded, never
   shortened. This matters in practice: 7 of the 141,600 transcript lines contain
   full-width Latin (`ｔ恤`, `ａ级`, `ｋ歌`, `ｃ罗`, `ｂ超`) and are excluded by
   this rule. Among the 100 shipped rows, `text_raw` contains zero non-CJK
   characters, so nothing was dropped;
4. sentences containing `<UNK>` are excluded from selection entirely, as are
   sentences shorter than 8 characters.

There is no ITN and no traditional/simplified conversion. A character-level CER
consumer should therefore compare against `text_norm` unchanged.

## How the two sets were selected

Eligibility, applied to all 7,176 utterances of the AISHELL-1 test split:

* duration within `[2.0, 5.0]` seconds, measured from the wav header;
* reference text passes the normalization rules above.

Then, on speakers sorted ascending: the first 8 eligible speakers feed
`near_clean` (60 utterances), the next 8 feed `noisy_far` (40). Within a set,
utterances are taken round-robin across speakers in ascending `utt_id` order,
which keeps per-speaker counts flat (8/8/8/8/7/7/7/7 and 5 each).

Selection never looks at the audio content, only at duration and text form.

## `near_clean` is a byte-for-byte copy

AISHELL-1 is natively 16 kHz / mono / 16-bit PCM, so "convert to 16 kHz mono
16-bit" is a no-op here: the upstream wav is copied byte for byte. No
transcoding, no resampling, no normalisation, no trimming. That is why each
`near_clean` sha256 in `sources.tsv` also identifies the upstream file.

## `noisy_far` is a simulation, not a far-field recording

Per utterance, deterministically:

1. take the clean utterance `s` and a 300 s THCHS-30 noise excerpt `n` of equal
   length, from `cafe` / `car` / `white` cycled in that fixed order, at a
   per-file random offset;
2. scale `n` so that `10*log10(P(s)/P(n))` equals exactly 10 dB;
3. `m = s + n`;
4. `m` is low-passed by a 129-tap Blackman-windowed-sinc FIR at 6 kHz;
5. the result is scaled by -6 dB (`0.501187`) and quantized to 16-bit PCM.

Measured in-file SNR (computed from the shipped wav against the upstream clean
signal, i.e. after the low-pass, not merely the mixing target):

| noise | n | min dB | mean dB | max dB |
|---|---|---|---|---|
| `cafe` | 14 | 9.964 | 9.994 | 10.011 |
| `car` | 13 | 9.952 | 9.989 | 9.998 |
| `white` | 13 | 11.114 | 11.271 | 11.299 |
| all | 40 | 9.952 | 10.408 | 11.299 |

The mixing target is exactly 10.000 dB for every file. `cafe` and `car` stay at
10 dB through the filter because their energy is already concentrated below
6 kHz; `white` rises to about 11.3 dB because broadband noise loses more energy
to the 6 kHz low-pass than speech does. This is a property of the filter, not a
mixing error.

This set must not be described as real far-field audio.

## Determinism

`SEED = 20261026`. Every choice is a pure function of the seed and of sorted
upstream identifiers: no wall clock, no unseeded RNG, no dict-order dependence.
Regenerating from the same upstream archives reproduces `sources.tsv` with an
identical sha256 for every one of the 100 files (verified).

## Regenerating

Upstream archives, with the sizes and sha256 that were actually used. OpenSLR
publishes no checksums, so these are pinned here for the first time:

| archive | bytes | sha256 |
|---|---|---|
| `https://openslr.org/resources/33/data_aishell.tgz` | 15,582,913,665 | `a4a0313cde0a933e0e01a451f77de0a23d6c942f4694af5bb7f40b9dc38143fe` |
| `https://openslr.org/resources/18/resource.tgz` | 24,813,708 | `5f10b11a86930d159250c979b2dccb26ac2527775b2aa628df08513af32d85ee` |

Noise excerpts used, from `resource.tgz`:

| file | bytes | sha256 |
|---|---|---|
| `resource/noise/cafe/cafe.wav` | 9,600,172 | `1fd63bab8ea6fc443bbe29df61dd8325780720d55b5172b1ae6e17809c767310` |
| `resource/noise/car/car.wav` | 9,600,172 | `ffe4cac12ed78705a52fce5d14f17a4ee838c6f5922446fc66258841323d4a59` |
| `resource/noise/white/white.wav` | 9,600,044 | `6ab0168081663e0d37119431f082965fad2b4a2a6bdfe9f5fe8b84602d9467a7` |

`data_aishell.tgz` holds one `wav/S####.tar.gz` per speaker; the train/dev/test
split is encoded **inside** each of those, as the first path component
(`test/S0901/...`). There is no `wav/test/` directory in the archive. The 20 test
speakers are S0764-S0770, S0901-S0908, S0912-S0916, and they contain exactly
7,176 utterances.

```bash
# 1. fetch (do this OUTSIDE the repo; the archives are 15.6 GB + 24 MB)
TMP=/d/tmp/wisp/20261026/15d2
mkdir -p "$TMP" && cd "$TMP"
curl -L --fail -o data_aishell.tgz https://openslr.org/resources/33/data_aishell.tgz
curl -L --fail -o resource.tgz     https://openslr.org/resources/18/resource.tgz

# 2. unpack just what is needed
tar -xzf data_aishell.tgz -C "$TMP/stage"
mkdir -p "$TMP/aishell/transcript" "$TMP/aishell/wav"
cp "$TMP/stage/data_aishell/transcript/aishell_transcript_v0.8.txt" "$TMP/aishell/transcript/"
for s in S0764 S0765 S0766 S0767 S0768 S0769 S0770 S0901 S0902 S0903 S0904 \
         S0905 S0906 S0907 S0908 S0912 S0913 S0914 S0915 S0916; do
  tar -xzf "$TMP/stage/data_aishell/wav/$s.tar.gz" -C "$TMP/aishell/wav"
done   # inner members are "test/S####/*.wav", so this yields wav/test/S####/*.wav
tar -xzf resource.tgz -C "$TMP"          # -> $TMP/resource/noise/{cafe,car,white}/*.wav

# 3. generate (needs python3 + numpy; writes the wavs, the .txt files,
#    sources.tsv and sources-extra.tsv)
cd /d/work/workspace/"projects plans"/Wisp
python3 testdata/asr-baseline/tools/make_asr_baseline.py \
  --aishell "$TMP/aishell" \
  --noise   "$TMP/resource/noise" \
  --out     testdata/asr-baseline
```

Step 3 writes both tables: `sources.tsv` (5 columns, the frozen manifest the CER
reader parses) and `sources-extra.tsv` (6 columns, the derivation columns that
reader never sees). Both are written with a fixed column order and a fixed row
order, so re-running step 3 over an already-generated tree rewrites identical
bytes: `sources.tsv` stays at sha256
`54b50a499e5b33c0e722f4a1995e02023f3cc25daec6aefa48dc6089f64db537` and
`sources-extra.tsv` at
`531677f225b7cc24f5530505fc058458a1bccd356eb96ef028bc27f91b2686dd`. The
reference text is not a column in either file: it is the sibling `<utt>.txt`.

`--n-near` / `--n-far` override the counts, for trial runs only.

## Cost model

16 kHz x 16-bit x 1 channel = 32,000 B/s. The 100 files total 12,971,844 bytes
for 405.23 s of audio, which is 32,011 B/s including the 44-byte headers: the
accounting is exact, there is no hidden overhead.

The research note that preceded this set
(`docs/evidence/s1/cer-a1-research.md`, section 5) budgeted about 9.2 MiB on the
assumption of 3.0 s per utterance. Real AISHELL-1 test utterances average
5.0 s (1,155,801,150 bytes / 7,176 files), and inside the mandated 2-5 s window
they still average 4.0 s, so a 60+40 cut lands at 12.37 MiB rather than 9.2 MiB.
The byte accounting is not the problem; the 3.0 s/utterance guess was.

## Licence

`LICENSE` in this directory carries the attribution for both upstream corpora
and the full text of the Apache License, Version 2.0. Both corpora are labelled
`License: Apache License v.2.0` by OpenSLR.

One open point, recorded rather than resolved: the SLR33 page also says, verbatim,
"The data is free for academic use.", which is inconsistent with its own
Apache-2.0 label. It was raised in `cer-a1-research.md` section 6 and is still
awaiting a decision on whether to register it as a deferred risk.
