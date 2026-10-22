//go:build !windows

package models

// POSIX has no alternate data streams, so this half of ticket 357's enumeration
// answers the only thing the report needs from it: nothing.
//
// It is NOT evidence that the Windows half works. On this platform the call can
// only ever return an empty list, so the case next to it (ads_other_test.go) is a
// control leg - it keeps the build and the caller honest on the leg CI runs
// without a DLL, and it proves the report stays silent where there is nothing to
// report. The tooth for ticket 357 lives on NTFS and cannot be moved here.
func listNamedStreams(string) ([]string, error) {
	return nil, nil
}
