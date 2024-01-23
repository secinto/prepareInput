package prepare

import (
	"bufio"
	"net/url"
	"os"
	"sort"
	"strings"
)

func CreateDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o770)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToTextFileInProject(filename string, data string) {
	writeFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	dataWriter := bufio.NewWriter(writeFile)

	if err != nil {
		log.Error(err)
	}
	dataWriter.WriteString(data)
	dataWriter.Flush()
	writeFile.Close()
}

func ConvertStringArrayToString(stringArray []string, separator string) string {
	sort.Strings(stringArray)
	justString := strings.Join(stringArray, separator)
	return justString
}

func AppendIfMissing(slice []string, key string) []string {
	for _, element := range slice {
		if element == key {
			return slice
		}
	}
	return append(slice, key)
}

func ExistsInArray(slice []string, key string) bool {
	for _, element := range slice {
		if element == key {
			return true
		}
	}
	return false
}

func GetHost(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return str
	}
	return u.Scheme + "://" + u.Host
}

func ConvertJSONLtoJSON(input string) string {

	var data []byte
	data = append(data, '[')

	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	isFirst := true
	for _, line := range lines {
		if !isFirst && strings.TrimSpace(line) != "" {
			data = append(data, ',')
			data = append(data, '\n')
		}
		if strings.TrimSpace(line) != "" {
			data = append(data, line...)
		}
		isFirst = false
	}
	data = append(data, ']')
	return string(data)
}
