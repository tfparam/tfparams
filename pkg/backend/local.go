package backend

import "os"

func fetchLocal(s Source) ([]byte, error) {
	return os.ReadFile(s.Path) //nolint:gosec // path is user-provided input
}
