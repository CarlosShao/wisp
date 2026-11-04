package panel

const (
	TreeKindNotRepo  = "not_repo"
	TreeKindNoGit    = "no_git"
	TreeKindTooLarge = "too_large"
	TreeKindCounted  = "counted"
)

type TreeEntry struct {
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type TreeView struct {
	Kind           string      `json:"kind"`
	Reason         string      `json:"reason"`
	Entries        []TreeEntry `json:"entries"`
	FilesTotal     int         `json:"filesTotal"`
	FilesTruncated bool        `json:"filesTruncated"`
	BytesTotal     int64       `json:"bytesTotal"`
	RecoverPointer string      `json:"recoverPointer"`
}

type TreeLimits struct {
	MaxEntries int   `json:"maxEntries"`
	MaxBytes   int64 `json:"maxBytes"`
	MaxDepth   int   `json:"maxDepth"`
}
