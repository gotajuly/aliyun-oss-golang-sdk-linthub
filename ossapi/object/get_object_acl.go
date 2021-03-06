/**
* Author: CZ cz.theng@gmail.com
 */

package object

import (
	"encoding/xml"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi/bucket"
	"path"
	"strconv"
)

/*
// define in get_bucket.go
type OwnerInfo struct {
	ID          string
	DisplayName string
}

type AccessControlListInfo struct {
	Grant string
}

type ACLInfo struct {
	XMLName           xml.Name `xml:"AccessControlPolicy"`
	Owner             OwnerInfo
	AccessControlList AccessControlListInfo
}
*/

// QueryACL Query bucket's ACL info
// @param objName : name of object
// @param bucketName : name of bucket
// @param locaton : location of bucket
// @return info: ACL info
// @retun ossapiError : nil on success
func QueryACL(objName, bucketName, location string) (info *bucket.ACLInfo, ossapiError *ossapi.Error) {
	host := bucketName + "." + location + ".aliyuncs.com"
	resource := path.Join("/", bucketName, objName)
	req := &ossapi.Request{
		Host:     host,
		Path:     "/" + objName + "?acl",
		Method:   "GET",
		Resource: resource,
		SubRes:   []string{"acl"}}
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
	info = new(bucket.ACLInfo)
	xml.Unmarshal(body, info)
	return
}
