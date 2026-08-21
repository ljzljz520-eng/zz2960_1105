package archive

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func Chunk(data []byte, size int) [][]byte {
	if size <= 0 {
		size = 1
	}
	chunks := make([][]byte, 0)
	for len(data) > 0 {
		end := size
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, append([]byte(nil), data[:end]...))
		data = data[end:]
	}
	return chunks
}

func Join(chunks [][]byte) []byte {
	var result []byte
	for _, chunk := range chunks {
		result = append(result, chunk...)
	}
	return result
}
