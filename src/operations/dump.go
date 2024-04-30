package operations

import (
	"encoding/json"
	"fmt"

	"github.com/noahc3/hatch-cli/src/types"
)

func Print(egg *types.Egg) error {
	eggJson, err := json.MarshalIndent(egg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", eggJson)
	return nil
}
