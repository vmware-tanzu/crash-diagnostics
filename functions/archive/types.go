package archive

// Params captures the argument for the command
type Params struct {
	SourcePaths []string
	OutputFile  string
	Size        uint64
}

// Archive is used to represent the output of
// an archive command
type Archive = Params
