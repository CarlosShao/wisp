module github.com/CarlosShao/wisp/probes175r1

go 1.27

require github.com/CarlosShao/wisp v0.0.0

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.48.0 // indirect
	modernc.org/libc v1.75.7 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
	modernc.org/sqlite v1.59.0 // indirect
)

// The rig lives inside the repo tree but in its own module, so the internal/*
// packages it drives come from the checked-out source, never from a published
// copy. Five levels up: r1 -> 175 -> probes -> wisp -> .scratch -> repo root.
replace github.com/CarlosShao/wisp => ../../../../../
