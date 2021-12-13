package zaplog

import (
	"encoding/json"
	"fmt"
)

// LogInfo log打印标记信息
type LogKeyInfo struct {
	MapInfo map[string]interface{}
}

// Clone 深拷贝MAP
func (s *LogKeyInfo) Clone() *LogKeyInfo {
	logKeyInfo := LogKeyInfo{MapInfo: make(map[string]interface{})}
	for k, v := range s.MapInfo {
		logKeyInfo.MapInfo[k] = v
	}
	return &logKeyInfo
}

func (s *LogKeyInfo) ToString() string {
	str := "["
	for k, v := range s.MapInfo {
		str += fmt.Sprintf("%s=%s,", k, v)
	}
	str += "]"
	return str
}

func (s *LogKeyInfo) ToJsonString() string {
	jsonBytes, _ := json.Marshal(s.MapInfo)
	return string(jsonBytes)
}

func (s *LogKeyInfo) Add(k string, v interface{}) {
	s.MapInfo[k] = v
}

func (s *LogKeyInfo) Del(k string) {
	delete(s.MapInfo, k)
}

func (s *LogKeyInfo) Clear() {
	m := make(map[string]interface{})
	s.MapInfo = m
}
