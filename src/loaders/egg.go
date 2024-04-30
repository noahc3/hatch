package loaders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/noahc3/hatch-cli/src/types"
)

func loadEggFromFile(path string) (types.Egg, error) {
	file, err := os.Open(path)
	if err != nil {
		return types.Egg{}, err
	}

	defer file.Close()

	var egg types.Egg
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&egg)
	if err != nil {
		return types.Egg{}, err
	}

	return egg, nil

}

func loadEggFromUrl(url string) (types.Egg, error) {
	resp, err := http.Get(url)
	if err != nil {
		return types.Egg{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Egg{}, fmt.Errorf("failed to fetch Egg file: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Egg{}, err
	}

	var egg types.Egg
	err = json.Unmarshal(data, &egg)
	if err != nil {
		return types.Egg{}, err
	}

	return egg, nil
}

func LoadEgg(tag string) (types.Egg, error) {
	if strings.HasPrefix(tag, "@") {
		pattern := `@github:([^\/]*)\/([^\/:]*)\/?([^:]*)?:(.*\.json)`
		re := regexp.MustCompile(pattern)

		match := re.FindStringSubmatch(tag)
		if match == nil {
			return types.Egg{}, fmt.Errorf("invalid GitHub tag: '%s'", tag)
		}

		owner := match[1]
		repo := match[2]
		ref := match[3]
		path := match[4]

		if ref == "" {
			ref = "master"
		}

		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)

		return loadEggFromUrl(url)
	} else if strings.HasPrefix(tag, "http://") || strings.HasPrefix(tag, "https://") {
		return loadEggFromUrl(tag)
	} else {
		_, err := os.Stat(tag)
		if os.IsNotExist(err) {
			return types.Egg{}, fmt.Errorf("failed to find Egg file: '%s'", tag)
		} else if err != nil {
			return types.Egg{}, err
		}

		return loadEggFromFile(tag)
	}
}
