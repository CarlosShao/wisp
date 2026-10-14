// Ticket 178 AC#2 mutation sample. NOT committed to the repo, never written into the
// working tree: it only exists inside the synthetic tree this leg builds from HEAD +
// this blob (objects under a scratch GIT_OBJECT_DIRECTORY in $TMPDIR). It is the shape
// AC#2 requires - a REAL open-only production call site (not a test): the pair census
// must still ring on it after the four-pile fix, otherwise the fix washed the
// "新增未成对" flavour away and ticket 178 fails.

package tools

// zz178r1OpenOnly opens one task's C25 taint scope; no call in this file, and no
// other file, ever closes it (the defect shape G6/G6neg must name).
func zz178r1OpenOnly(b *Bridge, taskID string) {
	b.OpenTask(taskID)
}
