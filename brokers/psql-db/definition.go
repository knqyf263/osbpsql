package psqldb

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets138f62064cf7da2f3bb42f0737f399a4e88eb299 = "id: \"6cb6dbdb-e4d2-4e37-aa91-48529a4e18b4\"\nname: \"psql-database\"\ndescription: \"Create database\"\nbindable: true\ntags: \n  - postgresql\nplanupdatable: false\nmetadata: \n  displayname: \"PostgreSQL Database\"\n  documentationurl: https://github.com/knqyf263/osbpsql\nplans: \n  - id: 53296ba3-38e6-4a49-acf4-c0cd75a88a91\n    name: standard\n    description: Standard\n    free: true\n    bindable: true"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"definition.yaml"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1537442790, 1537442790000000000),
		Data:     nil,
	}, "/definition.yaml": &assets.File{
		Path:     "/definition.yaml",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1537284871, 1537284871000000000),
		Data:     []byte(_Assets138f62064cf7da2f3bb42f0737f399a4e88eb299),
	}}, "")
