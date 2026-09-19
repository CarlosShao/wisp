//go:build !cgo_sherpa

package main

import "errors"

var errAsrCreate = errors.New("ASR session creation failed")

func openSession(kind, modelsDir string) (*sessionHandle, error) {
	return nil, errors.New("unload-test requires the cgo_sherpa build of xy-verdict")
}
