package migrate

import (
	"dg-server/core"
	"fmt"
	"io/ioutil"
	"strings"
)

func init() {
	loadTempleteFile()
}

func loadTempleteFile() {
	files, _ := ioutil.ReadDir("./templetes")
	for _, f := range files {
		fileName := f.Name()
		data, err := ioutil.ReadFile("./templetes/" + fileName)
		if err != nil {
			fmt.Println("Templete File reading error", err)
			continue
		}
		content := string(data)
		kv := strings.Split(fileName, ".")
		MigrationTempletes = append(MigrationTempletes, core.Templete{
			Name:     kv[0],
			DataBase: kv[1],
			Language: kv[2],
			Orm:      kv[3],
			Content:  content,
			Type:     "init",
		})
	}
}

// MigrationTempletes ...
var MigrationTempletes = []core.Templete{}
