/**
* Author: CZ cz.theng@gmail.com
 */

package bucket

import (
	"encoding/xml"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi"
	"path"
	"strconv"
)

/**
// redefine on put_bucket_lifecycle

const (
	LifecycleStatsEnable  = "Enabled"
	LifecycleStatsDisable = "Disabled"
)

type ExpirationDaysInfo struct {
	Days uint
}

type ExpirationDateInfo struct {
	Date string
}

type RuleInfo struct {
	ID         string
	Prefix     string
	Status     string
	Expiration ExpirationDaysInfo
}

type LifecycleConfiguration struct {
	XMLName xml.Name `xml:"LifecycleConfiguration"`
	Rule    []RuleInfo
}
*/

// QueryLifecycle Query bucket's life cycle
// @param name: name of bucket
// @param location:location of name
// @return infos : lifecycle info of buckets
// @return ossapiError : nil on success
func QueryLifecycle(name, location string) (infos []RuleInfo, ossapiError *ossapi.Error) {
	host := name + "." + location + ".aliyuncs.com"
	resource := path.Join("/", name) + "/"
	req := &ossapi.Request{
		Host:     host,
		Path:     "/?lifecycle",
		Method:   "GET",
		Resource: resource,
		SubRes:   []string{"lifecycle"}}
	rsp, err := req.Send()
	if err != nil {
		if _, ok := err.(*ossapi.Error); !ok {
			ossapi.Logger.Error("GetService's Send Error:%s", err.Error())
			ossapiError = ossapi.OSSAPIError
			return
		}
	}
	if rsp.Result != ossapi.ErrSUCC {
		ossapiError = err.(*ossapi.Error)
		return
	}
	bodyLen, err := strconv.Atoi(rsp.HTTPRsp.Header["Content-Length"][0])
	if err != nil {
		ossapi.Logger.Error("GetService's Send Error:%s", err.Error())
		ossapiError = ossapi.OSSAPIError
		return
	}
	body := make([]byte, bodyLen)
	rsp.HTTPRsp.Body.Read(body)
	info := new(LifecycleConfiguration)
	err = xml.Unmarshal(body, info)
	if err != nil {
		ossapi.Logger.Error("xml.Unmarshal body info Error:", err.Error())
		ossapiError = ossapi.OSSAPIError
		return
	}
	infos = info.Rule
	return
}
