package models

import "github.com/CarlosShao/wisp/internal/observe"

// Package-local aliases of the D37 classes used here (keeps call sites short
// without redefining the taxonomy).
const (
	ClassModel     = observe.ClassModel
	ClassNetwork   = observe.ClassNetwork
	ClassConfig    = observe.ClassConfig
	ClassCancelled = observe.ClassCancelled
)

// D37 error-class helpers for the models domain: every failure this package
// returns is a classified observe.Error. Mapping follows D37: model missing /
// hash / signature problems are ClassModel; transport problems are
// ClassNetwork; user or owner cancellation is ClassCancelled; a caller
// trying to construct the manager with signature verification disabled is
// ClassConfig (C29 hard error, mirrors the config-layer rejection).

func observeNew(class observe.ErrorClass, detail string) *observe.Error {
	return observe.New(class, detail)
}

func observeWrap(class observe.ErrorClass, detail string, err error) *observe.Error {
	return observe.Wrap(class, err, detail)
}
