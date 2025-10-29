package petfinder

import (
	"fmt"
)

type DownloadCommonCmd struct {
	Dir string `required:""`
}

func (dc *DownloadCommonCmd) Run() error {
	fmt.Printf("DOWNLOADING COMMON DATA INTO: %q\n", dc.Dir)
	return nil
}
