//go:build !windows

package manifestslot

import "errors"

// errNoJobObject is reported by openChildJob where no Job Object exists. The
// spawner then falls back to killing the direct child and SAYS SO through this
// variable rather than silently claiming a tree kill (SPEC-07 §7 wants the tree
// gone; on a platform without jobs we cannot promise that).
var errNoJobObject = errors.New("manifestslot: no Job Object on this platform (process-tree kill degrades to the direct child)")
