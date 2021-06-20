package utils

import "fmt"

func FormatAppInfoHeader(name, version, build string) string {
	return fmt.Sprintf("%s/%s (build %s)", name, version, build)
}
