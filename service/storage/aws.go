package storage

import "fmt"

func storeFileViaAWS(path2File string) (string, error) {
	return fmt.Sprintf("https://%s/aws_url", path2File), nil
}
