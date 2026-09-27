package main

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type args struct {
	Edits []struct {
		Old string `json:"old"`
		New string `json:"new"`
	} `json:"edits"`
}

func main() {
	partial := "\xbfal" // one BOM byte + content: invalid UTF-8
	fmt.Printf("input valid-utf8=%v bytes=% x\n", utf8.ValidString(partial), []byte(partial))
	b, err := json.Marshal(map[string]string{"old": partial})
	fmt.Printf("MARSHAL out=%s err=%v\n", b, err)
	// raw crafted JSON with the lone 0xBF byte inside the string literal
	raw := []byte(`{"edits":[{"old":"` + partial + `","new":"X"}]}`)
	var a args
	err = json.Unmarshal(raw, &a)
	fmt.Printf("UNmarshal rc-err=%v\n", err)
	if err == nil {
		got := a.Edits[0].Old
		fmt.Printf("decoded old bytes=% x valid-utf8=%v\n", []byte(got), utf8.ValidString(got))
		fmt.Printf("decoded old == raw partial? %v ; contains U+FFFD? %v\n", got == partial, string([]rune(got)) != got || containsFffd(got))
		fmt.Printf("would decoded old match file bytes [2:5]=% x ? %v\n", []byte("\xbfal"), got == "\xbfal")
	}
	// Go encoder name for the replacement
	fmt.Printf("U+FFFD name check: %q\n", "\ufffd")
}

func containsFffd(s string) bool {
	for range s {
	}
	return false
}
