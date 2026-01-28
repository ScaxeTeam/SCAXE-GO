package version

import "fmt"

const (
	Major = 0

	Minor = 2

	Patch = 0

	Codename = "Entity & Event"

	MinecraftVersion = "0.14.3"

	ProtocolVersion = 70
)

func String() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}

func Full() string {
	return fmt.Sprintf("Scaxe Go v%s \"%s\" (MCPE %s, Protocol %d)",
		String(), Codename, MinecraftVersion, ProtocolVersion)
}
