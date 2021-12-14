//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   yaml_test.go
// @Description:

package yaml_utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type YamlConfig struct {
	//配置文件要通过tag来指定配置文件中的名称
	Local struct {
		WsIPAddr string `yaml:"wsIPAddr"`
	}

	MySQL struct {
		Url string `yaml:"url"`
	} `yaml:"mysql" json:"mysql"`
}

func (y *YamlConfig) String() string {
	//jsonBytes, _ := json.Marshal(c)
	jsonBytes, _ := json.MarshalIndent(y, "", "    ") //格式化json
	return strings.Replace(string(jsonBytes), "\\u0026", "&", -1)
	//return string(jsonBytes) //json.Marshal 默认 escapeHtml 为true,会转义 <、>、&
}

func TestNewYamlConfig(t *testing.T) {
	var config YamlConfig
	NewYamlConfig("./test.yaml", &config)
	fmt.Println(config)
}
