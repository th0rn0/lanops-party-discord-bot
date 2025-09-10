package config

type Config struct {
	Discord Discord
	Lanops  Lanops
}

type Discord struct {
	CommandPrefix string
	Token         string
	GuildId       string
	AdminRoleId   string
}

type Lanops struct {
	StreamProxyApiAddress string
	JukeboxApiUsername    string
	JukeboxApiPassword    string
	JukeboxApiUrl         string
}
