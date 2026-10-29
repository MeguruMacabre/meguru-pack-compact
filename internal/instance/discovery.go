package instance

import "os"

// FindCandidates возвращает список директорий, доступных для выбора Minecraft-инстанса.
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
