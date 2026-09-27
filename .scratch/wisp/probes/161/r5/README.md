# probes/161/r5 - 161-r5's working directory

Three cells, one instrument. Read-only on everything it does not own.

| file | what it is |
|---|---|
| `attrib.sh` | r5's copy of `probes/161/r4/attrib.sh`, same two rulers, same rules, plus `--tracked-only` (what CI runs), `--classify-line <text>` (the classifier, drivable from a line of text) and `--self-test` (that classifier, checked in both directions). `r4/attrib.sh` was not edited. |
| `tracked-dirty-proof.sh` | cell 3: proves form (A)'s `TRACKED-DIRTY` leg and its zero-file refusal in an out-of-repo copy of the tracked set. |
| `negative-control/bad-sample.go` | byte-identical copy of `probes/999/bad-sample.go` (`sha256 8aa59c87a088247037141081cd308ad22d6fb4e1bc6c341f5a72b18eb69fc3c4` for both), under a path the attribution rule CAN resolve to a real ticket. |
| `logs/` | every reading this program pasted, in run order. |

## Why `negative-control/` exists, and why BOTH copies are untracked

The r4 program proved "the ruler can ring" by leaving a genuinely unattributable
`.go` file on disk at `probes/999/bad-sample.go`. That made form (B) red on every
run on this tree - which is the same defect as a permanently green ruler, pointed the
other way. 161-r5's ruling: the ring moves into `--self-test`, where it is
reproducible without a wound.

This repo's create-only rule means the old file stays where it is, so the copy is
**not a move in the filesystem sense**: `probes/999/bad-sample.go` still exists and
still forces form (B) to rc=1 until its owner retires it. Nothing below the owner may
delete it, and this program ran no delete commands at all.

Neither copy is committed, and neither may be: `probes/999/bad-sample.go` and
`probes/161/r5/negative-control/bad-sample.go` are both deliberately non-gofumpt
bytes, and a *tracked* file in that shape is exactly what form (A) - the leg CI
executes, see the `gofmt (gofumpt) - the tracked set is the denominator` step in
`.github/workflows/ci.yml` - goes red on. An untracked bad sample is a bench
artifact; a committed bad sample is a build break. `git ls-files` lists neither.
