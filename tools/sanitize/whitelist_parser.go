package sanitize

import(
	"os"
	"encoding/json"
)

// Load a new whitelist from a JSON file
func WhitelistFromFile(filepath string) (*Whitelist, error) {
	bytes, err := readFileToBytes(filepath)
	if err != nil {
		return nil, err
	}

	whitelist, err := NewWhitelist(bytes)
	return whitelist, nil
}

// helper function to read entirety of provided file into byte slice
func readFileToBytes(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	// prepare byte slice to read json file into
	fileInfo, err := f.Stat()
	bytes := make([]byte, fileInfo.Size())

	_, err = f.Read(bytes)
	return bytes, err
}

// Create a new whitelist from JSON configuration
func NewWhitelist(jsonData []byte) (*Whitelist, error) {
	// unmarshal json file into contract-free interface
	configuration := &Whitelist{}
	err := json.Unmarshal(jsonData, configuration)

	return configuration, err
}