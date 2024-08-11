package safelock

type fileChunk struct {
	Chunk  []byte
	Sought int
}
