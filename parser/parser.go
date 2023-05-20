package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Parse() (map[string]map[string]interface{}, error) {
	file, err := os.Open("dbList.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteResult, _ := ioutil.ReadAll(file)

	var res map[string]map[string]interface{}
	json.Unmarshal([]byte(byteResult), &res)

	return res, nil
}
