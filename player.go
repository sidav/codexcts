package main

type player struct {
	hand    []*card
	draw    []*card
	discard []*card
	heroes  [3]*card
	gold    int
	workers int
}
