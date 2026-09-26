package instance

import "os"

func FindCandidates(root string) ([]string, error) {

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var directories []string

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	return directories, nil
}
