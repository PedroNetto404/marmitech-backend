package ports

const DEFAULT_BUCKET = "default"

type (
	IFileStorage interface {
		Save(key string, bucket string, data []byte) error
		Get(key, bucket string) ([]byte, error)
	}
)