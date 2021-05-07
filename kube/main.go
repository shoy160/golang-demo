package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"

	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alimt20181012.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("mt.cn-hangzhou.aliyuncs.com")
	_result = &alimt20181012.Client{}
	_result, _err = alimt20181012.NewClient(config)
	return _result, _err
}

func translate(word string) string {
	accessKeyId := ""
	accessKeySecret := ""
	client, _err := CreateClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if _err != nil {
		log.Fatal(_err)
		return word
	}

	translateGeneralRequest := &alimt20181012.TranslateGeneralRequest{}
	translateGeneralRequest.FormatType = tea.String("text")
	translateGeneralRequest.SourceLanguage = tea.String("en")
	translateGeneralRequest.TargetLanguage = tea.String("zh")
	translateGeneralRequest.SourceText = tea.String(word)
	// 复制代码运行请自行打印 API 的返回值
	res, _err := client.TranslateGeneral(translateGeneralRequest)
	if _err != nil {
		log.Fatal(_err)
		return word
	}
	return *res.Body.Data.Translated
}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func foreach_json(json map[string]interface{}) {
	keys := []string{"description", "summary"}
	for key, val := range json {
		switch val.(type) {
		case map[string]interface{}:
			foreach_json(val.(map[string]interface{}))

		case []interface{}:
			foreach_array(val.([]interface{}))
		default:
			if in(key, keys) {
				// println(key, val.(string))
				log.Println(val.(string))
				json[key] = translate(val.(string))
				log.Println(json[key])
			}
		}
	}
}

func foreach_array(anArray []interface{}) {
	for _, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			foreach_json(val.(map[string]interface{}))
		case []interface{}:
			foreach_array(val.([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
		}
	}
}

func load_json(path string) map[string]interface{} {
	res := map[string]interface{}{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return res
	}

	err = json.Unmarshal(f, &res)
	if err != nil {
		log.Fatal(err)
		return res
	}
	return res
}

func main() {
	swagger := load_json("C:\\Users\\shay\\Desktop\\kubesphere\\swagger.json")
	// log.Println(swagger)
	foreach_json(swagger)
	data, err := json.Marshal(swagger)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("C:\\Users\\shay\\Desktop\\kubesphere\\swagger_t.json", data, 0755)
	// word := translate("query conditions,connect multiple conditions with commas, equal symbol for exact query, wave symbol for fuzzy query e.g. name~a")
	// log.Println(word)
}
