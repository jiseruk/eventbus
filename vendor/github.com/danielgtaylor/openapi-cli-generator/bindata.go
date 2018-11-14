// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/commands.tmpl
// templates/main.tmpl

package main


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataTemplatesCommandstmpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\xdf\x6f\xdb\x38\x12\x7e\x96\xfe\x8a\xa9\xd0\x2b\xa4\x8b\x23\x5f" +
	"\x71\x87\x7b\xf0\xc2\x0f\x69\xd2\x34\x01\x9a\xa4\x9b\x34\x4f\xd9\x00\x61\xac\x91\x4c\x94\x22\x15\x8a\x4e\x93\x15" +
	"\xf4\xbf\x2f\x86\xa4\x64\xd9\x56\xb6\x09\xb0\x7e\xb1\xc4\x19\xce\x8f\x6f\x3e\x7e\xe2\x74\x0a\x87\x2a\x43\x28\x50" +
	"\xa2\x66\x06\x33\xb8\x7f\x06\x55\xa1\x64\x15\xdf\x5f\x08\xbe\xef\x0d\x4a\xa7\x70\x74\x01\xe7\x17\xdf\xe1\xf3\xd1" +
	"\xe9\xf7\x34\x9c\x4e\xe1\x0a\x11\x96\xc6\x54\xf5\x6c\x3a\x2d\xb8\x59\xae\xee\xd3\x85\x2a\xa7\x19\x93\x1c\x45\x61" +
	"\xd8\xb3\x50\x7a\x3a\x1a\x2b\x0c\x2b\xb6\xf8\xc1\x0a\x84\x92\x71\x19\x86\xbc\xac\x94\x36\x10\x87\x41\xd3\x00\xcf" +
	"\x21\x3d\xb5\x0b\x75\x7a\x5c\x1a\x68\xdb\x28\x2f\x4d\xd4\x34\x80\x32\x83\xb6\xdd\x71\xba\x32\x9a\xcb\xa2\x26\xc7" +
	"\xda\x3d\x0e\x9c\xc3\x20\x7a\x5b\x6d\xd3\x85\xe0\xd1\xe6\x2e\x5d\x4f\xff\x44\xad\x84\x2a\xa6\x42\x15\x5b\xc6\xba" +
	"\xca\x3f\xfe\x77\xba\x50\xf7\x9a\x8d\x5a\x1e\x79\x85\x3a\x0a\x93\x30\x6c\x1a\x78\x2f\x59\x89\x30\x9b\x43\x7a\x4e" +
	"\x0f\x6d\x6b\x17\x59\xc5\xed\xda\x17\xd5\xad\x86\xf9\x4a\x2e\xa0\xb3\xb5\xed\x15\xea\x47\xd4\x75\x9c\xc0\xcd\x6d" +
	"\xc9\xaa\x1b\xd7\xe7\xad\xfb\x83\x26\x0c\x34\x9a\x95\x96\x63\xd6\x26\x0c\x08\x30\xcd\x64\x81\xf0\xbe\xb6\x81\x6c" +
	"\x36\x1f\xd3\x22\x1a\x04\xa3\xfb\x82\x20\xca\xb0\x5e\x68\x5e\x19\xae\x64\x34\x03\x02\xd6\xc7\x48\x8f\xd6\x16\x82" +
	"\x7e\xe2\xfc\x57\x5a\x6c\xf9\x5d\x5f\x7e\xed\xed\xed\xc4\x55\xd3\x4d\xb2\x0d\x77\x7b\xbd\xc4\x82\xd7\x06\x75\x5c" +
	"\xaf\xee\x17\xaa\x2c\x99\xcc\xe0\x5e\x29\x91\xd8\x3e\x95\x32\x54\xfd\x42\xf0\xf4\x52\x29\x13\x86\x01\xcf\x61\xe0" +
	"\x49\x65\x5b\xa7\x39\x7c\xb0\x43\x49\x0f\x9d\xc5\xf6\x73\x5d\xa3\x2f\x4e\x3a\xa4\x5d\x59\x57\x4b\xa5\x8d\x33\xa4" +
	"\xdf\xb9\x11\x6b\xcb\x57\x25\x8b\x99\xcd\x76\xc6\xf4\x8f\x4c\xfd\x94\xb1\xf5\xda\x6a\x3e\x21\x67\xea\x07\x50\xd4" +
	"\x68\x8b\xe8\x2a\x4c\x6d\x70\x98\x6f\x45\x1f\x7a\x50\x12\x98\xbf\x22\x0b\x01\x16\x0e\xa6\xa9\x2a\xe2\x2c\x59\x69" +
	"\xa0\x17\xdd\x9b\x9f\xa9\xc5\x35\xb6\xb0\x05\x41\xc5\x34\x2b\x6b\xf2\xb3\x84\x4c\xcf\xf1\x67\x9c\x84\x64\x79\x64" +
	"\x1a\xf0\x89\x95\x95\xc0\x1a\xdc\xec\xed\xfa\x3a\x0d\x3e\xd9\xf8\x9f\x3b\x27\xc7\x98\xa0\xdf\xb4\x37\x87\x08\x20" +
	"\x82\xbd\x7e\x2c\x1d\xe8\xdf\x98\x59\xc6\x09\xec\x41\x64\x07\xdc\xd7\x9b\x5e\xd7\x84\x82\x5d\xc4\x27\x68\xdb\x3f" +
	"\x64\xe4\x73\xf6\xe7\x36\x08\x16\x65\x46\x89\x47\xe6\xb8\x1e\xa4\x8f\xe4\xe9\xd7\x34\xfb\x56\x1a\x0e\x04\x67\xf5" +
	"\xba\xd2\xc0\xbf\xcf\xe0\x66\x83\xdd\x6e\x83\xeb\x72\x67\x4f\x10\xd8\xf8\xeb\xe0\xce\xbb\xa3\x6e\xd0\xd1\x79\x67" +
	"\x79\xc8\x26\x37\xfc\x3e\xc4\x4b\x74\xb2\x0c\xe8\x78\x14\x04\x1e\xe9\x59\x3f\x17\xb7\x7c\xa0\x8b\x7a\x06\x0e\x8d" +
	"\x33\x2e\x79\xb9\x2a\xcf\x69\x2d\x6e\x1a\x10\x28\x21\xbd\xc4\x87\x15\xd7\x98\x7d\x73\xc3\x6e\x5b\x1f\xef\x72\x25" +
	"\x67\x40\x74\x88\x09\xd2\x7f\x6f\xe0\x39\x01\xa6\x8b\xba\x07\xc6\xd3\x25\x08\xd6\x4a\xe1\x08\xf3\x05\x8d\x13\xda" +
	"\x38\x72\x26\xa2\x23\xfd\xe8\xfc\x39\xdf\xf9\x1c\xa2\xa8\xdb\xdf\x05\x98\x8f\xc9\xd8\x4d\x1f\xf3\x54\x9a\x2e\xe0" +
	"\x3e\x97\x19\x3e\x45\xc9\xed\x8d\x55\x91\x5b\x8f\x71\xe8\xfe\x57\x5a\x50\x31\xce\x75\xcf\xa2\x46\xec\xb2\x67\x29" +
	"\xf4\xd3\xe9\x08\xcb\x27\xf0\xde\x12\xde\xf2\x76\x07\x95\x7e\x9a\xc4\x15\x7c\xf0\xbe\xe9\xa9\x84\xa8\x62\x66\x19" +
	"\x0d\x28\x40\x59\xe7\xfe\x54\xd4\xe9\x25\x56\x82\x2d\x30\x5e\x69\x31\xa1\xf9\xde\x35\x77\x6d\x4b\xed\xb9\x00\x5e" +
	"\xbd\x9b\xe6\xae\xbd\xa3\x91\x5b\x64\x6f\xc8\x4e\xcd\xdf\x4e\xe0\x63\xb2\x4e\x3d\xe4\xd1\x26\xed\x83\x40\xe3\x43" +
	"\x27\x72\x87\x82\xa3\x34\x29\xb5\x7b\x86\x66\xa9\xc8\x2b\x4e\x48\x54\xa9\x8a\xe4\x1f\x6c\xfd\x61\x85\xfa\x79\xd8" +
	"\x3b\x55\x31\x07\x8d\x0f\xe9\x41\x96\xfd\x4e\x56\x47\xd6\xf3\x4e\x3a\xb7\xfa\x1b\x36\x47\x2a\xb8\x93\x61\x89\x2c" +
	"\x43\xfd\x62\x8a\x13\x6b\x7e\x7d\x8e\xbf\x01\x70\x70\xac\x2f\xac\x7c\x32\x31\x02\xc1\xfa\x83\x4b\x60\x39\x85\x24" +
	"\x4e\xba\x12\x0e\xbf\x9e\x76\x55\x24\x69\x6c\xc5\xfb\xb9\xa2\xd7\xae\x06\x9e\xc3\x66\x90\x77\x96\xeb\xd6\xed\x9c" +
	"\x0b\xab\x6f\x5d\xa3\x5e\x98\xf0\x01\xc6\xb1\x7e\x05\xd8\x79\x69\xd2\xab\x4a\x73\x69\xf2\x38\xfa\xd7\x63\x34\xd9" +
	"\xcc\x9e\x24\xc3\x5c\x83\x01\xbc\x00\xfd\x6b\xb0\x7f\x5b\xca\xc1\x3c\x82\x36\xdc\x5e\x0f\x07\xc4\x4b\x0f\x99\x3c" +
	"\x61\x8f\xf8\x49\x65\xcf\xeb\x3d\xf7\x2a\x7b\x9e\x00\x6a\xdd\x71\xff\x0b\x1a\xf2\x70\x55\x9d\x61\xc6\x99\x9f\xc0" +
	"\x80\x16\x23\xba\xd7\xb6\xb3\xdb\xc1\x8c\x28\xe0\xbb\x39\x48\x2e\xd6\xe3\x10\xaa\x48\x8f\x99\x61\x22\x4e\xd2\xcf" +
	"\x5a\xc7\xa8\x75\x92\x9e\xd5\x45\x1c\x5d\x4b\x76\x2f\x10\x8c\x82\x02\x0d\x50\x49\x9d\xce\xf5\x4a\x44\x41\xc9\x40" +
	"\x51\x07\x7a\x37\x0e\xe8\xa1\x92\x06\xa5\xd9\xa7\xca\xa3\x09\xec\xb6\x92\xa4\xd4\xa3\x97\x56\x0a\x9b\x6c\x43\xb8" +
	"\x89\x20\x01\x73\xc2\x64\x26\xf0\x13\xe6\x4a\x23\x49\xfa\xc4\x93\x77\x42\xd9\x93\x5e\x44\xea\xaa\x87\x93\xaa\x3a" +
	"\x52\xf1\x5a\xb2\x47\x50\x79\x19\x14\xc2\x17\x6b\x03\x39\xe3\x02\xb3\x0e\x90\xae\x22\xba\x40\x64\xb8\x50\x19\x66" +
	"\xc0\xa5\x41\x9d\xb3\x05\x36\xed\x46\x2a\x3f\xd1\x6b\x59\x32\x5d\x2f\x99\x88\x5d\x75\x1f\xfc\xbe\xe4\xb7\xb7\x15" +
	"\xd4\xc7\x11\x74\xff\xa5\x58\x4a\xd6\xf8\x42\x7d\x5d\x6d\xae\x04\x87\xdd\x41\x4e\x37\xcc\x2d\xe8\xa8\xa2\xae\xa0" +
	"\x01\xda\xc7\x4a\x97\xcc\x18\xd4\xfe\x29\xee\x7d\x82\xfe\x0e\x60\xbb\xa5\x7b\x27\xcd\xde\x7f\x59\x29\x7c\xe2\xaf" +
	"\x52\xbf\x90\xa3\x81\x3c\x58\x5e\x44\x74\xe1\x45\x26\xd7\x07\x76\x51\x66\xe9\xb1\x60\x45\x1d\x13\x61\x94\xd8\xd1" +
	"\xa8\x09\xe4\x4c\xd4\xe8\x29\xb6\x7b\x73\x1c\x11\x06\x97\x8a\x4b\xf3\xff\xff\x8d\x27\x3a\x25\xd3\x48\xa6\xff\xbc" +
	"\x3d\x4b\x2e\x14\x7b\x31\xcf\xb1\x33\x8e\x65\x4a\x5f\x97\x6b\x2c\x6a\x77\x5d\xd9\x09\x1a\x45\xbf\x8c\xd9\x6b\xd9" +
	"\x96\x82\x11\x1f\xae\xd0\x1c\xae\x6a\xa3\x4a\x97\xa8\x9f\x32\xcf\x61\x98\xfe\x84\xd5\xfe\xd1\x13\xda\x7f\x60\x3e" +
	"\x71\x99\x7d\xeb\xb7\x76\xee\x49\xc7\xa2\x96\x0e\xe9\xfa\xd0\xb7\xe1\x5f\x01\x00\x00\xff\xff\xb2\x5e\x21\xa3\xa3" +
	"\x0f\x00\x00")

func bindataTemplatesCommandstmplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTemplatesCommandstmpl,
		"templates/commands.tmpl",
	)
}



func bindataTemplatesCommandstmpl() (*asset, error) {
	bytes, err := bindataTemplatesCommandstmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "templates/commands.tmpl",
		size: 4003,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1539148995, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataTemplatesMaintmpl = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8e\x41\x4b\xc3\x40\x10\x85\xcf\x3b\xbf\x62\xc8\x41\x12\xb0\x9b\x7a" +
	"\xed\xad\x68\x0e\x5e\xac\x88\x78\x5f\x37\x93\xed\x60\x76\x66\xd9\x6c\x4a\x25\xe4\xbf\x4b\x2a\xe2\xed\xbd\xf7\xf1" +
	"\x1e\x2f\x39\xff\xe5\x02\x61\x74\x2c\x00\x1c\x93\xe6\x82\x35\x98\x2a\x70\x39\xcf\x9f\xd6\x6b\x6c\x7b\x27\x4c\x63" +
	"\x28\xee\x7b\xd4\xdc\x6a\x22\x71\x89\x77\x7e\xe4\x5d\x20\xa1\xec\x8a\xe6\xd6\x8f\x5c\x41\x03\x30\xcc\xe2\x6f\x63" +
	"\x75\x83\x0b\x18\x3f\xb2\x7d\x16\x2e\xf5\xdd\xa6\x1e\x55\x06\x0e\x0b\x18\x73\x4c\xe9\xc5\x45\x3a\x20\x62\xb5\x2c" +
	"\x68\x37\x83\xeb\x5a\xdd\x83\x31\x9d\x5c\x5e\x33\x0d\x7c\x3d\xfc\xb3\x4e\x2e\x7f\xf8\x83\xf2\xc4\x2a\xb7\xea\x83" +
	"\xdd\xdb\xfd\x96\xae\x0d\x80\x69\x5b\x7c\x3f\x3d\x9d\x0e\x78\xec\x7b\xcc\x14\x78\x2a\x94\xd1\x6b\x8c\x4e\xfa\x09" +
	"\xcf\x94\xc9\xc2\xef\xa7\x37\xd5\x62\xbb\x2b\xf9\xb9\x50\xdd\xc0\x0a\x3f\x01\x00\x00\xff\xff\xd7\x90\x9c\xb4\x08" +
	"\x01\x00\x00")

func bindataTemplatesMaintmplBytes() ([]byte, error) {
	return bindataRead(
		_bindataTemplatesMaintmpl,
		"templates/main.tmpl",
	)
}



func bindataTemplatesMaintmpl() (*asset, error) {
	bytes, err := bindataTemplatesMaintmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "templates/main.tmpl",
		size: 264,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1538448749, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}


//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"templates/commands.tmpl": bindataTemplatesCommandstmpl,
	"templates/main.tmpl":     bindataTemplatesMaintmpl,
}

//
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}


type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"templates": {Func: nil, Children: map[string]*bintree{
		"commands.tmpl": {Func: bindataTemplatesCommandstmpl, Children: map[string]*bintree{}},
		"main.tmpl": {Func: bindataTemplatesMaintmpl, Children: map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}