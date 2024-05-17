package storage

import (
	"context"
	"fmt"

	"github.com/khaledhikmat/threat-detection-shared/equates"
)

func storeClipViaAzure(_ context.Context, clip equates.RecordingClip) (string, error) {
	return fmt.Sprintf("https://%s/azure_url", clip.LocalReference), nil
}

func retrieveClipFromAzure(_ context.Context, _ equates.RecordingClip) ([]byte, error) {
	return []byte{}, nil
}

func downloadClipFromAzure(_ context.Context, _ equates.RecordingClip) ([]byte, error) {
	return []byte{}, nil
}
