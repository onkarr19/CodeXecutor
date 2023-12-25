package models

type CompilationResult struct {
	ExitCode int    // Indicates exit code of container
	Output   string // Compiler output or execution results
	Error    error  // Compilation or execution errors, if any
}
