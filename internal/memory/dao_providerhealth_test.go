package memory

import (
	"context"
	"testing"
	"time"
)

// DAO primitive tests for provider_health (ticket 09 defines the primitive;
// ticket 11 writes real probe results).

func TestProviderHealthProbeRoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	vision, fc := true, false
	err := s.UpsertProviderProbe(ctx, "openai", "gpt-4o", ProbeFlags{
		Vision: &vision, FC: &fc,
	}, true, 640, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	row, err := s.ProviderHealthRow(ctx, "openai", "gpt-4o")
	if err != nil {
		t.Fatal(err)
	}
	if row.Provider != "openai" || row.Model != "gpt-4o" || row.QuotaState != QuotaOK {
		t.Errorf("row = %+v", row)
	}
	if row.Probe.Vision == nil || !*row.Probe.Vision || row.Probe.FC == nil || *row.Probe.FC {
		t.Errorf("probe flags = %+v", row.Probe)
	}
	if row.LastProbeAt.IsZero() || row.LatencyMSP50 != 640 {
		t.Errorf("probe metadata = %+v latency=%d", row.LastProbeAt, row.LatencyMSP50)
	}

	// Re-probe updates in place.
	err = s.UpsertProviderProbe(ctx, "openai", "gpt-4o", ProbeFlags{FC: &fc}, false, 700, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	row, err = s.ProviderHealthRow(ctx, "openai", "gpt-4o")
	if err != nil {
		t.Fatal(err)
	}
	if row.Probe.Vision != nil { // flags are the latest probe's verdict
		t.Errorf("re-probe should replace flags, got %+v", row.Probe)
	}
	if row.QuotaState != QuotaUnknown || row.LatencyMSP50 != 700 {
		t.Errorf("row after re-probe = %+v", row)
	}
}

func TestProviderHealthRecordErrorAndList(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Recording an error for an unseen pair creates the row.
	if err := s.RecordProviderError(ctx, "deepseek", "deepseek-chat", "rate_limit", "HTTP 429: slow down", QuotaThrottled); err != nil {
		t.Fatal(err)
	}
	row, err := s.ProviderHealthRow(ctx, "deepseek", "deepseek-chat")
	if err != nil {
		t.Fatal(err)
	}
	if row.LastError == "" || row.LastErrorAt.IsZero() || row.QuotaState != QuotaThrottled {
		t.Errorf("row = %+v", row)
	}
	if row.LastProbeAt.IsZero() == false {
		t.Errorf("probe timestamp must stay empty: %+v", row.LastProbeAt)
	}

	// A later probe result survives error recording (upsert-then-update).
	fc := true
	if err := s.UpsertProviderProbe(ctx, "deepseek", "deepseek-chat", ProbeFlags{FC: &fc}, true, 320, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordProviderError(ctx, "deepseek", "deepseek-chat", "provider", "HTTP 500", ""); err != nil {
		t.Fatal(err)
	}
	row, err = s.ProviderHealthRow(ctx, "deepseek", "deepseek-chat")
	if err != nil {
		t.Fatal(err)
	}
	if row.QuotaState != QuotaOK || row.LatencyMSP50 != 320 || row.LastError != "HTTP 500" {
		t.Errorf("row after probe+error = %+v", row)
	}

	// D37 enum is enforced at the DAO layer.
	if err := s.RecordProviderError(ctx, "p", "m", "not_a_class", "x", ""); err == nil {
		t.Error("invalid error_class accepted")
	}
	if err := s.RecordProviderError(ctx, "p", "m", "provider", "x", "bogus_quota"); err == nil {
		t.Error("invalid quota_state accepted")
	}

	// A second pair makes the sorted-list check meaningful.
	if err := s.UpsertProviderProbe(ctx, "openai", "gpt-4o", ProbeFlags{}, true, 100, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}

	// List returns both rows sorted.
	list, err := s.ListProviderHealth(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("list = %d rows, want 2", len(list))
	}
	if list[0].Provider != "deepseek" || list[1].Provider != "openai" {
		t.Errorf("list order = %s, %s", list[0].Provider, list[1].Provider)
	}
}

func TestProviderHealthMissingRowIsNoRows(t *testing.T) {
	s := openTestStore(t)
	if _, err := s.ProviderHealthRow(context.Background(), "nobody", "nothing"); err == nil {
		t.Error("expected sql.ErrNoRows for an unprobed pair")
	}
}
