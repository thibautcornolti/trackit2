package parsator

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

func transformArrayToMap(in *interface{}) {
	typpedIn := (*in).([]interface{})
	newMap := make(map[string]interface{})
	for _, field := range typpedIn {
		typpedMap := field.(map[string]interface{})
		for name, fieldMap := range typpedMap {
			newMap[name] = fieldMap
		}
	}
	*in = newMap
}

func recurMap(in map[string]interface{}) {
	for name, mapField := range in {
		if name[0] != '&' && name[0] != '*' {
			continue
		}
		typpedMapField := mapField.(map[string]interface{})
		splittedName := strings.Split(name[1:], ".")
		aggSect := splittedName[0]
		typpedSectMapField := typpedMapField[aggSect]
		if name[0] == '&' {
			recurTab(&typpedSectMapField)
		}
		in[in["key"].(string)] = typpedSectMapField
		delete(in, "key")
		delete(in, name)
	}
}

func recurTab(in *interface{}) {
	typpedIn := (*in).([]interface{})
	transformable := true
	for _, field := range typpedIn {
		typpedField := field.(map[string]interface{})
		delete(typpedField, "doc_count")
		recurMap(typpedField)
		if len(typpedField) > 1 {
			transformable = false
		}
	}
	if transformable {
		transformArrayToMap(in)
	}
}

func Test(esResult *elastic.SearchResult) map[string]interface{} {
	res := make(map[string]interface{})
	resJson, _ := json.MarshalIndent(*esResult, "", "  ")
	fmt.Printf("%v\n", string(resJson))
	for name, agg := range esResult.Aggregations {
		var t map[string]interface{}
		err := json.Unmarshal(*agg, &t)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		splittedName := strings.Split(name[1:], ".")
		realName := splittedName[len(splittedName)-1]
		aggSect := splittedName[0]
		sectT := t[aggSect]
		if name[0] == '&' || name[0] == '*' {
			recurTab(&sectT)
			res[realName] = sectT
		} else {
			res[name] = t
		}
	}
	return res
}
