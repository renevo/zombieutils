server "Burpcraft" {
  path            = "./burpcraft/zombie/stable"
  config          = "./burpcraft/zombie/stable.xml"
  save_folder     = "./burpcraft/zombie/stable/GameData"
  admin_file_name = "admin.xml"
  steam           = "/usr/games/steamcmd"
  clean_mods      = false

  admin "Dante" {
    id         = 76561197969618392
    permission = 0
    platform   = "Steam"
  }

  whitelist "Dante" {
    id       = 76561197969618392
    platform = "Steam"
  }
  whitelist "fallendice" {
    id       = 76561198165130563
    platform = "Steam"
  }
  whitelist "ShoeDawg" {
    id       = 76561198126980759
    platform = "Steam"
  }
  whitelist "rmdashrrootsplat" {
    id       = 76561197971008541
    platform = "Steam"
  }
  whitelist "Acatera" {
    id       = 76561198067220232
    platform = "Steam"
  }
  whitelist "atownmanx" {
    id       = 76561198056550775
    platform = "Steam"
  }
  whitelist "Afka" {
    id       = 76561199039705922
    platform = "Steam"
  }

  permission "chunkcache" { level = 500 }
  permission "debugshot" { level = 500 }
  permission "debugweather" { level = 500 }
  permission "getgamepref" { level = 500 }
  permission "getgamestat" { level = 500 }
  permission "getoptions" { level = 500 }
  permission "gettime" { level = 500 }
  permission "gfx" { level = 500 }
  permission "help" { level = 500 }
  permission "memcl" { level = 500 }
  permission "settempunit" { level = 500 }
  permission "listplayerids" { level = 500 }
  permission "listthreads" { level = 500 }
  permission "graph" { level = 500 }
  permission "loot" { level = 500 }
  permission "uioptions" { level = 500 }
  permission "createwebuser" { level = 500 }
  permission "decomgr" { level = 500 }
  permission "getlogpath" { level = 500 }

  modpack "burpcraft" {
    url = "https://github.com/renevo/burpmod-7days/releases/download/1.0.0/burpcraft-1.0.0.zip"
  }
}
