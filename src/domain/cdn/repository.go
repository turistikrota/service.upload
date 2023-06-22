package cdn

type Repository interface {
	Upload(file []byte, filename string, path ...string) (string, bool)
}
