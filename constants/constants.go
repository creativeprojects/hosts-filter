package constants

const (
	BufferInitialLines   = 20000
	BufferInitialEntries = 30000
)

var (
	IgnoreDomains = []string{
		"localhost",
		"localhost.localdomain",
		"broadcasthost",
		"local",
		"ip6-localhost",
		"ip6-loopback",
		"ip6-localnet",
		"ip6-mcastprefix",
		"ip6-allnodes",
		"ip6-allrouters",
		"ip6-allhosts",
	}
)
