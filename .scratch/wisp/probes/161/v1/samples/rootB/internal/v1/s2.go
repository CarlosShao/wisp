package v1

import fp "path/filepath"

func norm(p string) string {
	return fp.Clean(p)
}
