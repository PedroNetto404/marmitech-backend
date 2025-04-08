package ports

const DEFAULT_BUCKET = "default"

type (
	IBlockStorage interface {
		Save(key string, bucket string, data []byte) (string, error)
		Get(key, bucket string) ([]byte, error)
		Delete(key, bucket string) error
	}
)
