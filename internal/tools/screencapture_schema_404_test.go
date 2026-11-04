package tools

import (
	"encoding/json"
	"testing"
)

// Test404ScreenCaptureSchemaIsValidJSON pins ticket 404: screenCaptureSchema's
// "target" field was concatenated straight onto windowTargetProperties with no
// separating comma, so the assembled schema was invalid JSON
// (`..."region"}"hwnd":{...}`) and marshalling any tool list that includes the
// screen-capture tool failed with `invalid character '"' after object key:value
// pair`. This is the pre-existing defect ticket 378's R4 surfaced: it blocks
// every real "wisp run 跑一发" test (e.g. Test233HotTierAgentReachesTheNextTask).
func Test404ScreenCaptureSchemaIsValidJSON(t *testing.T) {
	if !json.Valid([]byte(screenCaptureSchema)) {
		t.Fatalf("screenCaptureSchema is not valid JSON: %s", string(screenCaptureSchema))
	}
}
