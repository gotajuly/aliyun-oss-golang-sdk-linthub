/**
* Author: CZ cz.theng@gmail.com
 */

package mupload

import (
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi"
	"path"
	"strconv"
)

//UploadPartRstInfo  Return Etag
type UploadPartRstInfo struct {
	Etag string
}

//UploadPartInfo  upload part info
type UploadPartInfo struct {
	ObjectName string
	BucketName string
	Location   string
	PartNumber int
	UploadID   string
	Data       []byte
	CntType    string
}

// Append a data slice
// @param partInfo : data slice ifo
// @return rstInfo : return response
// @return ossapiError : nil on success
func Append(partInfo *UploadPartInfo) (rstInfo *UploadPartRstInfo, ossapiError *ossapi.Error) {
	if partInfo == nil {
		ossapiError = ossapi.ArgError
		return
	}
	resource := path.Join("/", partInfo.BucketName, partInfo.ObjectName)
	host := partInfo.BucketName + "." + partInfo.Location + ".aliyuncs.com"
	req := &ossapi.Request{
		Host:     host,
		Path:     "/" + partInfo.ObjectName + "?partNumber=" + strconv.FormatUint(uint64(partInfo.PartNumber), 10) + "&uploadId=" + partInfo.UploadID,
		Method:   "PUT",
		Body:     partInfo.Data,
		SubRes:   []string{"partNumber=" + strconv.FormatUint(uint64(partInfo.PartNumber), 10) + "&uploadId=" + partInfo.UploadID},
		Resource: resource}
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
	rstInfo = new(UploadPartRstInfo)
	rstInfo.Etag = rsp.HTTPRsp.Header["Etag"][0]
	return
}
