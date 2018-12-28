package util

import (
	"bobthebuilder/logging"
	"io/ioutil"
	"path"
	"strings"
)

func GetFilenameListInFolder(folder, suffix string) ([]string, error) {
	output := []string{}
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		logging.Error("file-util", err)
		return nil, err
	}

	for _, file := range files {
		if (!file.IsDir()) && strings.HasSuffix(file.Name(), suffix) {

			p := file.Name()
			if !path.IsAbs(p) {
				p = path.Join(folder, file.Name())
			}
			output = append(output, p)
		}
	}
	return output, nil
}
