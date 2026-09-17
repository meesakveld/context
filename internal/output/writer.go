package output

import (
	"archive/zip"
	"bytes"
	"os"

	"github.com/meesakveld/context/internal/scanner"
)

func WriteFile(
	path string,
	data []byte,
) error {
	return os.WriteFile(
		path,
		data,
		0644,
	)
}

func WriteZip(
	path string,
	result scanner.ContextResult,
) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, file := range result.Files {
		w, err := zipWriter.Create(file.Path)
		if err != nil {
			zipWriter.Close()
			return err
		}
		_, err = w.Write([]byte(file.Content))
		if err != nil {
			zipWriter.Close()
			return err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}