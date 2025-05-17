package rao

import (
	c_rand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math"
	"scene/internal/biz/log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"

	"github.com/mritd/chinaid"
)

const (
	FuncAssistantPrefix = "__"
	CustomerDefPrefix   = "@_"

	ChangeCase = "ChangeCase"
	UPPER      = "UPPER"
	LOWER      = "LOWER"

	TimeFormat      = "TimeFormat"
	RandomInt       = "RandomInt"
	RandomStr       = "RandomStr"
	RandomChooseStr = "RandomChooseStr"
	SubStr          = "SubStr"
	UUID            = "UUID"
	MD5             = "MD5"
	Base64Encode    = "Base64Encode"
	Base64Decode    = "Base64Decode"
	RSA             = "RSA"

	Calculate = "Calculate"

	IDCard      = "IDCard"
	ChineseName = "ChineseName"
	Mobile      = "Mobile"
)

var timeFormatMap = map[string]string{
	"YYYY-MM-DD HH:MM:SS":    "2006-01-02 15:04:05",
	"MM/DD/YYYY HH:MM:SS":    "01/02/2006 15:04:05",
	"DD-MM-YYYY HH:MM:SS":    "02-01-2006 15:04:05",
	"MM/DD/YYYY":             "01/02/2006",
	"DD-MM-YYYY":             "02-01-2006",
	"YYYY/MM/DD":             "2006/01/02",
	"YYYY年MM月DD日 HH时MM分SS秒":  "2006年01月02日 15时04分05秒",
	"YYYY年MM月DD日":            "2006年01月02日",
	"YYYY年M月D日 H时M分S秒":       "2006年1月2日 3时4分5秒",
	"HH:MM:SS":               "15:04:05",
	"HH:MM":                  "15:04",
	"HH时MM分SS秒":              "15时04分05秒",
	"HH时MM分":                 "15时04分",
	"HH时":                    "15时",
	"YYYY-MM-DD":             "2006-01-02",
	"YYYY/MM/DD HH:MM:SS":    "2006/01/02 15:04:05",
	"YYYY/MM/DD HH:MM":       "2006/01/02 15:04",
	"YYYY/MM/DD HH":          "2006/01/02 15",
	"YYYY/MM/DD HH:MM:SS AM": "2006/01/02 15:04:05 PM",
	"YYYY/MM/DD HH:MM AM":    "2006/01/02 15:04 PM",
	"YYYY/MM/DD HH AM":       "2006/01/02 15 PM",
	"YYYY/MM/DD HH:MM:SS PM": "2006/01/02 15:04:05 PM",
	"YYYY/MM/DD HH:MM PM":    "2006/01/02 15:04 PM",
	"YYYY/MM/DD HH PM":       "2006/01/02 15 PM",
}

// replaceVariables
// key: ${xxxx}
func replaceVariables(key *string, variablePool *sync.Map) {
	var sb strings.Builder // 使用strings.Builder提高字符串拼接性能
	s := *key
	startIndex := 0
	for {
		openIndex := strings.Index(s[startIndex:], "${")
		if openIndex == -1 { // 没有找到"${"，则剩余部分无需处理，直接添加到结果中
			sb.WriteString(s[startIndex:])
			break
		}

		closeIndex := strings.Index(s[startIndex+openIndex+2:], "}")
		if closeIndex == -1 {
			// 没有找到"}"
			sb.WriteString(s[startIndex:])
			break
		}
		closeIndex += startIndex + openIndex + 2
		// 获取变量名
		variableName := s[startIndex+openIndex+2 : closeIndex]
		newVal := ""
		if strings.HasPrefix(variableName, FuncAssistantPrefix) {
			// 走函数助手逻辑
			newVal = FunctionAssistant(variableName, variablePool)
		} else {
			// 走正常变量替换逻辑
			val, ok := variablePool.Load(variableName)
			if !ok {
				newVal = "${" + variableName + "}"
			} else {
				newVal = val.(string)
			}
		}
		sb.WriteString(s[startIndex : startIndex+openIndex])
		sb.WriteString(newVal)
		// 更新startIndex为当前匹配结束的位置之后
		startIndex = closeIndex + 1
	}
	*key = sb.String()
}

func FunctionAssistant(key string, variablePool *sync.Map) (newVal string) {
	if !strings.HasPrefix(key, FuncAssistantPrefix) {
		return key
	}

	leftIndex := strings.Index(key, "(")
	if leftIndex == -1 {
		return key
	}
	funcType := key[2:leftIndex]
	params := strings.Split(key[leftIndex+1:len(key)-1], ",")
	// 遍历params，去除所有空格
	for i, param := range params {
		params[i] = strings.TrimSpace(param)
	}
	newVal = key
	// 如需另存为新变量，该变量在params内的索引
	saveKeyIndex := 0
	var err error
	switch funcType {
	case ChangeCase:
		switch params[1] {
		case UPPER:
			newVal = strings.ToUpper(params[0])
		case LOWER:
			newVal = strings.ToLower(params[0])
		default:
			newVal = params[0]
		}
		saveKeyIndex = 2
	case TimeFormat:
		format := params[0]
		cstSh, _ := time.LoadLocation(time.Local.String())
		if format == "" {
			// 获取毫秒时间戳
			newVal = strconv.Itoa(int(time.Now().In(cstSh).UnixMilli()))
		} else if format == "/1000" {
			newVal = strconv.Itoa(int(time.Now().In(cstSh).Unix()))
		} else {
			newVal = time.Now().In(cstSh).Format(timeFormatMap[format])
		}
		saveKeyIndex = 1
	case RandomInt:
		min := math.MinInt / 2
		max := math.MaxInt / 2
		if params[0] != "" {
			min, err = strconv.Atoi(params[0])
			if err != nil {
				return
			}
		}
		if params[1] != "" {
			max, err = strconv.Atoi(params[1])
			if err != nil {
				return
			}
		}
		newVal = strconv.Itoa(random.RandInt(min, max))
		saveKeyIndex = 2
	case RandomStr:
		length := 1
		if params[0] != "" {
			length, _ = strconv.Atoi(params[0])
		}
		randomStrPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		if len(params) > 1 && params[1] != "" {
			randomStrPool = params[1]
		}
		list := strings.Split(randomStrPool, "")
		str := ""
		// TODO: 优化为byte替换
		//for i := range b {
		//	b[i] = s[rand.Int63()%int64(len(s))]
		//}
		for i := 0; i < length; i++ {
			index := random.RandInt(0, len(list))
			str += list[index]
		}
		newVal = str
		saveKeyIndex = 2
	case RandomChooseStr:
		list := strings.Split(params[0], "|")
		newVal = list[random.RandInt(0, len(list))]
		saveKeyIndex = 1
	case SubStr:
		if len(params) < 3 || params[0] == "" {
			return
		}
		left, right := 0, len(params[0])
		if params[1] != "" {
			left, err = strconv.Atoi(params[1])
			if err != nil {
				return
			}
		}
		if params[2] != "" {
			le, err := strconv.Atoi(params[2])
			if err != nil {
				return
			}
			if right >= le {
				right = le
			}
		}
		newVal = params[0][left:right]
		saveKeyIndex = 3
	case UUID:
		newVal, _ = random.UUIdV4()
		saveKeyIndex = 0
	case MD5:
		getCustomerDefineVal(params[0], variablePool)
		newVal = cryptor.Md5String(params[0])
		saveKeyIndex = 1
	case Base64Encode:
		newVal = cryptor.Base64StdEncode(params[0])
		saveKeyIndex = 1
	case Base64Decode:
		newVal = cryptor.Base64StdDecode(params[0])
		saveKeyIndex = 1
	case RSA:
		publicKey := getCustomerDefineVal(params[0], variablePool)
		block, _ := pem.Decode([]byte(publicKey))
		if block == nil {
			log.Logger.Error("model.func_assistant.FunctionAssistant.pemDecode err")
			break
		}
		pubKey, parseErr := x509.ParsePKCS1PublicKey(block.Bytes)
		if parseErr != nil {
			log.Logger.Error("model.func_assistant.FunctionAssistant.x509ParseKey err: ", parseErr)
			break
		}
		encryptedData, enErr := rsa.EncryptPKCS1v15(c_rand.Reader, pubKey, []byte(getCustomerDefineVal(params[1], variablePool)))
		if enErr != nil {
			log.Logger.Error("model.func_assistant.FunctionAssistant.enErr err: ", enErr)
			break
		}
		newVal = base64.StdEncoding.EncodeToString(encryptedData)
		saveKeyIndex = 2
	case Calculate:
		if len(params) < 3 || params[0] == "" || params[1] == "" || params[2] == "" {
			return
		}
		obj1, err := strconv.Atoi(params[1])
		if err != nil {
			return
		}
		obj2, err := strconv.Atoi(params[2])
		if err != nil {
			return
		}
		switch params[0] {
		case "+":
			newVal = strconv.Itoa(obj1 + obj2)
		case "-":
			newVal = strconv.Itoa(obj1 - obj2)
		case "*":
			newVal = strconv.Itoa(obj1 * obj2)
		case "/":
			newVal = strconv.Itoa(obj1 / obj2)
		}
		saveKeyIndex = 3
	case IDCard:
		newVal = chinaid.IDNo()
		saveKeyIndex = 0
	case ChineseName:
		newVal = chinaid.Name()
		saveKeyIndex = 0
	case Mobile:
		newVal = chinaid.Mobile()
		saveKeyIndex = 0
	}
	if len(params) == saveKeyIndex+1 && params[saveKeyIndex] != "" {
		variablePool.Store(params[saveKeyIndex], newVal)
	}
	return
}

// 函数入参内支持@_XX声明，意味当前key为变量
func getCustomerDefineVal(key string, variablePool *sync.Map) (val string) {
	if !strings.HasPrefix(key, CustomerDefPrefix) {
		return key
	}
	newKey := strings.TrimPrefix(key, CustomerDefPrefix)
	v, ok := variablePool.Load(newKey)
	if !ok {
		return key
	} else {
		return v.(string)
	}
}
