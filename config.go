package main

import (
	"bufio"
	"os"
	"strings"
)

//readConfigFile takes a path to a configuration file and returns a map of configuration parameters
func readConfigFile(path string) map[string]string {

	configMap := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "//") {
			fields := strings.SplitN(scanner.Text(), "=", 2)

			configMap[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}
	}

	return configMap
}
