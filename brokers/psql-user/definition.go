package psqluser

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets138f62064cf7da2f3bb42f0737f399a4e88eb299 = "id: \"17cfc1d5-a2ef-47bb-b3c2-b3fe07cb5b41\"\nname: \"psql-user\"\ndescription: \"Create user\"\nbindable: true\ntags: \n  - postgresql\nplanupdatable: false\nmetadata: \n  displayname: \"PostgreSQL User\"\n  documentationurl: https://github.com/knqyf263/osbpsql\nplans: \n  - id: 4cbdb95c-26a2-4968-aef1-a392ea23debc\n    name: standard\n    description: Standard\n    free: true\n    bindable: true"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"definition.yaml"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1537447501, 1537447501000000000),
		Data:     nil,
	}, "/definition.yaml": &assets.File{
		Path:     "/definition.yaml",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1537447577, 1537447577000000000),
		Data:     []byte(_Assets138f62064cf7da2f3bb42f0737f399a4e88eb299),
	}}, "")
