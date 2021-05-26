package archive

// Args captures the argument for the command
type Args struct {
	SourcePaths []string `arg:"source_paths"`
	OutputFile  string   `arg:"output_file" optional:"true"`
}

type Result struct {
	OutputFile string
	Error      string
	Size       int64
}
