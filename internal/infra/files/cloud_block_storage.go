package files

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