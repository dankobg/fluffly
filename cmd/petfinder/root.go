package petfinder

type RootCmd struct {
	Auth            AuthCmd                  `cmd:"" help:"get access token for petfinder"`
	DownloadOrgs    DownloadOrganizationsCmd `cmd:"" help:"Download petfinder organizations data"`
	DownloadAnimals DownloadAnimalsCmd       `cmd:"" help:"Download petfinder animals data"`
	DownloadCommon  DownloadCommonCmd        `cmd:"" help:"Download petfinder common data"`
}
