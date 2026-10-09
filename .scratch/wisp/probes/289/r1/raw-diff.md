diff --git a/cmd/wisp/subagent_carrier_197_test.go b/cmd/wisp/subagent_carrier_197_test.go
index 370cda0b..b0ded836 100644
--- a/cmd/wisp/subagent_carrier_197_test.go
+++ b/cmd/wisp/subagent_carrier_197_test.go
@@ -20,4 +20,4 @@ package main
-// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,
-// internal/agent/loop.go:647, which is also what makes the roster's parent link read
-// as a task id). The spawn therefore waits for the ROOT's own roster row - proof the
-// root loop is live - and only then dispatches through rt.bridge.
+// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
+// NOT the loop's: since bd124b2a callCorr mints one corr per call
+// (<taskID>#<callID>, internal/agent/loop.go:603, used :676). The spawn waits
+// for the ROOT's own roster row - proof it is live - then uses rt.bridge.
diff --git a/internal/panel/subagent_blocked_220_test.go b/internal/panel/subagent_blocked_220_test.go
index a6b552e5..dbdce69d 100644
--- a/internal/panel/subagent_blocked_220_test.go
+++ b/internal/panel/subagent_blocked_220_test.go
@@ -155,2 +155,2 @@ func TestAL1WindowInWaitingLightsItsRow(t *testing.T) {
-	// The production pairing (the loop dispatches with CorrelationID == TaskID,
-	// internal/agent/loop.go:647) lights through the correlation id alone.
+	// Since bd124b2a the loop's corr is per call (<taskID>#<callID>, loop.go:603,
+	// used :676), never == taskID: the corr-only pairing below is this fixture's.
diff --git a/internal/panel/subagent_roster_197.go b/internal/panel/subagent_roster_197.go
index 341878c8..35ce6ffb 100644
--- a/internal/panel/subagent_roster_197.go
+++ b/internal/panel/subagent_roster_197.go
@@ -56,3 +56,3 @@ type TaskRow struct {
-	// The entity leg files the deriving call's correlation id, and the agent
-	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),
-	// so on the production path this reads as the parent's task id.
+	// The entity leg files the deriving call's task id (TaskID(ctx),
+	// internal/tools/subagent_197.go:262), never its correlation id: since
+	// bd124b2a loop.go's callCorr mints one per call (<taskID>#<callID>, :603/:676).
diff --git a/internal/panel/subagent_roster_197_test.go b/internal/panel/subagent_roster_197_test.go
index f803179d..1808b13d 100644
--- a/internal/panel/subagent_roster_197_test.go
+++ b/internal/panel/subagent_roster_197_test.go
@@ -466,5 +466,5 @@ func TestDroppedStreamsAreNamedOnTheWire(t *testing.T) {
-//  3. The blockedOnApproval join is asserted against a card this file made. On the
-//     production path the pairing depends on the loop dispatching with
-//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a
-//     different pairing would report every row as not blocked, and the only thing
-//     that would catch it is the cmd/wisp case, which uses the real gate.
+//  3. The blockedOnApproval join is asserted against a card this file made. On
+//     the production path it pairs the row's task id, not the corr: since
+//     bd124b2a callCorr mints one corr per call (loop.go:603, used :676),
+//     shape <taskID>#<callID>, never == taskID. A host that paired on the
+//     corr alone would miss every row; only cmd/wisp's real gate catches it.
