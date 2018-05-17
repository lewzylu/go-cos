package main

import (
	"fmt"
	"context"
	"net/url"
	"strings"
	"net/http"
	"encoding/json"
 	"github.com/lewzylu/go-cos"
        "github.com/lewzylu/go-cos/debug"
	"github.com/QcloudApi/qcloud_sign_golang"

)
type Credent struct{
        SessionToken string `json:"sessionToken"`
        TmpSecretId  string `json:"tmpSecretId"`
        TmpSecretKey string `json:"tmpSecretKey"`
}
type Data struct{
        Credentials Credent `json:"credentials`
                
}
type Response struct{
        Dat Data `json:"data"`
}
func main() {
        // 替换实际的 SecretId 和 SecretKey
        secretId := ""
        secretKey := ""
    
        // 配置
        config := map[string]interface{} {"secretId" : secretId, "secretKey" : secretKey, "debug" : false}
    
        // 请求参数
        params := map[string]interface{} {"Region" : "gz", "Action" : "GetFederationToken","name":"alantong","policy":"{\"statement\": [{\"action\": [\"name/cos:GetObject\",\"name/cos:PutObject\"],\"effect\": \"allow\",\"resource\":[\"qcs::cos:ap-guangzhou:uid/1251668577:prefix//1251668577/alantest/*\"]}],\"version\": \"2.0\"}"   }
    
        // 发送请求
        retData, err := QcloudApi.SendRequest("sts", params, config)
        if err != nil{
            fmt.Print("Error.", err)
            return
        }
        r := &Response{}
        err = json.Unmarshal([]byte(retData), r)
        if err != nil {
            fmt.Println(err);
            return
        }
        //获取临时ak、sk、token
        tmp_secId := r.Dat.Credentials.TmpSecretId
        tmp_secKey := r.Dat.Credentials.TmpSecretKey
        token := r.Dat.Credentials.SessionToken

    	//fmt.Println("token:", token)
	u, _ := url.Parse("https://alangz-1251668577.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: tmp_secId, 
			SecretKey: tmp_secKey,
			SessionToken: token,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/objectPut.go"
	f := strings.NewReader("test")

	_, err = c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	name = "test/put_option.go"
	f = strings.NewReader("test xxx")
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			//XCosACL: "public-read",
			XCosACL: "private",
		},
	}
	_, err = c.Object.Put(context.Background(), name, f, opt)
	if err != nil {
		panic(err)
	}
}
