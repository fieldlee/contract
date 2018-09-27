package common

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Json2map(str string) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func SetStructByJsonName(ptr interface{}, fields map[string]interface{}) error {
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("json")
		if name == "" {
			//无tag的情况，将field首字母变小写
			name = strings.ToLower(fieldInfo.Name[0:1]) + fieldInfo.Name[1:]
		}
		//去掉逗号后面内容 如 `json:"txId,omitempty"`
		name = strings.Split(name, ",")[0]
		if value, ok := fields[name]; ok {
			//给结构体赋值
			//保证赋值时数据类型一致
			if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func GetUserFromCertification(stub shim.ChaincodeStubInterface) (username string) {
	userBytes, err := stub.GetCreator()
	if err != nil {
		return
	}
	certStart := bytes.IndexAny(userBytes, "-----") // Devin:I don't know why sometimes -----BEGIN is invalid, so I use -----
	if certStart == -1 {
		fmt.Errorf("No certificate found")
		return
	}
	certText := userBytes[certStart:]

	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
		return
	}
	fmt.Println(string(certText))
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
		return
	}
	fmt.Println(cert)

	username = cert.Subject.CommonName
	fmt.Println("Name:" + username)

	return
}
