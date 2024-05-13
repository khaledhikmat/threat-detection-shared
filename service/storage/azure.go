package storage

import "fmt"

func storeFileViaAzure(path2File string) (string, error) {
	return fmt.Sprintf("https://%s/azure_url", path2File), nil
}
