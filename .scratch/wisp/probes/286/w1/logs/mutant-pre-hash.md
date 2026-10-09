cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189 *internal/agent/loop.go
=== sed -n '600,612p' BEFORE ===
// still goes through taskID - and appends the provider's call id, the id the
// wire, the history and the tool_call rows already key each call by. An
// id-less call falls back to its position inside the turn.
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return fmt.Sprintf("%s#call-%d", taskID, index)
	}
	return taskID + "#" + callID
}

// executeCalls runs a turn's tool calls: risk policy first (pass-through until
// ticket 21), then D38d-capped concurrent execution with the C22 per-tool
// timeout, then D15(3) spill, then the results back into the history. Tool
