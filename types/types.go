package types

// PutObjectInput structure for PutObject
type PutObjectInput struct {
	Key         string
	Value       string
	Application string
	KmsKeyAlias string
}

// GetObjectInput structure for GetObject
type GetObjectInput struct {
	Key     string
	Version string
}

// GetObjectOuput structure for GetObject
type GetObjectOutput struct {
	Value string
}
