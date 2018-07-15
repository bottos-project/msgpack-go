// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  msgpack go
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package msgpack

import (
	"fmt"
	"reflect"
	"bytes"
	"io"
)

const (
	//BIN16 is byte array type identifier
	BIN16 = 0xc5
	//UINT8 is uint8
	UINT8  = 0xcc
	//UINT16 is uint16
	UINT16 = 0xcd
	//UINT32 is uint32
	UINT32 = 0xce
	//UINT64 is uint64
	UINT64 = 0xcf
	//STR16 is string type identifier
	STR16   = 0xda
	//ARRAY16 is array size type identifier
	ARRAY16 = 0xdc
	//LEN_INT32 value
	LEN_INT32 = 4
	//LEN_INT64 value
	LEN_INT64 = 8
	//MAX16BIT value
	MAX16BIT = 2 << (16 - 1)
	//REGULAR_UINT7_MAX value
	REGULAR_UINT7_MAX  = 2 << (7 - 1)
	//REGULAR_UINT8_MAX value
	REGULAR_UINT8_MAX  = 2 << (8 - 1)
	//REGULAR_UINT16_MAX value
	REGULAR_UINT16_MAX = 2 << (16 - 1)
	//REGULAR_UINT32_MAX value
	REGULAR_UINT32_MAX = 2 << (32 - 1)

	//SPECIAL_INT8 value
	SPECIAL_INT8  = 32
	//SPECIAL_INT16 value
	SPECIAL_INT16 = 2 << (8 - 2)
	//SPECIAL_INT32 value
	SPECIAL_INT32 = 2 << (16 - 2)
	//SPECIAL_INT64 value
	SPECIAL_INT64 = 2 << (32 - 2)
)

//Bytes is []byte type
type Bytes []byte

//ABIAction abi Action(Method)
type ABIAction struct {
	ActionName string `json:"action_name"`
	Type       string `json:"type"`
}

//ABIStruct parameter struct for abi Action(Method)
type ABIStruct struct {
	Name   string    `json:"name"`
	Base   string    `json:"base"`
	Fields *FeildMap `json:"fields"`
}

//ABI struct for abi
type ABI struct {
	Types   []interface{} `json:"types"`
	Structs []ABIStruct   `json:"structs"`
	Actions []ABIAction   `json:"actions"`
	Tables  []interface{} `json:"tables"`
}

//ABIStructs structs for ABI
type ABIStructs struct {
	Structs []struct {
		Name   string            `json:"name"`
		Base   string            `json:"base"`
		Fields map[string]string `json:"fields"`
	} `json:"structs"`
}

//PackUint8 is to pack a given value and writes it into the specified writer.
func PackUint8(writer io.Writer, value uint8) (n int, err error) {
	return writer.Write(Bytes{UINT8, value})
}

//PackUint16 is to pack a given value and writes it into the specified writer.
func PackUint16(writer io.Writer, value uint16) (n int, err error) {
	return writer.Write(Bytes{UINT16, byte(value >> 8), byte(value)})
}

//PackUint32 is to pack a given value and writes it into the specified writer.
func PackUint32(writer io.Writer, value uint32) (n int, err error) {
	return writer.Write(Bytes{UINT32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackUint64 is to pack a given value and writes it into the specified writer.
func PackUint64(writer io.Writer, value uint64) (n int, err error) {
	return writer.Write(Bytes{UINT64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackBin16 is to pack a given value and writes it into the specified writer.
func PackBin16(writer io.Writer, value []byte) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{BIN16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write(value)
	return n1 + n2, err
}

//PackStr16 is to pack a given value and writes it into the specified writer.
func PackStr16(writer io.Writer, value string) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{STR16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write([]byte(value))
	return n1 + n2, err
}

//PackArraySize is to pack a given value and writes it into the specified writer.
func PackArraySize(writer io.Writer, length uint16) (n int, err error) {
	n, err = writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
	if err != nil {
		return n, err
	}
	return n, nil
}

//MarshalAbi is to serialize the message
func MarshalAbi(v interface{}, Abi *ABI, contractName string, method string) ([]byte, error) {
	var err error
	var abi ABI

	if Abi == nil {
		return []byte{}, err
	}
	
	abi = *Abi
	

	writer := &bytes.Buffer{}
	err = EncodeAbi(contractName, method, writer, v, abi, "")
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}


func getAbiFieldsByAbi(contractname string, method string, abi ABI, subStructName string) map[string]interface{} {
	for _, subaction := range abi.Actions {
		if subaction.ActionName != method {
			continue
		}

		structname := subaction.Type

		for _, substruct := range abi.Structs {
			if subStructName != "" {
				if substruct.Name != subStructName {
					continue
				}
			} else if structname != substruct.Name {
				continue
			}

			return substruct.Fields.values
		}
	}

	return nil
}

//EncodeAbi is to encode message
func EncodeAbi(contractName string, method string, w io.Writer, value interface{}, abi ABI, subStructName string) error {
	abiFields := getAbiFieldsByAbi(contractName, method, abi, subStructName)
	if abiFields == nil {
		return fmt.Errorf("EncodeAbi: getAbiFieldsByAbi failed: %s", abi)
	}

	v := reflect.ValueOf(value)
	vt := reflect.TypeOf(value)

	if !v.IsValid() {
		return fmt.Errorf("Not Valid %T\n", value)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		vt = vt.Elem()
		if !v.IsValid() {
			return fmt.Errorf("Nil Ptr: %T\n", value)
		}
	}

	count := v.NumField()
	PackArraySize(w, uint16(count))

	for i := 0; i < count; i++ {
		fieldname := vt.Field(i).Tag.Get("json")
		vals := v.Field(i).Interface()

		types := reflect.TypeOf(vals)
		val := reflect.ValueOf(vals)

		if _, ok := abiFields[fieldname]; !ok {
			return fmt.Errorf("%s is not in abiFields [%s]!", fieldname, abiFields)
		}

		switch abiFields[fieldname] {
		case "string":
			PackStr16(w, val.String())
		case "uint8":
			PackUint8(w, uint8(val.Uint()))
		case "uint16":
			PackUint16(w, uint16(val.Uint()))
		case "uint32":
			PackUint32(w, uint32(val.Uint()))
		case "uint64":
			PackUint64(w, uint64(val.Uint()))
		case "bytes":
			t := reflect.TypeOf(v.Field(i).Interface())
			if t.Elem().Kind() == reflect.Uint8 {
				PackBin16(w, val.Bytes())
			} else {
				return fmt.Errorf("Unsupported Slice Type")
			}
		default:
			t := reflect.TypeOf(v.Field(i).Interface())
			if t.Kind() == reflect.Struct || t.Kind() == reflect.Ptr {
				EncodeAbi(contractName, method, w, v.Field(i).Interface(), abi, fieldname)
			} else {
				return fmt.Errorf("Unsupported Type: %v", types)
			}
		}
	}

	return nil
}

//getAbiFieldsByAbiEx function
func getAbiFieldsByAbiEx(contractname string, method string, abi ABI, subStructName string) *FeildMap {
	for _, subaction := range abi.Actions {
		if subaction.ActionName != method {
			continue
		}
		structname := subaction.Type

		for _, substruct := range abi.Structs {
			if subStructName != "" {
				if substruct.Name != subStructName {
					continue
				}
			} else if structname != substruct.Name {
				continue
			}

			return substruct.Fields
		}
	}

	return nil
}

//EncodeAbiEx is to encode message
func EncodeAbiEx(contractName string, method string, w io.Writer, value map[string]interface{}, abi ABI, subStructName string) error {
        abiFieldsAttr := getAbiFieldsByAbiEx(contractName, method, abi, subStructName)
	if abiFieldsAttr == nil {
		return fmt.Errorf("EncodeAbiEx: getAbiFieldsByAbi failed: %s", abi)

	}

	abiFields := abiFieldsAttr.GetStringPair()
	
	count  := len(abiFields)
	count2 := len(value)
	
	if count != count2 {
		return fmt.Errorf("EncodeAbiEx: fields number mismatch! count: %d, count2: %d", count, count2)
	}
	
	if (count <= 0) {
		return fmt.Errorf("EncodeAbiEx: count is 0!", count)
	}

	PackArraySize(w, uint16(count))

		for _, abiValTypeAttr := range abiFields {
			
			abiValKey   := abiValTypeAttr.Key
			abiValType := abiValTypeAttr.Value

			val, ok := value[abiValKey]
			
			if !ok {
				return fmt.Errorf("EncodeAbiEx: value abiValKey %s not found in map", abiValKey)
			}
			
			valType := reflect.TypeOf(val).Name()
			
			if reflect.ValueOf(val).Kind() == reflect.Slice {
				valType = reflect.TypeOf(val).Elem().Name()
				if valType == "uint8"	{
					valType = "bytes"
				}
			}
				
			if valType != abiValType {
				return fmt.Errorf("EncodeAbiEx: abiValType %s mismatch to valType %s", abiValType, valType)
			}

			switch abiValType {
				case "string":
					PackStr16(w, val.(string))
				case "uint8":
					PackUint8(w, val.(uint8))
				case "uint16":
					PackUint16(w, val.(uint16))
				case "uint32":
					PackUint32(w, val.(uint32))
				case "uint64":
					PackUint64(w, val.(uint64))
				case "bytes":
					PackBin16(w, val.([]byte))
				default:
					if reflect.ValueOf(value[abiValKey]).Kind() == reflect.Struct {
						EncodeAbi(contractName, method, w, value[abiValKey], abi, abiValKey)
					} else {
						return fmt.Errorf("Unsupported Type: %v | %v", valType, abiValType)
					}
				}
		}

	return nil
}

func Setmapval(structmap map[string]interface{}, key string, val interface{}) {
        structmap[key] = val
}

//MarshalAbiEx is to serialize the message
func MarshalAbiEx(v map[string]interface{}, Abi *ABI, contractName string, method string) ([]byte, error) {
	var err error
	var abi ABI
	
	
	if Abi == nil {
		return []byte{}, err
	}
	
	abi = *Abi

	writer := &bytes.Buffer{}
	err = EncodeAbiEx(contractName, method, writer, v, abi, "")
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}
