package petfinder

import (
	"fmt"
)

type DownloadAnimalsCmd struct {
	Dir string `required:""`
}

func (dc *DownloadAnimalsCmd) Run() error {
	fmt.Printf("DOWNLOADING ANIMALS DATA INTO: %q\n", dc.Dir)
	return nil
}
