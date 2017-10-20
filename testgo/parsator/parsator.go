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

func cleanField(name string, sectMapField interface{}, in map[string]interface{}) {
	if _, ok := in["key"]; ok {
		if typpedSectMapField, ok := sectMapField.(map[string]interface{}); ok && name[0] == '*' {
			sectMapField := typpedSectMapField["buckets"]
			transformArrayToMap(&sectMapField)
			in[in["key"].(string)] = sectMapField
		} else {
			in[in["key"].(string)] = sectMapField
		}
		delete(in, "key")
	}
	delete(in, name)
}

func recurMap(in map[string]interface{}) {
	for name, mapField := range in {
		if name[0] != '&' && name[0] != '*' {
			continue
		}
		typpedMapField := mapField.(map[string]interface{})
		splittedName := strings.Split(name[1:], ".")
		aggSect := strings.Replace(splittedName[0], ",", ".", -1)
		typpedSectMapField := typpedMapField[aggSect]
		if name[0] == '&' {
			recurTab(&typpedSectMapField)
		} else if name[0] == '*' {
			recurMap(typpedMapField)
		}
		cleanField(name, typpedSectMapField, in)
	}
}

func recurTab(in *interface{}) {
	typpedIn := (*in).([]interface{})
	for _, field := range typpedIn {
		typpedField := field.(map[string]interface{})
		delete(typpedField, "doc_count")
		recurMap(typpedField)
	}
	transformArrayToMap(in)
}

// GetParsedElasticSearchResult ...
func GetParsedElasticSearchResult(esResult *elastic.SearchResult) map[string]interface{} {
	res := make(map[string]interface{})
	resJSON, _ := json.MarshalIndent(*esResult, "", "  ")
	fmt.Printf("%v\n", string(resJSON))
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
