package files

import (
	"fmt"
	"os"
)

type CloudBlockStorage struct {
}

func NewCloudBlockStorage() *CloudBlockStorage {
	return &CloudBlockStorage{}
}

func (c *CloudBlockStorage) Save(key string, bucket string, data []byte) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CloudBlockStorage) Get(key, bucket string) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CloudBlockStorage) Delete(key, bucket string) error {
	// TODO implement me
	panic("implement me")
}

type DiskStorage struct {
	basePath string
}

func NewDiskStorage(
	basePath string,
) *DiskStorage {
	return &DiskStorage{
		basePath: basePath,
	}
}

func (d *DiskStorage) Save(key string, bucket string, data []byte) (string, error) {
	bucketPath := fmt.Sprintf("%s/%s", d.basePath, bucket)
	if err := d.createIfNotExists(bucketPath); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", bucketPath, key)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (d *DiskStorage) Get(key, bucket string) ([]byte, error) {
	bucketPath := fmt.Sprintf("%s/%s", d.basePath, bucket)
	filePath := fmt.Sprintf("%s/%s", bucketPath, key)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *DiskStorage) Delete(key, bucket string) error {
	bucketPath := fmt.Sprintf("%s/%s", d.basePath, bucket)
	filePath := fmt.Sprintf("%s/%s", bucketPath, key)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiskStorage) createIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
