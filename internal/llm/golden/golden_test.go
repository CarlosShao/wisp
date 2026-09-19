package golden

import (
	"net/http"
	"strings"
	"testing"
)

func TestParseDirectiveSections(t *testing.T) {
	src := "# wisp golden sse v1\n" +
		"# @scenario ladder\n" +
		"# @response 429 retry-after:2\n" +
		"{\"error\":\"slow down\"}\n" +
		"# @response 200\n" +
		"data: {\"a\":1}\n" +
		"\n" +
		"data: [DONE]\n"

	rs, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 2 {
		t.Fatalf("sections = %d, want 2", len(rs))
	}
	if rs[0].Status != 429 || rs[0].Header.Get("Retry-After") != "2" {
		t.Errorf("section 0 = %+v", rs[0])
	}
	if ct := rs[0].Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("section 0 content-type = %q, want application/json", ct)
	}
	if string(rs[0].Body) != "{\"error\":\"slow down\"}\n" {
		t.Errorf("section 0 body = %q", rs[0].Body)
	}
	if rs[1].Status != 200 || rs[1].Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("section 1 = %+v", rs[1])
	}
	want := "data: {\"a\":1}\n\ndata: [DONE]\n"
	if string(rs[1].Body) != want {
		t.Errorf("section 1 body = %q, want %q", rs[1].Body, want)
	}
}

func TestParseNoDirectivesIsSingleResponse(t *testing.T) {
	src := "# just a header comment\n\ndata: x\n\ndata: [DONE]\n"
	rs, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 {
		t.Fatalf("sections = %d, want 1", len(rs))
	}
	if rs[0].Status != 200 || rs[0].Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("implicit response = %+v", rs[0])
	}
	if !strings.Contains(string(rs[0].Body), "data: [DONE]") {
		t.Errorf("body = %q", rs[0].Body)
	}
}

func TestParseLatencyBeforeAndAfterResponse(t *testing.T) {
	src := "# @latency 30\n# @response 200\n# @latency 70\ndata: x\n"
	rs, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 {
		t.Fatalf("sections = %d", len(rs))
	}
	if rs[0].LatencyMS != 70 {
		t.Errorf("latency = %d, want 70 (the directive after @response wins)", rs[0].LatencyMS)
	}
}

func TestParseEscapedDirectiveBodyLine(t *testing.T) {
	src := "# @response 200\n# @@data: literal\n"
	rs, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if string(rs[0].Body) != "# @data: literal\n" {
		t.Errorf("body = %q, want escaped line unescaped", rs[0].Body)
	}
}

func TestParseRejectsUnknownDirectiveAndBadStatus(t *testing.T) {
	if _, err := Parse([]byte("# @wat 1\ndata: x\n")); err == nil {
		t.Error("unknown directive accepted")
	}
	if _, err := Parse([]byte("# @response 9999\ndata: x\n")); err == nil {
		t.Error("status 9999 accepted")
	}
}

func TestWriteParseRoundTrip(t *testing.T) {
	rs := []Response{
		{Status: 429, Header: http.Header{"Retry-After": []string{"3"}}, Body: []byte("{\"e\":1}\n")},
		{Status: 200, LatencyMS: 25, Body: []byte("data: hi\n\ndata: [DONE]\n")},
	}
	var sb strings.Builder
	Write(&sb, "round trip", rs)
	back, err := Parse([]byte(sb.String()))
	if err != nil {
		t.Fatalf("reparse: %v\n%s", err, sb.String())
	}
	if len(back) != 2 {
		t.Fatalf("round trip sections = %d", len(back))
	}
	if back[0].Status != 429 || back[0].Header.Get("Retry-After") != "3" || string(back[0].Body) != "{\"e\":1}\n" {
		t.Errorf("round trip section 0 = %+v", back[0])
	}
	if back[1].LatencyMS != 25 || string(back[1].Body) != "data: hi\n\ndata: [DONE]\n" {
		t.Errorf("round trip section 1 = %+v", back[1])
	}
}

func TestReplayerSequentialConsumption(t *testing.T) {
	rs, err := Parse([]byte("# @response 500\nbad\n# @response 200\ndata: ok\n"))
	if err != nil {
		t.Fatal(err)
	}
	rep := NewReplayer(rs)
	srv := rep.Server()
	defer srv.Close()

	resp1, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp1.StatusCode != 500 {
		t.Errorf("first response status = %d, want 500", resp1.StatusCode)
	}
	resp1.Body.Close()

	resp2, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Errorf("second response status = %d, want 200", resp2.StatusCode)
	}

	resp3, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != 599 {
		t.Errorf("exhausted response status = %d, want 599", resp3.StatusCode)
	}
	if len(rep.Requests) != 3 {
		t.Errorf("requests recorded = %d, want 3", len(rep.Requests))
	}
}
