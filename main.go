package main

import (
	"context"

	"github.com/kevpar/repl-go"
)

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	return repl.Run(&state{systems: make(map[string]*cs)}, allCommands(), func(state *state) string { return state.def })
}
