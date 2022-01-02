package server

type Server struct {
	Name   string `hcl:"name,label"`
	Path   string `hcl:"path"`
	Stable bool   `hcl:"stable"`
	Steam  string `hcl:"steam"`

	FixesVersion string `hcl:"server_fixes_version"`
}
