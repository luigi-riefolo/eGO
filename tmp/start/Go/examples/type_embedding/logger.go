package main

import (
	"log"
	"os"
)

// Job ...
type Job struct {
	Command string
	Logger  *log.Logger
}

// JobOne ...
type JobOne struct {
	Command string
	*log.Logger
}

func main() {
	job := &Job{"demo", log.New(os.Stderr, "Job: ", log.Ldate)}
	jobOne := &JobOne{"demo", log.New(os.Stderr, "Job: ", log.Ldate)}
	// same as
	// job := &Job{Command: "demo",
	// Logger: log.New(os.Stderr, "Job: ", log.Ldate)}
	job.Logger.Print("test")
	// Or
	jobOne.Print("test one")
}
