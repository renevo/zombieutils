server "Burpcraft" {
  path            = "./burpcraft/zombie/stable"
  config          = "./burpcraft/zombie/stable.xml"
  save_folder     = "./burpcraft/zombie/stable/GameData"
  admin_file_name = "admin.xml"
  experimental    = false
  steam           = "/usr/games/steamcmd"

  server_fixes_version = "22.24.39"

  admin "Dante" {
    id         = 76561197969618392
    permission = 0
  }

  whitelist "Dante" { id = 76561197969618392 }
  whitelist "fallendice" { id = 76561198165130563 }
  whitelist "ShoeDawg" { id = 76561198126980759 }
  whitelist "rmdashrrootsplat" { id = 76561197971008541 }
  whitelist "Acatera" { id = 76561198067220232 }

  #whitelist "Darkside916" { id = 76561198169213629 }
  #whitelist "Bruneo" { id = 76561198093708037 }
  #whitelist "grinanberrett" { id = 76561198052121335 }
  #whitelist "Thur" { id = 76561197961074073 }
  #whitelist "Wrestleeagle" { id = 76561198068378636 }
  #whitelist "Beanie The Alien" { id = 76561198070203187 }
  #whitelist "Lucy Shepard" { id = 76561197972274364 }
  #whitelist "N7omad" { id = 76561197968569725 }
  #whitelist "Rysarth" { id = 76561198074646046 }
  #whitelist "phancy" { id = 76561198440484535 }
  #whitelist "Ceanox" { id = 76561198069242633 }
  #whitelist "Krayvin" { id = 76561198067189139 }
  #whitelist "Decaine" { id = 76561197970820129 }

  permission "chunkcache" { level = 1000 }
  permission "debugshot" { level = 1000 }
  permission "debugweather" { level = 1000 }
  permission "getgamepref" { level = 1000 }
  permission "getgamestat" { level = 1000 }
  permission "getoptions" { level = 1000 }
  permission "gettime" { level = 1000 }
  permission "gfx" { level = 1000 }
  permission "help" { level = 1000 }
  permission "memcl" { level = 1000 }
  permission "settempunit" { level = 1000 }
  permission "listplayerids" { level = 1000 }
  permission "listthreads" { level = 1000 }

  webpermission "web.map" { level = 2000 }
  webpermission "webapi.getstats" { level = 2000 }
  webpermission "webapi.getplayersonline" { level = 2000 }
  webpermission "webapi.getplayerslocation" { level = 2000 }
  webpermission "webapi.getlandclaims" { level = 2000 }
  webpermission "webapi.viewallplayers" { level = 2000 }
  webpermission "webapi.viewallclaims" { level = 2000 }

  webtoken "CLI" {
    token = "supersecret"
    level = 0
  }
}
