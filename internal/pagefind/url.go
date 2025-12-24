package pagefind

import "fmt"

const pagefindVersion = "v1.4.0"

func downloadURL() (string, bool, error) {
	key, err := platformKey()
	if err != nil {
		return "", false, err
	}

	if key == "x86_64-pc-windows-msvc" {
		return fmt.Sprintf(
			"https://github.com/Pagefind/pagefind/releases/download/%s/pagefind-%s-%s.zip",
			pagefindVersion,
			pagefindVersion,
			key,
		), true, nil
	}

	return fmt.Sprintf(
		"https://github.com/Pagefind/pagefind/releases/download/%s/pagefind-%s-%s.tar.gz",
		pagefindVersion,
		pagefindVersion,
		key,
	), false, nil
}
