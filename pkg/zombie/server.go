package zombie

type Server struct {
	Name          string `hcl:"name,label"`
	Path          string `hcl:"path"`
	Experimental  bool   `hcl:"experimental,optional"`
	Steam         string `hcl:"steam"`
	Config        string `hcl:"config"`
	SaveFolder    string `hcl:"save_folder"`
	AdminFileName string `hcl:"admin_file_name,optional"`

	FixesVersion string `hcl:"server_fixes_version,optional"`

	Admins         []ServerAdmin          `hcl:"admin,block"`
	Permissions    []ServerPermission     `hcl:"permission,block"`
	WebPermissions []ServerWebPermission  `hcl:"webpermission,block"`
	WebTokens      []ServerWebToken       `hcl:"webtoken,block"`
	Whitelist      []ServerWhitelistEntry `hcl:"whitelist,block"`
	Mods           []ServerMod            `hcl:"mod,block"`
	CleanMods      bool                   `hcl:"clean_mods,optional"`
}

type ServerAdmin struct {
	Name       string `hcl:"name,label" xml:"name,attr"`
	ID         int    `hcl:"id" xml:"steamID,attr"`
	Permission int    `hcl:"permission" xml:"permission_level,attr"`
}

type ServerPermission struct {
	Command string `hcl:"cmd,label" xml:"cmd,attr"`
	Level   int    `hcl:"level" xml:"permission_level,attr"`
}

type ServerWebPermission struct {
	Module string `hcl:"module,label" xml:"module,attr"`
	Level  int    `hcl:"level" xml:"permission_level,attr"`
}

type ServerWebToken struct {
	Name  string `hcl:"name,label" xml:"name,attr"`
	Token string `hcl:"token" xml:"token,attr"`
	Level int    `hcl:"level" xml:"permission_level,attr"`
}

type ServerWhitelistEntry struct {
	Name string `hcl:"name,label" xml:"name,attr"`
	ID   int    `hcl:"id" xml:"steamID,attr"`
}

func Default() *Server {
	return &Server{
		Name:          "Burpcraft",
		Path:          "./burpcraft/zombie/stable",
		Experimental:  false,
		Steam:         "steamcmd",
		Config:        "./burpcraft/zombie/stable/serverconfig.xml",
		AdminFileName: "serveradmin.xml",
	}
}
