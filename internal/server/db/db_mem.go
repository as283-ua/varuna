package db

import (
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
	Username string   `json:"username"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

type File struct {
	Name      string    `json:"name"`
	Owner     string    `json:"owner,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
}

type DataBase struct {
	Users     map[string]User   `json:"users"`
	Files     []File            `json:"files"`
	RoleFiles map[Role][]int    `json:"roleFiles"`
	RoleUsers map[Role][]string `json:"roleUsers"`
	UserFiles map[string][]int  `json:"userFiles"`
}

func (db *DataBase) AddRole(role Role) {
	db.RoleFiles[role] = make([]int, 0)
	db.RoleUsers[role] = make([]string, 0)
}

func (db *DataBase) AddUser(user User) {
	db.Users[user.Username] = user
	db.UserFiles[user.Username] = make([]int, 0)
	for _, v := range user.Roles {
		db.RoleUsers[Role(v)] = append(db.RoleUsers[Role(v)], user.Username)
	}
}

func (db *DataBase) AddFile(file File) {
	db.Files = append(db.Files, file)
	fileIdx := len(db.Files) - 1
	db.UserFiles[file.Owner] = append(db.UserFiles[file.Owner], fileIdx)
	for _, v := range file.Roles {
		db.RoleFiles[Role(v)] = append(db.RoleFiles[Role(v)], fileIdx)
	}
}

var DB DataBase

func init() {
	DB = DataBase{
		Users:     make(map[string]User),
		Files:     make([]File, 0),
		RoleFiles: make(map[Role][]int),
		RoleUsers: make(map[Role][]string),
		UserFiles: make(map[string][]int),
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
		Roles:    []string{string(RoleSoftware), string(RoleQA)},
	})
	DB.AddUser(User{
		Username: "dlc5",
		Email:    "correo2@alu.ua.es",
		Password: "password",
		Roles:    []string{string(RoleDevops)},
	})
	DB.AddUser(User{
		Username: "aic32",
		Email:    "correo3@alu.ua.es",
		Password: "password",
		Roles:    []string{string(RoleFinance), string(RoleSoftware)},
	})
	DB.AddUser(User{
		Username: "rafica",
		Email:    "correo4@alu.ua.es",
		Password: "password",
		Roles:    []string{string(RoleAdmin)},
	})

	DB.AddFile(File{
		Name:      "Archivo sw",
		Owner:     "as283",
		Roles:     []string{string(RoleSoftware), string(RoleQA)},
		CreatedAt: time.Now(),
	})
	DB.AddFile(File{
		Name:      "Archivo devops",
		Owner:     "dlc5",
		Roles:     []string{string(RoleDevops)},
		CreatedAt: time.Now(),
	})
	DB.AddFile(File{
		Name:      "Archivo finanza",
		Owner:     "aic32",
		Roles:     []string{string(RoleFinance)},
		CreatedAt: time.Now(),
	})
	DB.AddFile(File{
		Name:      "Importante pa todos",
		Owner:     "aic32",
		Roles:     []string{string(RoleFinance), string(RoleDevops), string(RoleSoftware)},
		CreatedAt: time.Now(),
	})
}
