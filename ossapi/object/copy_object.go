/**
* Author: CZ cz.theng@gmail.com
 */

package object

import (
	"encoding/xml"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi"
	"path"
	"strconv"
)

// enmu
const (
	//COPY will not replace if dest is exist
	COPY = "COPY"
	// REPLACE wiil repalce .
	REPLACE = "REPLACE"
)

// CopyConditionInfo is copy condition
type CopyConditionInfo struct {
	ETAG         string
	Date         string
	LastModify   string
	LastUnModify string
}

//CopyInfo is Copy Info with source and dest
type CopyInfo struct {
	ObjectName string
	BucketName string
	Location   string
	Source     string
	Directive  string
	Encryption string
	ACL        string
}

// CopyResultInfo is return XML info
type CopyResultInfo struct {
	XMLName      xml.Name `xml:"CopyObjectResult"`
	ETag         string   `xml:"ETag"`
	LastModified string   `xml:"LastModified"`
}

// Copy a object to creat a bucket
// @param copyInfo : source and dest object and bucket
// @param copyConnInfo : conndition to copy
// @return rstInfo: return info
// @return ossapiError : nil on success
func Copy(copyInfo *CopyInfo, copyConnInfo *CopyConditionInfo) (rstInfo *CopyResultInfo, ossapiError *ossapi.Error) {
	if copyInfo == nil {
		ossapiError = ossapi.ArgError
		return
	}
	resource := path.Join("/", copyInfo.BucketName, copyInfo.ObjectName)
	host := copyInfo.BucketName + "." + copyInfo.Location + ".aliyuncs.com"
	req := &ossapi.Request{
		Host:     host,
		Path:     "/" + copyInfo.ObjectName,
		Method:   "PUT",
		Resource: resource}
	if copyConnInfo != nil {
		if copyConnInfo.ETAG != "" {
			req.AddXOSS("x-oss-copy-source-if-match", copyConnInfo.ETAG)
		}
		if copyConnInfo.Date != "" {
			req.AddXOSS("x-oss-copy-source-if-none-match", copyConnInfo.Date)
		}
		if copyConnInfo.LastUnModify != "" {
			req.AddXOSS("x-oss-copy-source-if-unmodified-since", copyConnInfo.LastUnModify)
		}
		if copyConnInfo.LastModify != "" {
			req.AddXOSS("x-oss-copy-source-if-modified-since", copyConnInfo.LastModify)
		}
	}
	if copyInfo.ObjectName != "" {
		req.AddXOSS("x-oss-copy-source", copyInfo.Source)
	}
	if copyInfo.Directive != "" {
		req.AddXOSS("x-oss-metadata-directive", copyInfo.Directive)
	}
	if copyInfo.Encryption != "" {
		req.AddXOSS("x-oss-server-side-encryption", copyInfo.Encryption)
	}
	if copyInfo.ACL != "" {
		req.AddXOSS("x-oss-object-acl", copyInfo.ACL)
	}

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
	info := new(CopyResultInfo)
	xml.Unmarshal(body, info)
	rstInfo = info
	return
}
