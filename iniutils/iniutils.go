package iniutils

import "github.com/go-ini/ini"

//https://www.jianshu.com/p/c1dd0f88abe1

type IniParser struct {
	confReader *ini.File // config reader
}

type IniParserError struct {
	errorInfo string
}

func (e *IniParserError) Error() string { return e.errorInfo }

func (this *IniParser) Load(configFileName string) error {
	conf, err := ini.Load(configFileName)
	if err != nil {
		this.confReader = nil
		return err
	}
	this.confReader = conf
	return nil
}

func (this *IniParser) GetString(section string, key string) string {
	if this.confReader == nil {
		return ""
	}

	s := this.confReader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

func (this *IniParser) GetInt32(section string, key string) int32 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int()

	return int32(valueInt)
}

func (this *IniParser) GetUint32(section string, key string) uint32 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint()

	return uint32(valueInt)
}

func (this *IniParser) GetInt64(section string, key string) int64 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int64()
	return valueInt
}

func (this *IniParser) GetUint64(section string, key string) uint64 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint64()
	return valueInt
}

func (this *IniParser) GetFloat32(section string, key string) float32 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return float32(valueFloat)
}

func (this *IniParser) GetFloat64(section string, key string) float64 {
	if this.confReader == nil {
		return 0
	}

	s := this.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return valueFloat
}

func (this *IniParser) GetFileHandle() *ini.File {
	if this.confReader == nil {
		return nil
	}

	return this.confReader
}

//Section为空时，默认值为DEFAULT
//DEFAULT.default=kv
//server.ip=0.0.0.0
func (this *IniParser) GetValueMap() map[string]interface{} {
	if this.confReader == nil {
		return nil
	}

	s := this.confReader.Sections()
	if s == nil {
		return nil
	}

	valueMap := make(map[string]interface{})
	for i := range s {
		sn := s[i].Name()
		k := s[i].Keys()
		for j := range k {
			kn := k[j].Name()
			valueMap[sn+"."+kn] = k[j].Value()
		}
	}

	return valueMap
}
