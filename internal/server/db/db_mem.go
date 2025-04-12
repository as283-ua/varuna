package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"time"
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
	RoleKeys  map[Role][]byte `json:"roles"` // role -> encrypted key with role key
	CreatedAt time.Time       `json:"createdAt"`
}

type DataBase struct {
	Users     map[string]User   `json:"users"`
	Files     []File            `json:"files"`
	RoleFiles map[Role][]int    `json:"roleFiles"`
	RoleUsers map[Role][]string `json:"roleUsers"`
	UserFiles map[string][]int  `json:"userFiles"`
	RoleKeys  map[Role][]byte   `json:"roleKeys"`
}

func (db *DataBase) AddRole(role Role) {
	db.RoleFiles[role] = make([]int, 0)
	db.RoleUsers[role] = make([]string, 0)
	db.RoleKeys[role] = make([]byte, 32)
	rand.Read(db.RoleKeys[role])
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
		file.RoleKeys[v], err = Encrypt(key, DB.RoleKeys[v])
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

func Encrypt(data, key []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	out = append(nonce, ciphertext...)
	return out, nil
}

func Decrypt(data, key []byte) (out []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	out, err = gcm.Open(nil, nonce, ciphertext, nil)
	return
}

var DB DataBase

func init() {
	DB = DataBase{
		Users:     make(map[string]User),
		Files:     make([]File, 0),
		RoleFiles: make(map[Role][]int),
		RoleUsers: make(map[Role][]string),
		UserFiles: make(map[string][]int),
		RoleKeys:  make(map[Role][]byte, 0),
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
