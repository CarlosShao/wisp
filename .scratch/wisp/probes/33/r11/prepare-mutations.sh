#!/bin/sh
# 33-r11 mutation preparation: every mutated copy is written under this probe
# directory and reaches the compiler only through -overlay, so no tracked file is
# ever edited for a mutation and there is nothing to revert.
cd "D:/work/workspace/projects plans/Wisp" || exit 9
R=.scratch/wisp/probes/33/r11/mutations
ABS="D:/work/workspace/projects plans/Wisp"
HOST=cmd/wisp/panel_host_windows.go

# M1: put the old name back on the self-locking round-trip method (and its call site,
# so it still compiles) - the shape this leg renamed away.
sed -e 's|^func (m \*PanelManager) firstRoundTrip(|func (m *PanelManager) firstRoundTripLocked(|' \
    -e 's|m\.firstRoundTrip(ctx, t0)|m.firstRoundTripLocked(ctx, t0)|' \
    "$HOST" > "$R/mut-m1-old-name-back.go"

# M2: a brand new Locked-suffixed method that does NOT lock - only the roster can
# catch this one, so it proves the exemption list has its own tooth.
cat "$HOST" > "$R/mut-m2-unrostered-locked-name.go"
cat >> "$R/mut-m2-unrostered-locked-name.go" <<'EOF'

// noteColdMsLocked is the M2 mutation: a Locked-suffixed method that takes no lock
// of its own, i.e. a name the roster has never heard of.
func (m *PanelManager) noteColdMsLocked(v float64) { m.lastColdMs = v }
EOF

# M3: make the rostered, honest method take the lock itself - the census and the
# behavioural case must both go red on this one.
sed -e 's|^func (m \*PanelManager) setPriorFocusLocked(prior windows.HWND) {|func (m *PanelManager) setPriorFocusLocked(prior windows.HWND) {\n\tm.mu.Lock()\n\tdefer m.mu.Unlock()|' \
    "$HOST" > "$R/mut-m3-rostered-method-self-locks.go"

# CONTROL: rename only the declaration, leave the call site pointing at the old name.
# Compilation must fail, which is what proves -overlay really feeds this compiler and
# that a green run is a green run of the mutated tree.
sed -e 's|^func (m \*PanelManager) firstRoundTrip(|func (m *PanelManager) firstRoundTripMovedAwayControl(|' \
    "$HOST" > "$R/mut-control-symbol-moved-outside-the-overlay-copy.go"

i=0
for f in "$R"/mut-*.go; do
  i=$((i+1))
  base=$(basename "$f")
  name=${base#mut-}; name=${name%%-*}
  printf '{\n  "Replace": {\n    "%s/%s": "%s/%s/%s"\n  }\n}\n' "$ABS" "$HOST" "$ABS" "$R" "$base" > "$R/overlay-$name.json"
done
echo "prepare_rc=$?"
ls -1 "$R"
echo "listing_rc=$?"
echo "--- what each mutated copy differs by, against the current worktree file ---"
for f in "$R"/mut-m*.go "$R"/mut-c*.go; do
  echo "## $(basename "$f")"
  diff "$HOST" "$f" | head -8
done
echo "diff_rc=$?"
