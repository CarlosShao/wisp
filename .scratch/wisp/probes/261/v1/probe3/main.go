// Ticket 261 leg v1 throwaway probe - three production entries into the gate,
// measured live (dispatch order 3). Deleted-nothing discipline: this file
// stays (build-only-never-deleted), it is scratch, zero tracked code touched.
//
// Fixture: hand-written TOML via config.LoadFile (the ticket-257 first-run
// shape). base_url is intentionally unreachable: the gate must refuse BEFORE
// any dial, and the flipped-true control only proves selection admits.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
)

const tomlT = `schema_version = 2

[llm]
text_chain = ["mock261/t-off"]

[llm.providers.mock261]
protocol = "openai-chat"
base_url = "http://127.0.0.1:9/v1"

[llm.providers.mock261.models.t-nokey]
context_window = 4096

[llm.providers.mock261.models.t-off]
context_window = 4096
enabled = %s

[llm.providers.mock261.models.t-on]
context_window = 4096
enabled = true
`

func load(path, flag string) *config.Config {
	body := fmt.Sprintf(tomlT, flag)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		panic(err)
	}
	cfg, _, err := config.LoadFile(path, nil)
	if err != nil {
		panic(fmt.Sprintf("LoadFile(%s): %v", flag, err))
	}
	return cfg
}

func main() {
	dir, err := os.MkdirTemp("", "t261v1")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "config.toml")

	// ---- pass 1: t-off disabled (as the user left it) ----
	cfg := load(path, "false")
	res := llm.NewResolver(cfg, nil)

	// Entry 1: text_chain names the disabled model.
	eps, err := res.ResolveChain()
	fmt.Printf("ENTRY1 text_chain ResolveChain: err=%q eps=%v\n", showErr(err), eps)

	// Entry 2: roles.chat unset -> fallback onto the chain head (run.go:436 shape).
	_, _, err = res.ResolveRole(llm.RoleChat)
	fmt.Printf("ENTRY2 role-fallback ResolveRole(chat): err=%q\n", showErr(err))

	// Entry 2b (bonus): roles.chat DIRECTLY names the disabled model.
	cfg2 := load(path, "false")
	cfg2.LLM.Roles.Chat.Provider = "mock261"
	cfg2.LLM.Roles.Chat.Model = "t-off"
	res2 := llm.NewResolver(cfg2, nil)
	_, _, err = res2.ResolveRole(llm.RoleChat)
	fmt.Printf("ENTRY2b role-direct ResolveRole(chat): err=%q\n", showErr(err))

	// Entry 3: enumeration.
	fmt.Printf("ENTRY3 DiscoveredModels(mock261) = %v\n", res.DiscoveredModels("mock261"))

	// ---- pass 2: the SAME file, only the flag flipped true (hand re-enable) ----
	cfg3 := load(path, "true")
	res3 := llm.NewResolver(cfg3, nil)
	eps3, err := res3.ResolveChain()
	fmt.Printf("FLIP text_chain ResolveChain: err=%q eps=%v\n", showErr(err), eps3)
	_, _, err = res3.ResolveRole(llm.RoleChat)
	fmt.Printf("FLIP role-fallback ResolveRole(chat): err=%q\n", showErr(err))
	fmt.Printf("FLIP DiscoveredModels(mock261) = %v\n", res3.DiscoveredModels("mock261"))
}

func showErr(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}
