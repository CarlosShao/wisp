//go:build windows

package models

// Ticket 357's second face, the read-only half: enumerate the named alternate
// data streams (ADS) standing on one installed file, on the only platform that
// has them.
//
// WHY THIS FILE EXISTS AND WHY IT IS THIS SMALL. Ticket 357's ruling (the form
// the ticket calls 丁) is a report, not a refusal: `wisp models verify` must say
// one sentence when a directory it has just called consistent also holds a
// second data stream the signed manifest never named - the reading 357-a1 took
// (D1: a 12-byte model.onnx plus a 19-byte model.onnx:attacker) printed
// 与已验签清单一致 with rc=0 and said nothing about the stream. So this file adds
// one enumeration and no policy: it decides nothing about what "unverified"
// means, it does not touch the hash face, and no download or hand-off path runs
// it.
//
// THE MECHANISM, AND WHAT IT DELIBERATELY IS NOT. GetFileInformationByHandleEx
// with FileStreamInfo is already wrapped by golang.org/x/sys/windows
// (zsyscall_windows.go:2432 for the call, types_windows.go:2765 for the class
// constant), so the only thing written by hand here is the 24-byte
// FILE_STREAM_INFO header - a data layout, not a syscall. That is this
// repository's existing shape: internal/ball/win32_windows.go:160 onward mirrors
// a dozen Win32 structs and sizes them with unsafe.Sizeof. FindFirstStreamW and
// FindNextStreamW, which the dispatch first named, are NOT exported by x/sys
// v0.48.0 at all (measured: 0 hits across its 27 files, with CreateFile,
// FindClose, GetFileInformationByHandleEx and FileStreamInfo as live controls),
// so reaching for them would have meant a hand-written syscall for an answer an
// exported wrapper already gives.
//
// THE TWO MEASURED BOUNDARIES THIS SHAPE COMES FROM:
//  1. one 4096-byte buffer is not enough once a file carries many streams: the
//     call answers ERROR_MORE_DATA and an EMPTY record set. Reading that as "no
//     streams" would be a false negative on exactly the file worth reporting, so
//     the buffer grows and the call is retried - with a ceiling, and with a named
//     error past it, never with a silent empty answer.
//  2. a directory handle answers with an empty list and no error, so this can
//     only be asked one file at a time; UnnamedStreams walks the manifest's
//     installed files rather than handing a directory to this function.

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// fileStreamInfo mirrors FILE_STREAM_INFO as the OS lays it out on amd64:
// 4 + 4 + 8 + 8 bytes, then the stream name as UTF-16 immediately after. This is
// a data layout, not a syscall - the same kind of hand-mirrored struct
// internal/ball/win32_windows.go carries.
type fileStreamInfo struct {
	NextEntryOffset      uint32
	StreamNameLength     uint32
	StreamSize           int64
	StreamAllocationSize int64
}

const (
	// streamInfoHeaderBytes is the fixed part of one record; the name follows it.
	streamInfoHeaderBytes = int(unsafe.Sizeof(fileStreamInfo{}))

	// streamStartBuf and streamCeiling bound the ERROR_MORE_DATA retry. 4096 is
	// the measured breaking point (200 short streams overflow it and come back as
	// an empty list with ERROR_MORE_DATA); 16 MiB is 4096 times that, far past any
	// directory a model install could hold, and still a ceiling, so a corrupt
	// length field cannot make this allocate without a bound.
	streamStartBuf = 4 << 10
	streamCeiling  = 16 << 20
)

// listNamedStreams returns the named streams of one file as the OS spells them
// (":name:$DATA"), excluding the file's own unnamed "::$DATA" stream. An empty
// result means "this file carries no second stream", and it is only ever
// returned by a call that actually completed.
func listNamedStreams(path string) ([]string, error) {
	return listNamedStreamsBounded(path, streamStartBuf, streamCeiling)
}

// listNamedStreamsBounded is the same enumeration with the retry bounds as
// arguments, so the ceiling's behaviour ("name the failure, never answer empty")
// has a case that reaches it without planting 16 MiB of stream names.
func listNamedStreamsBounded(path string, start, ceiling int) ([]string, error) {
	if start <= 0 {
		start = streamStartBuf
	}
	if ceiling < start {
		ceiling = start
	}
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	// FILE_READ_ATTRIBUTES is the least right that answers this question, and the
	// sharing mask keeps a file another process has open openable: this is a
	// report, so it must never be the reason a check fails.
	h, err := windows.CreateFile(p,
		windows.FILE_READ_ATTRIBUTES,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(h)

	size := start
	for {
		buf := make([]byte, size)
		callErr := windows.GetFileInformationByHandleEx(h, windows.FileStreamInfo, &buf[0], uint32(len(buf)))
		if callErr == nil {
			return parseFileStreamInfo(buf)
		}
		if callErr != windows.ERROR_MORE_DATA {
			return nil, callErr
		}
		if size >= ceiling {
			return nil, fmt.Errorf("named streams of %s did not fit in %d bytes (%v): "+
				"the list is incomplete, so this is an error rather than an empty answer",
				path, size, callErr)
		}
		next := size * 4
		if next > ceiling {
			next = ceiling
		}
		size = next
	}
}

// parseFileStreamInfo walks the records one call produced. Every length it
// trusts is checked against the buffer first: the bytes came from the OS, and a
// record that does not fit is a failure to read, not an empty list.
func parseFileStreamInfo(buf []byte) ([]string, error) {
	var out []string
	for off := 0; ; {
		if off+streamInfoHeaderBytes > len(buf) {
			return nil, fmt.Errorf("stream record at offset %d runs past the %d-byte buffer", off, len(buf))
		}
		rec := (*fileStreamInfo)(unsafe.Pointer(&buf[off]))
		nameAt := off + streamInfoHeaderBytes
		if rec.StreamNameLength%2 != 0 || nameAt+int(rec.StreamNameLength) > len(buf) {
			return nil, fmt.Errorf("stream record at offset %d names %d bytes of UTF-16, which does not fit the %d-byte buffer",
				off, rec.StreamNameLength, len(buf))
		}
		name := ""
		if chars := int(rec.StreamNameLength) / 2; chars > 0 {
			name = windows.UTF16ToString(unsafe.Slice((*uint16)(unsafe.Pointer(&buf[nameAt])), chars))
		}
		if name != "" && name != "::$DATA" {
			out = append(out, name)
		}
		if rec.NextEntryOffset == 0 {
			return out, nil
		}
		step := int(rec.NextEntryOffset)
		if step <= 0 || off+step > len(buf) {
			return nil, fmt.Errorf("stream record at offset %d points %d bytes ahead, past the %d-byte buffer", off, step, len(buf))
		}
		off += step
	}
}
