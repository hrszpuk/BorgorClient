package main

// RPS! (thats all)

// the programs entry point
func main() {
	// i stole some code from rgoc because i really like its cli
	// all these functions can be found in cli.go

	// Init defines all the flags and initializes them
	Init()

	// ProcessFlags does as it says, takes the flags from Init and uses them to run parts of the compiler
	ProcessFlags()
}
