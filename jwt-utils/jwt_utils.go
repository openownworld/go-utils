//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   jwt_utils.go
// @Description:

package jwt_utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	_ "github.com/gbrlsnchs/jwt/v3/jwtutil"
	"io/ioutil"
	_ "strconv"
	"strings"
	"time"
)

type JwtUtils struct {
	Alg           int   //算法类型
	ExpTimeSecond int64 //失效时间，单位秒 如：2小时后
	HsSecret      string
	rsaPrivateKey string
	rsaPublicKey  string
	rsaPriv       *rsa.PrivateKey
	rsaPub        *rsa.PublicKey
}

var (
	AlgorithmHS256 = 1
	AlgorithmRS256 = 2
)

func NewJwtUtils() *JwtUtils {
	obj := new(JwtUtils)
	obj.Alg = AlgorithmHS256
	obj.ExpTimeSecond = 60 * 60
	return obj
}

func genRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return priv, &priv.PublicKey
}

func ReadFile(fileName string) ([]byte, error) {
	rsaKey, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}

// GetPubKey 获取RSA公钥长度
// 参数：public
// 返回：成功则返回 RSA 公钥长度，失败返回 error 错误信息
func GetPubKey(pubKey []byte) (*rsa.PublicKey, error) {
	if pubKey == nil {
		return nil, errors.New("input arguments error")
	}
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public rsaKey error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	fmt.Println("pub rsa bit ", pub.N.BitLen())
	//return pub.N.BitLen(), nil
	return pub, nil
}

// GetPriKey 这里只识别RSA PKCS#1传统的，不能识别 PKCS8格式私钥
// 获取RSA私钥长度
// PriKey
// 成功返回 RSA 私钥长度，失败返回error
func GetPriKey(priKey []byte) (*rsa.PrivateKey, error) {
	if priKey == nil {
		return nil, errors.New("input arguments error")
	}
	//
	//re, _ := regexp.Compile("\\-*BEGIN.*KEY\\-*")
	//re2, _ := regexp.Compile("\\-*END.*KEY\\-*")
	//s := string(priKey)
	//pemStrTmp := re.ReplaceAllString(s, "");
	//pemStr := re2.ReplaceAllString(pemStrTmp, "");
	//pemStr = strings.Replace(pemStr, "\n", "", 1)
	//
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private rsaKey error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("priv rsa bit", priv.N.BitLen())
	//return priv.N.BitLen(), nil
	return priv, nil
}

func (j *JwtUtils) SetRsaPubKey(key []byte) error {
	rsa, err := GetPubKey(key)
	if err != nil {
		return err
	}
	j.rsaPub = rsa
	return nil
}

func (j *JwtUtils) SetRsaPrivKey(key []byte) error {
	rsa, err := GetPriKey(key)
	if err != nil {
		return err
	}
	j.rsaPriv = rsa
	return nil
}

// CreateJwt {"alg":"HS256","typ":"JWT"}.{"exp":1570772471,"iat":1570765271,"jti":"1200","uid":"100000001"}
func (j *JwtUtils) CreateJwt(jwtID string, claims map[string]interface{}) (string, error) {
	/*
		type CustomPayload struct {
			jwt.Payload
			Foo string `json:"foo,omitempty"`
		}
		now := time.Now()
		pl := CustomPayload{
			Payload: jwt.Payload{
				//Issuer:         "gbrlsnchs",
				//Subject:        "someone",
				//Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
				ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
				//NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
				IssuedAt: jwt.NumericDate(now),
				JWTID:    jwtID,
			},
			Foo: "foo1",
		}
	*/
	//
	alg := j.Alg
	var algKey jwt.Algorithm
	if alg == 1 {
		algKey = jwt.NewHS256([]byte(j.HsSecret))
	} else if alg == 2 {
		algKey = jwt.NewRS256(jwt.RSAPrivateKey(j.rsaPriv))
	} else {
		return "", errors.New("无效的算法")
	}

	payloadClaims := make(map[string]interface{})
	//payload_claims["exp"] = strconv.FormatInt(time.Now().Add(time.Hour*time.Duration(2)).Unix(), 10)
	//payload_claims["iat"] = strconv.FormatInt(time.Now().Unix(), 10)
	//payload_claims["exp"] = time.Now().Add(time.Hour * time.Duration(2)).Unix()
	payloadClaims["exp"] = time.Now().Add(time.Second * time.Duration(j.ExpTimeSecond)).Unix()
	payloadClaims["iat"] = time.Now().Unix() //unix时间戳是从1970年1月1日（UTC/GMT的午夜）开始所经过的秒数
	payloadClaims["jti"] = jwtID
	for key, val := range claims {
		payloadClaims[key] = val
	}
	token, err := jwt.Sign(payloadClaims, algKey)
	if err != nil {
		return "", err
	}
	//encodeString := base64.StdEncoding.EncodeToString(token)
	return string(token), nil
}

// VerifyParseTokenAlg 验证token并解析payload
func (j *JwtUtils) VerifyParseTokenAlg(alg int, token string) (map[string]interface{}, error) {
	payloadClaims := make(map[string]interface{})
	//
	var algKey jwt.Algorithm
	if alg == 1 {
		algKey = jwt.NewHS256([]byte(j.HsSecret))
	} else if alg == 2 {
		algKey = jwt.NewRS256(jwt.RSAPublicKey(j.rsaPub))
	} else {
		return payloadClaims, errors.New("无效的算法")
	}
	_, err := jwt.Verify([]byte(token), algKey, &payloadClaims)
	if err != nil {
		return payloadClaims, err
	}
	return payloadClaims, nil
}

// ParseToken token解析payload
func (j *JwtUtils) ParseToken(token string) (map[string]interface{}, error) {
	payloadClaims := make(map[string]interface{})
	header := make(map[string]interface{})
	ss := strings.Split(token, ".")
	if len(ss) > 1 {
		//encodeString := base64.StdEncoding.EncodeToString(token)
		decodeString, _ := base64.RawURLEncoding.DecodeString(ss[0])
		err := json.Unmarshal([]byte(decodeString), &header)
		if err != nil {
			return payloadClaims, err
		}
		//
		decodeString, _ = base64.RawURLEncoding.DecodeString(ss[1])
		err = json.Unmarshal([]byte(decodeString), &payloadClaims)
		if err != nil {
			return payloadClaims, err
		}
		for k, v := range header {
			payloadClaims[k] = v
		}
		return payloadClaims, nil
	} else {
		return payloadClaims, fmt.Errorf("jwt string err")
	}
}

// VerifyParseToken 验证token并解析payload
func (j *JwtUtils) VerifyParseToken(token string) (map[string]interface{}, error) {
	payloadClaims := make(map[string]interface{})
	//
	var algKey jwt.Algorithm
	var alg = ""
	header := make(map[string]interface{})
	ss := strings.Split(token, ".")
	if len(ss) > 1 {
		//encodeString := base64.StdEncoding.EncodeToString(token)
		decodeString, _ := base64.RawURLEncoding.DecodeString(ss[0])
		err := json.Unmarshal([]byte(decodeString), &header)
		if err != nil {
			return payloadClaims, err
		}
		alg = header["alg"].(string) //强制转换
	}
	if alg == "HS256" {
		algKey = jwt.NewHS256([]byte(j.HsSecret))
	} else if alg == "RS256" {
		algKey = jwt.NewRS256(jwt.RSAPublicKey(j.rsaPub))
	} else {
		return payloadClaims, errors.New("无效的算法")
	}
	_, err := jwt.Verify([]byte(token), algKey, &payloadClaims)
	if err != nil {
		return payloadClaims, err
	}
	return payloadClaims, nil
}

func (j *JwtUtils) VerifyToken(token string) error {
	//
	var algKey jwt.Algorithm
	var alg = ""
	header := make(map[string]interface{})
	ss := strings.Split(token, ".")
	if len(ss) > 1 {
		decodeString, _ := base64.RawURLEncoding.DecodeString(ss[0])
		err := json.Unmarshal([]byte(decodeString), &header)
		if err != nil {
			return err
		}
		alg = header["alg"].(string) //强制转换
	}
	if alg == "HS256" {
		algKey = jwt.NewHS256([]byte(j.HsSecret))
	} else if alg == "RS256" {
		algKey = jwt.NewRS256(jwt.RSAPublicKey(j.rsaPub))
	} else {
		return errors.New("无效的算法")
	}
	var (
		now = time.Now()
		//aud = jwt.Audience{"https://golang.org"}
		// Validate claims "iat", "exp" and "aud".
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		//audValidator = jwt.AudienceValidator(aud)
		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		pl              jwt.Payload
		validatePayload = jwt.ValidatePayload(&pl, iatValidator, expValidator)
	)
	_, err := jwt.Verify([]byte(token), algKey, &pl, validatePayload)
	if err != nil {
		return err
	}
	return nil
}
