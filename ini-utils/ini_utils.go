//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   ini_utils.go
// @Description:

package ini_utils

import (
	"fmt"
	"io"

	"github.com/go-ini/ini"
)

type Parser struct {
	reader *ini.File // config reader
}

type ParserError struct {
	errorInfo string
}

func (p *ParserError) Error() string { return p.errorInfo }

func NewIniParserByIO(reader io.Reader) (*Parser, error) {
	p := Parser{}
	conf, err := ini.Load(reader)
	if err != nil {
		p.reader = nil
		return nil, fmt.Errorf("load config stream failed: %v", err)
	}
	p.reader = conf
	return &p, nil
}

func NewIniParser(fileName string) (*Parser, error) {
	p := Parser{}
	conf, err := ini.Load(fileName)
	if err != nil {
		p.reader = nil
		return nil, fmt.Errorf("load config file %s, error %s", fileName, err.Error())
	}
	p.reader = conf
	return &p, nil
}

func (p *Parser) GetString(section string, key string) string {
	if p.reader == nil {
		return ""
	}

	s := p.reader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

func (p *Parser) GetInt32(section string, key string) int32 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int()

	return int32(valueInt)
}

func (p *Parser) GetUint32(section string, key string) uint32 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint()

	return uint32(valueInt)
}

func (p *Parser) GetInt64(section string, key string) int64 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int64()
	return valueInt
}

func (p *Parser) GetUint64(section string, key string) uint64 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint64()
	return valueInt
}

func (p *Parser) GetFloat32(section string, key string) float32 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return float32(valueFloat)
}

func (p *Parser) GetFloat64(section string, key string) float64 {
	if p.reader == nil {
		return 0
	}

	s := p.reader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return valueFloat
}

func (p *Parser) GetFileHandle() *ini.File {
	if p.reader == nil {
		return nil
	}

	return p.reader
}

// GetValueMap ini-utils to map
// Section为空时，默认值为DEFAULT
// DEFAULT.default=kv
// server.ip=0.0.0.0
func (p *Parser) GetValueMap() map[string]interface{} {
	if p.reader == nil {
		return nil
	}

	s := p.reader.Sections()
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
