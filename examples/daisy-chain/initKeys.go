package main

import (
	"os"
	"path/filepath"
)

func initKeys(keys map[string]string) string {
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "tf", "inception", "config.json")
	keys["pv"] = "b773bde3-e73b-4914-8fac-3513ca76a596"
	keys["pvc"] = "e614bfad-436e-4e26-b6b5-41384f2260e6"
	keys["pods"] = "86f96730-44c7-4942-ba59-9cc711143ffa"
	keys["svc"] = "58a60855-ca76-45c8-a8cf-ad2b03362db9"
	keys["ds"] = "4d415f93-0a19-4037-839e-00bd7b049eae"
	keys["jobs"] = "567605cf-3ae7-41ef-8a49-31bb1f6bb3cc"

	return configFile
}
