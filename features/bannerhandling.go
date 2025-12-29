package asciiart

import (
	"os"
)

func ReadBanner(banner string) string {
	data, err := os.ReadFile("banner/" + banner + ".txt")
	CheckError(err)
	stringData := string(data)
	return stringData
}
