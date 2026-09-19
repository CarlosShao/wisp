//go:build !windows

package audio

// Non-Windows placeholder: Wisp is a Windows product (D41d ships
// windows/amd64 only), but the package must still compile on other GOOS
// values so accidental cross-packages references fail loudly, not silently.

import (
	"context"
	"errors"
)

// WASAPIMicrophone is unavailable off Windows; Start fails loudly.
type WASAPIMicrophone struct{}

// NewWASAPIMicrophone returns the placeholder.
func NewWASAPIMicrophone() *WASAPIMicrophone { return &WASAPIMicrophone{} }

// Start always fails off Windows.
func (m *WASAPIMicrophone) Start(ctx context.Context, buf chan<- []byte) error {
	return errors.New("audio: WASAPI capture requires Windows (D41d)")
}

// Stop is a no-op off Windows.
func (m *WASAPIMicrophone) Stop() error { return nil }

// Stats reports no counters off Windows.
func (m *WASAPIMicrophone) Stats() Stats { return Stats{} }
