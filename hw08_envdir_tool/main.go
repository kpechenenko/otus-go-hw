package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Println("Usage: go-envdir dir command [args]")
		os.Exit(ReturnCodeErr)
	}
	// log.Printf("args: %v\n", args)
	env, err := ReadDir(args[1])
	if err != nil {
		log.Printf("Error reading env from dir %s: %s\n", args[1], err)
		os.Exit(ReturnCodeErr)
	}
	log.Printf("readed env: %+v\n", env)
	exitCode := RunCmd(args[2:], env)
	if exitCode != ReturnCodeOk {
		log.Println("Error running command.")
	}
	os.Exit(exitCode)
}
