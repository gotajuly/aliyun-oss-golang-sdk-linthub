/**
* Author: CZ cz.theng@gmail.com
 */

package mupload

import (
	"fmt"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi"
	"github.com/cz-it/aliyun-oss-golang-sdk/ossapi/bucket"
	"testing"
)

func TestCompleteMultipartUpload(t *testing.T) {
	if nil != ossapi.Init("v8P430U3UcILP6KA", "EB9v8yL2aM07YOgtO1BdfrXtdxa4A1") {
		t.Fail()
	}
	initInfo := &InitInfo{
		CacheControl:       "no-cache",
		ContentDisposition: "attachment;filename=oss_download.jpg",
		ContentEncoding:    "utf-8",
		Expires:            "Fri, 28 Feb 2012 05:38:42 GMT",
		Encryption:         "AES256"}
	var info *InitRstInfo
	var err *ossapi.Error
	if info, err = Init("a.c", "test-mupload", bucket.LHangzhou, initInfo); err != nil {
		fmt.Println(err.ErrNo, err.HTTPStatus, err.ErrMsg, err.ErrDetailMsg)
	} else {
		t.Log("Init Multiple Upload Success!")
		fmt.Println(info)
	}
	var partData []byte
	for i := 0; i < 10250; i++ {
		partData = append(partData, "1234567890"...)
	}

	partInfo := &UploadPartInfo{
		ObjectName: "a.c",
		BucketName: "test-mupload",
		Location:   bucket.LHangzhou,
		UploadID:   info.UploadId,
		PartNumber: 1,
		Data:       partData[:100*1024],
		CntType:    "text/html"}

	var i1 PartInfo
	if info, err := Append(partInfo); err != nil {
		fmt.Println(err.ErrNo, err.HTTPStatus, err.ErrMsg, err.ErrDetailMsg)
	} else {
		i1.ETag = info.Etag
		i1.PartNumber = 1
		t.Log("UploadPart Success!")
	}

	partInfo.PartNumber = 2
	var i2 PartInfo
	if info, err := Append(partInfo); err != nil {
		fmt.Println(err.ErrNo, err.HTTPStatus, err.ErrMsg, err.ErrDetailMsg)
	} else {
		i2.ETag = info.Etag
		i2.PartNumber = 2
		t.Log("UploadPart Success!")
	}

	partsInfo := &PartsInfo{Part: []PartInfo{i1, i2}}
	if info, err := Complete("a.c", "test-mupload", bucket.LHangzhou, info.UploadId, partsInfo); err != nil {
		fmt.Println(err.ErrNo, err.HTTPStatus, err.ErrMsg, err.ErrDetailMsg)
	} else {
		t.Log(" CompleteMultipartUpload Success!")
		fmt.Println(info)
	}

}
