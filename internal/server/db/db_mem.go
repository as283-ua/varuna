package db

import (
	"bytes"
	"crypto/rand"
	"log"
	"os"
	"time"
	"varuna-openapi/internal/server/util"
)

type Role string

const (
	RoleSoftware Role = "software"
	RoleHardware Role = "hardware"
	RoleDevops   Role = "devops"
	RoleHR       Role = "hr"
	RoleFinance  Role = "finance"
	RoleQA       Role = "qa"
	RoleAdmin    Role = "admin"
)

var Roles = []Role{
	RoleSoftware,
	RoleHardware,
	RoleDevops,
	RoleHR,
	RoleFinance,
	RoleQA,
	RoleAdmin,
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`
}

type File struct {
	Name      string          `json:"name"`
	Owner     string          `json:"owner,omitempty"`
	Roles     []Role          `json:"roles"`
	RoleKeys  map[Role][]byte `json:"rolekeys"` // role -> encrypted key with role key
	CreatedAt time.Time       `json:"createdAt"`
}

type DataBase struct {
	Users     map[string]User   `json:"users"`
	Files     []File            `json:"files"`
	RoleFiles map[Role][]int    `json:"roleFiles"`
	RoleUsers map[Role][]string `json:"roleUsers"`
	UserFiles map[string][]int  `json:"userFiles"`
	RoleKeys  map[Role][]byte   `json:"roleKeys"`
	KMS       map[Role][]byte   `json:"kms"`
}

func (db *DataBase) AddRole(role Role) {
	db.RoleFiles[role] = make([]int, 0)
	db.RoleUsers[role] = make([]string, 0)
	db.RoleKeys[role] = make([]byte, 32)
	rand.Read(db.RoleKeys[role])

	b := make([]byte, 32)
	rand.Read(b)
	db.KMS[role] = b
}

func (db *DataBase) AddUser(user User) {
	db.Users[user.Username] = user
	db.UserFiles[user.Username] = make([]int, 0)
	for _, v := range user.Roles {
		db.RoleUsers[Role(v)] = append(db.RoleUsers[Role(v)], user.Username)
	}
}

func (db *DataBase) AddFile(file File, key []byte) {
	file.RoleKeys = make(map[Role][]byte)
	var err error
	for _, v := range file.Roles {
		file.RoleKeys[v], err = util.Encrypt(key, DB.RoleKeys[v])
		if err != nil {
			log.Fatal(err)
		}
	}
	db.Files = append(db.Files, file)
	fileIdx := len(db.Files) - 1
	db.UserFiles[file.Owner] = append(db.UserFiles[file.Owner], fileIdx)
	for _, v := range file.Roles {
		db.RoleFiles[Role(v)] = append(db.RoleFiles[Role(v)], fileIdx)
	}
}

var DB DataBase

const DB_FILE = "varuna.db"

func ImportDb(key []byte) error {
	dbEnc, err := os.ReadFile(DB_FILE)
	if err != nil {
		return err
	}
	dbCompressed, err := util.Decrypt(dbEnc, key)
	if err != nil {
		return err
	}
	dbJson := util.Decompress(dbCompressed)
	err = util.DecodeJSON(bytes.NewReader(dbJson), &DB)
	if err != nil {
		return err
	}

	return nil
}

func ExportDb(key []byte) error {
	dbJson := util.EncodeJSON(DB)
	dbCompressed := util.Compress(dbJson)
	dbEnc, err := util.Encrypt(dbCompressed, key)
	if err != nil {
		return err
	}

	err = os.WriteFile(DB_FILE, dbEnc, 0600)

	return err
}

func CleanDb() {
	DB = DataBase{
		Users:     make(map[string]User),
		Files:     make([]File, 0),
		RoleFiles: make(map[Role][]int),
		RoleUsers: make(map[Role][]string),
		UserFiles: make(map[string][]int),
		RoleKeys:  make(map[Role][]byte, 0),
		KMS:       make(map[Role][]byte),
	}

	DB.AddRole(RoleSoftware)
	DB.AddRole(RoleDevops)
	DB.AddRole(RoleHardware)
	DB.AddRole(RoleQA)
	DB.AddRole(RoleFinance)
	DB.AddRole(RoleHR)

	DB.AddUser(User{
		Username: "as283",
		Email:    "correo1@alu.ua.es",
		Password: "password",
		Roles:    []Role{RoleSoftware, RoleQA},
	})
	DB.AddUser(User{
		Username: "dlc5",
		Email:    "correo2@alu.ua.es",
		Password: "password",
		Roles:    []Role{RoleDevops},
	})
	DB.AddUser(User{
		Username: "aic32",
		Email:    "correo3@alu.ua.es",
		Password: "password",
		Roles:    []Role{RoleFinance, RoleSoftware},
	})
	DB.AddUser(User{
		Username: "rafica",
		Email:    "correo4@alu.ua.es",
		Password: "password",
		Roles:    []Role{RoleAdmin},
	})

	key := make([]byte, 32)
	rand.Read(key)

	DB.AddFile(File{
		Name:      "Archivo sw",
		Owner:     "as283",
		Roles:     []Role{RoleSoftware, RoleQA},
		CreatedAt: time.Now(),
	}, key)

	key = make([]byte, 32)
	rand.Read(key)
	DB.AddFile(File{
		Name:      "Archivo devops",
		Owner:     "dlc5",
		Roles:     []Role{RoleDevops},
		CreatedAt: time.Now(),
	}, key)

	key = make([]byte, 32)
	rand.Read(key)
	DB.AddFile(File{
		Name:      "Archivo finanza",
		Owner:     "aic32",
		Roles:     []Role{RoleFinance},
		CreatedAt: time.Now(),
	}, key)

	key = make([]byte, 32)
	rand.Read(key)
	DB.AddFile(File{
		Name:      "Importante pa todos",
		Owner:     "aic32",
		Roles:     []Role{RoleFinance, RoleDevops, RoleSoftware},
		CreatedAt: time.Now(),
	}, key)
}
