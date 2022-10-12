package main

var (
	conf = &config{
		Port:         8081,
		Host:         "",
		Mode:         "development",
		KeyDir:       "/etc/sro/auth",
		DBFile:       "/etc/sro/db.yaml",
		AccountsPort: 8080,
		AccountsHost: "",
	}
)
