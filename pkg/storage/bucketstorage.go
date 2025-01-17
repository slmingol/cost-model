package storage

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// StorageProvider is the type of provider used for storage if not leveraging a file implementation.
type StorageProvider string

const (
	S3 StorageProvider = "S3"
	// AZURE StorageProvider = "AZURE"
	// GCS   StorageProvider = "GCS"
)

// StorageConfig is the configuration type used as the "parent" configuration. It contains a type, which will
// specify the bucket storage implementation, and a configuration object specific to that storage implementation.
type StorageConfig struct {
	Type   StorageProvider `yaml:"type"`
	Config interface{}     `yaml:"config"`
}

// NewBucketStorage initializes and returns new Storage implementation leveraging the storage provider
// configuration. This configuration type uses the layout provided in thanos: https://thanos.io/tip/thanos/storage.md/
func NewBucketStorage(config []byte) (Storage, error) {
	storageConfig := &StorageConfig{}
	if err := yaml.UnmarshalStrict(config, storageConfig); err != nil {
		return nil, errors.Wrap(err, "parsing config YAML file")
	}

	// Because the Config property is specific to the storage implementation, we'll marshal back into yaml, and allow
	// the specific implementation to unmarshal back into a concrete configuration type.
	config, err := yaml.Marshal(storageConfig.Config)
	if err != nil {
		return nil, errors.Wrap(err, "marshal content of storage configuration")
	}

	var storage Storage
	switch strings.ToUpper(string(storageConfig.Type)) {
	case string(S3):
		storage, err = NewS3Storage(config)
	//case string(GCS):
	//	storage, err = NewGCSStorage(config)
	//case string(AZURE):
	//	storage, err = NewAzureStorage(config)
	default:
		return nil, errors.Errorf("storage with type %s is not supported", storageConfig.Type)
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("create %s client", storageConfig.Type))
	}

	return storage, nil
}
