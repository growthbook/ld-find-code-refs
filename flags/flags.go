package flags

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/growthbook/gb-find-code-refs/internal/log"
	"github.com/growthbook/gb-find-code-refs/options"
)

const (
	minFlagKeyLen = 3 // Minimum flag key length helps reduce the number of false positives
)

func GetFlagKeys(opts options.Options) []string {
	flags, err := getFlags(opts.FlagsPath)
	if err != nil {
		log.Error.Fatal(fmt.Errorf("could not parse flag keys: %w", err))
	}

	filteredFlags, omittedFlags := filterShortFlagKeys(flags)
	if len(filteredFlags) == 0 {
		log.Info.Printf("no flag keys longer than the minimum flag key length (%v) were found, exiting early",
			minFlagKeyLen)
		os.Exit(0)
	} else if len(omittedFlags) > 0 {
		log.Warning.Printf("omitting %d flags with keys less than minimum (%d)", len(omittedFlags), minFlagKeyLen)
	}
	return filteredFlags
}

// Very short flag keys lead to many false positives when searching in code,
// so we filter them out.
func filterShortFlagKeys(flags []string) (filtered []string, omitted []string) {
	filteredFlags := []string{}
	omittedFlags := []string{}
	for _, flag := range flags {
		if len(flag) >= minFlagKeyLen {
			filteredFlags = append(filteredFlags, flag)
		} else {
			omittedFlags = append(omittedFlags, flag)
		}
	}
	return filteredFlags, omittedFlags
}

func getFlags(flagsPath string) ([]string, error) {
	jsonFile, err := os.Open(flagsPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var flags []string

	err = json.NewDecoder(jsonFile).Decode(&flags)
	if err != nil {
		return nil, err
	}

	return flags, nil
}
