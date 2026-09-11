package output

import "os"

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
