package skipfixture

import "testing"

func TestSkipsOnPurpose(t *testing.T) {
	t.Skip("AC#3 fixture: runtests.sh must treat SKIP as fatal, not as a pass")
}
