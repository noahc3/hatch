package loaders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/noahc3/hatch-cli/src/log"
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
		log.Error("  Failed to parse Egg file: %s\n", err)
		return types.Egg{}, err
	}

	return egg, nil

}

func loadEggFromUrl(url string) (types.Egg, error) {
	log.Info("  Fetching from: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return types.Egg{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Egg{}, log.Error("  Failed to fetch Egg file from remote (%s)", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Egg{}, err
	}

	var egg types.Egg
	err = json.Unmarshal(data, &egg)
	if err != nil {
		log.Error("  Failed to parse Egg file: %s\n", err)
		return types.Egg{}, err
	}

	log.Info("  Fetch successful\n")

	return egg, nil
}

func loadEggFromGithub(tag string) (types.Egg, error) {
	pattern := `@github:([^\/]*)\/([^\/:]*)\/?([^:]*)?:(.*\.json)`
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(tag)
	if match == nil {
		return types.Egg{}, log.Error("  GitHub provider tag does not match the required format @github:owner/repo[/ref]:path/to/file.json\n")
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
}

func LoadEgg(tag string) (types.Egg, error) {

	log.Info("Loading Egg: %s\n", tag)

	if strings.HasPrefix(tag, "@") {
		provider := strings.Split(tag[1:], ":")[0]

		switch provider {
		case "github":
			log.Info("  Using GitHub provider\n")
			return loadEggFromGithub(tag)
		}

		return types.Egg{}, log.Error("  Unknown remote provider: @%s\n", provider)
	} else if strings.HasPrefix(tag, "http://") || strings.HasPrefix(tag, "https://") {
		log.Info("  Fetching from remote\n")

		return loadEggFromUrl(tag)
	} else {
		log.Info("  Loading from disk\n")

		_, err := os.Stat(tag)
		if os.IsNotExist(err) {
			return types.Egg{}, log.Error("  Failed to find Egg file: '%s'\n", tag)
		} else if err != nil {
			return types.Egg{}, log.Error("  Failed to load Egg file: '%s'\n", err)
		}

		return loadEggFromFile(tag)
	}
}
