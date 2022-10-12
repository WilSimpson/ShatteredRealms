package option

import (
    "flag"
    "fmt"
    "os"
    "reflect"
)

// ConfigOption Stores the value for the configuration
type ConfigOption struct {
    Value   string
    Default string
    EnvVar  string
    Flag    string
    Usage   string
}

// Config Represent all the configuration options there are for the application
type Config struct {
    // Port The port for the application
    Port ConfigOption

    // Host The host address for the application
    Host ConfigOption

    // Mode The running mode for the application
    Mode ConfigOption

    // KeyDir The path to the keys for JWT auth
    KeyDir ConfigOption

    AccountsPort ConfigOption
    AccountsHost ConfigOption

    CharactersPort ConfigOption
    CharactersHost ConfigOption

    AgonesKeyFile       ConfigOption
    AgonesCertFile      ConfigOption
    AgonesCaCertFile    ConfigOption
    AgonesNamespace     ConfigOption
    AgonesAllocatorHost ConfigOption
    AgonesAllocatorPort ConfigOption
}

const (

    // LocalhostMode A config more for development where a localhost connections will be returned instead of finding
    // services within kubernetes
    LocalhostMode = "localhost"

    // DebugMode A config mode allowing for the most verbose logging. Used heavily for development.
    DebugMode = "debug"

    // ReleaseMode A config mode allowing for minimal logging. Intended to be used for production.
    ReleaseMode = "release"
)

var (
    DefaultConfig = Config{
        Port: ConfigOption{
            Default: "8888",
            EnvVar:  "SRO_GAMEBACKEND_PORT",
            Flag:    "port",
            Usage:   "The port for the application",
        },
        Host: ConfigOption{
            Default: "",
            EnvVar:  "SRO_GAMEBACKEND_HOST",
            Flag:    "host",
            Usage:   "The host address for the application",
        },
        Mode: ConfigOption{
            Default: LocalhostMode,
            EnvVar:  "SRO_GAMEBACKEND_MODE",
            Flag:    "mode",
            Usage:   "The running mode for the application",
        },
        KeyDir: ConfigOption{
            Default: "/etc/sro/auth",
            EnvVar:  "SRO_KEY_DIR",
            Flag:    "keys",
            Usage:   "The path to the keys for JWT auth",
        },
        AccountsPort: ConfigOption{
            Default: "8080",
            EnvVar:  "SRO_ACCOUNTS_PORT",
            Flag:    "accounts_port",
            Usage:   "The port for the application",
        },
        AccountsHost: ConfigOption{
            Default: "",
            EnvVar:  "SRO_ACCOUNTS_HOST",
            Flag:    "accounts_host",
            Usage:   "The host address for the application",
        },
        CharactersPort: ConfigOption{
            Default: "8081",
            EnvVar:  "SRO_CHARACTERS_PORT",
            Flag:    "characters_port",
            Usage:   "The port for the application",
        },
        CharactersHost: ConfigOption{
            Default: "",
            EnvVar:  "SRO_CHARACTERS_HOST",
            Flag:    "characters_host",
            Usage:   "The host address for the application",
        },
        AgonesAllocatorHost: ConfigOption{
            Default: "",
            EnvVar:  "SRO_AGONES_HOST",
            Flag:    "agones_host",
            Usage:   "host for agones allocator server",
        },
        AgonesAllocatorPort: ConfigOption{
            Default: "443",
            EnvVar:  "SRO_AGONES_PORT",
            Flag:    "agones_port",
            Usage:   "port for agones allocator server",
        },
        AgonesNamespace: ConfigOption{
            Default: "default",
            EnvVar:  "SRO_AGONES_NS",
            Flag:    "agones_ns",
            Usage:   "kubernetes namespace to search for game servers",
        },
        AgonesKeyFile: ConfigOption{
            Default: "/etc/sro/auth/agones/client/key",
            EnvVar:  "SRO_AGONES_KEY",
            Flag:    "agones_key",
            Usage:   "the public key file for the client certificate in PEM format",
        },
        AgonesCertFile: ConfigOption{
            Default: "/etc/sro/auth/agones/client/cert",
            EnvVar:  "SRO_AGONES_CERT",
            Flag:    "agones_cert",
            Usage:   "the public key file for the client certificate in PEM format",
        },
        AgonesCaCertFile: ConfigOption{
            Default: "/etc/sro/auth/agones/ca/ca",
            EnvVar:  "SRO_AGONES_CA",
            Flag:    "agones_ca",
            Usage:   "the CA cert for server signing certificate in PEM format",
        },
    }
)

func NewConfig() Config {
    config := DefaultConfig
    config.readFlags()
    config.readEnvs()
    return config
}

// Address Gets the full address for the HTTP server
func (c *Config) Address() string {
    return fmt.Sprintf("%s:%s", c.Host.Value, c.Port.Value)
}

func (c *Config) AccountsAddress() string {
    return fmt.Sprintf("%s:%s", c.AccountsHost.Value, c.AccountsPort.Value)
}

func (c *Config) CharactersAddress() string {
    return fmt.Sprintf("%s:%s", c.CharactersHost.Value, c.CharactersPort.Value)
}

func (c *Config) AgonesAllocatorAddress() string {
    return fmt.Sprintf("%s:%s", c.AgonesAllocatorHost.Value, c.AgonesAllocatorPort.Value)
}

// IsRelease Returns true if the application mode is release
func (c *Config) IsRelease() bool {
    return c.Mode.Value == ReleaseMode
}

func (c *Config) readFlags() {
    v := reflect.ValueOf(c)

    temp := make([]*string, v.Elem().NumField())

    for i := 0; i < v.Elem().NumField(); i++ {
        field := v.Elem().Field(i)
        if field.Type() == reflect.TypeOf(ConfigOption{}) {
            option := field.Interface().(ConfigOption)
            temp[i] = flag.String(option.Flag, option.Default, option.Usage)
        }
    }

    flag.Parse()

    for i := 0; i < v.Elem().NumField(); i++ {
        coField := v.Elem().Field(i)
        if coField.Type() == reflect.TypeOf(ConfigOption{}) {
            valField := coField.FieldByName("Value")
            if temp[i] != nil {
                valField.SetString(*temp[i])
            }
        }
    }
}

func (c *Config) readEnvs() {
    v := reflect.ValueOf(c)
    for i := 0; i < v.Elem().NumField(); i++ {
        field := v.Elem().Field(i)
        if field.Type() == reflect.TypeOf(ConfigOption{}) {
            option := field.Interface().(ConfigOption)
            env, found := os.LookupEnv(option.EnvVar)
            if found && !isFlagPassed(option.Flag) {
                field.FieldByName("Value").SetString(env)
            }
        }
    }
}

func isFlagPassed(name string) bool {
    found := false
    flag.Visit(func(f *flag.Flag) {
        if f.Name == name {
            found = true
        }
    })

    return found
}
