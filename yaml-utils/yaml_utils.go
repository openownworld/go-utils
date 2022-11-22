package yaml_utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// NewYamlConfig 读取配置文件，填充yamlConfig 结构体，传入地址
func NewYamlConfig(confFileName string, yamlConfig interface{}) error {
	f, err := os.OpenFile(confFileName, os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("config file %s, error %s ", confFileName, err.Error())
	}
	defer f.Close()
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read config file %s, error %s ", confFileName, err.Error())
	}
	//读配置文件
	err = yaml.Unmarshal([]byte(contentByte), yamlConfig)
	if err != nil {
		return fmt.Errorf("unmarshal config file %s, error %s ", confFileName, err.Error())
	}
	return nil
}

func NewYamlConfigByIO(iniReader io.Reader, yamlConfig interface{}) error {
	contentByte, err := ioutil.ReadAll(iniReader)
	if err != nil {
		return fmt.Errorf("read config stream, error %s ", err.Error())
	}
	//读配置文件
	err = yaml.Unmarshal([]byte(contentByte), yamlConfig)
	if err != nil {
		return fmt.Errorf("unmarshal config stream, error %s ", err.Error())
	}
	return nil
}
