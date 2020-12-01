package lib

import (
	"github.com/acciaioli/mono/internal/common"
)

func List(bs common.BlobStorage) ([]ServiceState, error) {
	services, err := LoadServices()
	if err != nil {
		return nil, err
	}

	return GetServicesState(bs, services)
}
