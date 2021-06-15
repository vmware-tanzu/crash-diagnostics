package archive

// Args captures the argument for the command
type Args struct {
	SourcePaths []string `name:"source_paths"`
	OutputFile  string   `name:"output_file" optional:"true"`
}

type Result struct {
	Error   string  `name:"error"`
	Archive Archive `name:"archive"`
}

type Archive struct {
	OutputFile string `name:"output_file"`
	Error      string `name:"error"`
	Size       int64  `name:"size"`
}
