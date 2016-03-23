/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Testing GetObjectInfo().
func TestGetObjectInfo(t *testing.T) {
	directory, e := ioutil.TempDir("", "minio-get-objinfo-test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(directory)

	// Create the filesystem.
	fs, err := New(directory, 0)
	if err != nil {
		t.Fatal(err)
	}
	// This bucket is used for testing getObjectInfo operations.
	err = fs.MakeBucket("test-getobjectinfo")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.CreateObject("test-getobjectinfo", "Asia/asiapics.jpg", "", int64(len("asiapics")), bytes.NewBufferString("asiapics"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resultCases := []ObjectInfo{
		// ObjectInfo -1.
		// ObjectName set to a existing object in the test case (Test case 14).
		{Bucket: "test-getobjectinfo", Name: "Asia/asiapics.jpg", ContentType: "image/jpeg", IsDir: false},
	}
	testCases := []struct {
		rootPath   string
		bucketName string
		objectName string

		// Expected output of GetObjectInfo.
		result ObjectInfo
		err    error
		// Flag indicating whether the test is expected to pass or not.
		shouldPass bool
	}{
		// Test cases with invalid bucket names ( Test number 1-4 ).
		{fs.path, ".test", "", ObjectInfo{}, BucketNameInvalid{Bucket: ".test"}, false},
		{fs.path, "Test", "", ObjectInfo{}, BucketNameInvalid{Bucket: "Test"}, false},
		{fs.path, "---", "", ObjectInfo{}, BucketNameInvalid{Bucket: "---"}, false},
		{fs.path, "ad", "", ObjectInfo{}, BucketNameInvalid{Bucket: "ad"}, false},
		// Test cases with valid but non-existing bucket names (Test number 5-7).
		{fs.path, "abcdefgh", "abc", ObjectInfo{}, BucketNotFound{Bucket: "abcdefgh"}, false},
		{fs.path, "ijklmnop", "efg", ObjectInfo{}, BucketNotFound{Bucket: "ijklmnop"}, false},
		// Test cases with valid but non-existing bucket names and invalid object name (Test number 8-9).
		{fs.path, "abcdefgh", "", ObjectInfo{}, ObjectNameInvalid{Bucket: "abcdefgh", Object: ""}, false},
		{fs.path, "ijklmnop", "", ObjectInfo{}, ObjectNameInvalid{Bucket: "ijklmnop", Object: ""}, false},
		// Test cases with non-existing object name with existing bucket (Test number 10-12).
		{fs.path, "test-getobjectinfo", "Africa", ObjectInfo{}, ObjectNotFound{Bucket: "test-getobjectinfo", Object: "Africa"}, false},
		{fs.path, "test-getobjectinfo", "Antartica", ObjectInfo{}, ObjectNotFound{Bucket: "test-getobjectinfo", Object: "Antartica"}, false},
		{fs.path, "test-getobjectinfo", "Asia/myfile", ObjectInfo{}, ObjectNotFound{Bucket: "test-getobjectinfo", Object: "Asia/myfile"}, false},
		// Test case with existing bucket but object name set to a directory (Test number 13).
		{fs.path, "test-getobjectinfo", "Asia", ObjectInfo{}, ObjectNotFound{Bucket: "test-getobjectinfo", Object: "Asia"}, false},
		// Valid case with existing object (Test number 14).
		{fs.path, "test-getobjectinfo", "Asia/asiapics.jpg", resultCases[0], nil, true},
	}
	for i, testCase := range testCases {
		result, err := fs.GetObjectInfo(testCase.bucketName, testCase.objectName)
		if err != nil && testCase.shouldPass {
			t.Errorf("Test %d: Expected to pass, but failed with: <ERROR> %s", i+1, err.Cause.Error())
		}
		if err == nil && !testCase.shouldPass {
			t.Errorf("Test %d: Expected to fail with <ERROR> \"%s\", but passed instead", i+1, testCase.err.Error())
		}
		// Failed as expected, but does it fail for the expected reason.
		if err != nil && !testCase.shouldPass {
			if testCase.err.Error() != err.Cause.Error() {
				t.Errorf("Test %d: Expected to fail with error \"%s\", but instead failed with error \"%s\" instead", i+1, testCase.err.Error(), err.Cause.Error())
			}
		}

		// Test passes as expected, but the output values are verified for correctness here.
		if err == nil && testCase.shouldPass {
			if testCase.result.Bucket != result.Bucket {
				t.Fatalf("Test %d: Expected Bucket name to be '%s', but found '%s' instead", i+1, testCase.result.Bucket, result.Bucket)
			}
			if testCase.result.Name != result.Name {
				t.Errorf("Test %d: Expected Object name to be %s, but instead found it to be %s", i+1, testCase.result.Name, result.Name)
			}
			if testCase.result.ContentType != result.ContentType {
				t.Errorf("Test %d: Expected Content Type of the object to be %v, but instead found it to be %v", i+1, testCase.result.ContentType, result.ContentType)
			}
			if testCase.result.IsDir != result.IsDir {
				t.Errorf("Test %d: Expected IsDir flag of the object to be %v, but instead found it to be %v", i+1, testCase.result.IsDir, result.IsDir)
			}
		}
	}
}

// Testing getObjectInfo().
func TestGetObjectInfoCore(t *testing.T) {
	directory, e := ioutil.TempDir("", "minio-get-objinfo-test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(directory)

	// Create the filesystem.
	fs, err := New(directory, 0)
	if err != nil {
		t.Fatal(err)
	}
	// This bucket is used for testing getObjectInfo operations.
	err = fs.MakeBucket("test-getobjinfo")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.CreateObject("test-getobjinfo", "Asia/asiapics.jpg", "", int64(len("asiapics")), bytes.NewBufferString("asiapics"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resultCases := []ObjectInfo{
		// ObjectInfo - 1.
		// ObjectName object name set to a existing directory in the test case.
		{Bucket: "test-getobjinfo", Name: "Asia", Size: 0, ContentType: "application/octet-stream", IsDir: true},
		// ObjectInfo -2.
		// ObjectName set to a existing object in the test case.
		{Bucket: "test-getobjinfo", Name: "Asia/asiapics.jpg", Size: int64(len("asiapics")), ContentType: "image/jpeg", IsDir: false},
		// ObjectInfo-3.
		// Object name set to a non-existing object in the test case.
		{Bucket: "test-getobjinfo", Name: "Africa", Size: 0, ContentType: "image/jpeg", IsDir: false},
	}
	testCases := []struct {
		rootPath   string
		bucketName string
		objectName string

		// Expected output of getObjectInfo.
		result ObjectInfo
		err    error

		// Flag indicating whether the test is expected to pass or not.
		shouldPass bool
	}{
		// Testcase with object name set to a existing directory ( Test number 1).
		{fs.path, "test-getobjinfo", "Asia", resultCases[0], nil, true},
		// ObjectName set to a existing object ( Test number 2).
		{fs.path, "test-getobjinfo", "Asia/asiapics.jpg", resultCases[1], nil, true},
		// Object name set to a non-existing object. (Test number 3).
		{fs.path, "test-getobjinfo", "Africa", resultCases[2], fmt.Errorf("%s", filepath.FromSlash("test-getobjinfo/Africa")), false},
	}
	for i, testCase := range testCases {
		result, err := getObjectInfo(testCase.rootPath, testCase.bucketName, testCase.objectName)
		if err != nil && testCase.shouldPass {
			t.Errorf("Test %d: Expected to pass, but failed with: <ERROR> %s", i+1, err.Cause.Error())
		}
		if err == nil && !testCase.shouldPass {
			t.Errorf("Test %d: Expected to fail with <ERROR> \"%s\", but passed instead", i+1, testCase.err.Error())
		}
		// Failed as expected, but does it fail for the expected reason.
		if err != nil && !testCase.shouldPass {
			if !strings.Contains(err.Cause.Error(), testCase.err.Error()) {
				t.Errorf("Test %d: Expected to fail with error \"%s\", but instead failed with error \"%s\" instead", i+1, testCase.err.Error(), err.Cause.Error())
			}
		}

		// Test passes as expected, but the output values are verified for correctness here.
		if err == nil && testCase.shouldPass {
			if testCase.result.Bucket != result.Bucket {
				t.Fatalf("Test %d: Expected Bucket name to be '%s', but found '%s' instead", i+1, testCase.result.Bucket, result.Bucket)
			}
			if testCase.result.Name != result.Name {
				t.Errorf("Test %d: Expected Object name to be %s, but instead found it to be %s", i+1, testCase.result.Name, result.Name)
			}
			if testCase.result.ContentType != result.ContentType {
				t.Errorf("Test %d: Expected Content Type of the object to be %v, but instead found it to be %v", i+1, testCase.result.ContentType, result.ContentType)
			}
			if testCase.result.IsDir != result.IsDir {
				t.Errorf("Test %d: Expected IsDir flag of the object to be %v, but instead found it to be %v", i+1, testCase.result.IsDir, result.IsDir)
			}
		}
	}
}

func BenchmarkGetObject(b *testing.B) {
	// Make a temporary directory to use as the filesystem.
	directory, e := ioutil.TempDir("", "minio-benchmark-getobject")
	if e != nil {
		b.Fatal(e)
	}
	defer os.RemoveAll(directory)

	// Create the filesystem.
	filesystem, err := New(directory, 0)
	if err != nil {
		b.Fatal(err)
	}

	// Make a bucket and put in a few objects.
	err = filesystem.MakeBucket("bucket")
	if err != nil {
		b.Fatal(err)
	}

	text := "Jack and Jill went up the hill / To fetch a pail of water."
	hasher := md5.New()
	hasher.Write([]byte(text))
	sum := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	for i := 0; i < 10; i++ {
		_, err = filesystem.CreateObject("bucket", "object"+strconv.Itoa(i), sum, int64(len(text)), bytes.NewBufferString(text), nil)
		if err != nil {
			b.Fatal(err)
		}
	}

	var w bytes.Buffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := filesystem.GetObject(&w, "bucket", "object"+strconv.Itoa(i%10), 0, 0)
		if err != nil {
			b.Error(err)
		}
		if n != int64(len(text)) {
			b.Errorf("GetObject returned incorrect length %d (should be %d)\n", n, int64(len(text)))
		}
	}
}