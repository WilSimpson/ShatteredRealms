package main

var (
	conf = &config{
		Port:                8082,
		Host:                "",
		Mode:                "development",
		KeyDir:              "/etc/sro/auth",
		DBFile:              "/etc/sro/db.yaml",
		AccountsPort:        8080,
		AccountsHost:        "",
		CharactersPort:      8081,
		CharactersHost:      "",
		AgonesKeyFile:       "/etc/sro/auth/agones/client/key",
		AgonesCertFile:      "/etc/sro/auth/agones/client/cert",
		AgonesCaCertFile:    "/etc/sro/auth/agones/client/ca",
		AgonesNamespace:     "default",
		AgonesAllocatorHost: "",
		AgonesAllocatorPort: 442,
	}
)

func loadConfig() {

}
