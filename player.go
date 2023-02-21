package main

type player struct {
	hand    []card
	draw    []card
	discard []card
	// heroes  [3]*card need a separate type
	gold    int
	workers int
}
