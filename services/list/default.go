package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
)

func List() error {
	services, err := listServices("./")
	if err != nil {
		return err
	}

	display(services)

	return nil
}

func listServices(root string) ([]string, error) {
	var services []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "error walking directory")
		}

		for _, service := range services {
			if strings.Contains(path, service) {
				return nil // we don't support nested services
			}
		}

		if common.IsServiceDir(path) {
			services = append(services, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return services, nil
}

func display(services []string) {
	for _, service := range services {
		fmt.Println(service)
	}
}
