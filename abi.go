package msgpack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
	"github.com/bitly/go-simplejson"
)

//ParseAbi parse abiraw to struct for contracts
func ParseAbi(abiRaw []byte) (*ABI, error) {
        abis := &ABIStructs{}
        err := json.Unmarshal(abiRaw, abis)
        if err != nil {
                return &ABI{}, err
        }

        abi := &ABI{}
        abi.Structs = make([]ABIStruct, len(abis.Structs))
        for i := range abi.Structs {
                abi.Structs[i].Fields = New()
        }
        err = json.Unmarshal(abiRaw, abi)
        if err != nil {
                return &ABI{}, err
        }

        return abi, nil
}

//GetAbibyContractName function
func GetAbibyContractName(contractname string) (ABI, error) {
	var abistring string
	NodeIp := "127.0.0.1"
	addr := "http://" + NodeIp + ":8080/rpc"
	params := `service=bottos&method=CoreApi.QueryAbi&request={
			"contract":"%s"}`
	s := fmt.Sprintf(params, contractname)
	respBody, err := http.Post(addr, "application/x-www-form-urlencoded",
		strings.NewReader(s))

	if err != nil {
		fmt.Println(err)
		return ABI{}, err
	}

	defer respBody.Body.Close()
	body, err := ioutil.ReadAll(respBody.Body)
	if err != nil {
		fmt.Println(err)
		return ABI{}, err
	}

	jss, _ := simplejson.NewJson([]byte(body))
	abistring = jss.Get("result").MustString()
	if len(abistring) <= 0 {
		fmt.Println(err)
		return ABI{}, err
	}

	Abi, err := ParseAbi([]byte(abistring))
	if err != nil {
		fmt.Println("Parse abistring", abistring, " to abi failed!")
		return ABI{}, err
	}

	return *Abi, nil
}
