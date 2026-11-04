package panel

const (
	TreeKindNotRepo  = "not_repo"
	TreeKindNoGit    = "no_git"
	TreeKindTooLarge = "too_large"
	TreeKindCounted  = "counted"
)

type TreeEntry struct {
	Path  string
	IsDir bool
	Size  int64
}

type TreeView struct {
	Kind           string
	Reason         string
	Entries        []TreeEntry
	FilesTotal     int
	FilesTruncated bool
	BytesTotal     int64
	RecoverPointer string
}

type TreeLimits struct {
	MaxEntries int
	MaxBytes   int64
	MaxDepth   int
}
