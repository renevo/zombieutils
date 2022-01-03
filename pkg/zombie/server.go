package zombie

type Server struct {
	Name         string `hcl:"name,label"`
	Path         string `hcl:"path"`
	Experimental bool   `hcl:"experimental"`
	Steam        string `hcl:"steam"`
	Config       string `hcl:"config"`

	FixesVersion string `hcl:"server_fixes_version"`
}

func Default() *Server {
	return &Server{
		Name:         "Burpcraft",
		Path:         "./burpcraft/zombie/stable",
		Experimental: false,
		Steam:        "steamcmd",
		Config:       "./burpcraft/zombie/stable/serverconfig.xml",
	}
}
