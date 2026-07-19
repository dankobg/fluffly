package identities

type RootCmd struct {
	ImportIdentities ImportIdentitiesCmd `cmd:"" help:"Import identities and permissions"`
}
