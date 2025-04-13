package main

import (
	"bufio"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"varuna-openapi/internal/client"
	"varuna-openapi/internal/server/db"

	"github.com/antihax/optional"
	"golang.org/x/term"
)

const USERAPI_HELP = `Options:
	0 - Create users
	1 - Login
	2 - List users
	3 - Get users by name
	4 - List users by role
	5 - Update user password`

const DOCAPI_HELP = `Options:
	0 - Upload
	1 - Download
	2 - Delete
	3 - List documents by role
	4 - Get doc permissions
	5 - Change doc permissions
	6 - Get doc meta-data`

type UserOp int

const (
	CREATE_USER UserOp = iota
	LOGIN
	LIST_USERS
	GET_USER
	USERS_BY_ROLE
	PASSWD
)

type DocOp int

const (
	UPLOAD DocOp = iota
	DOWNLOAD
	DELETE
	DOC_ROLES
	GETMOD
	CHMOD
	GETDOC
)

var ctx context.Context

func handleUserApi(taskId int) {
	var err error
	cfg := client.NewConfiguration()
	apiClient := client.NewAPIClient(cfg)

	cmdReader := bufio.NewReader(os.Stdin)

	switch UserOp(taskId) {
	case CREATE_USER:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Email: ")
		email, _, _ := cmdReader.ReadLine()
		fmt.Print("Password: ")
		pass, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		req := client.RegisterReq{Email: string(email), Password: string(pass), Username: string(username)}
		apiClient.UsersApi.CreateUsers(ctx, req)
	case LOGIN:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Password: ")
		pass, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		req := client.LoginReq{Username: string(username), Password: string(pass)}
		login, _, err := apiClient.UsersApi.LoginPost(ctx, req)
		if err != nil {
			var swagErr *client.GenericSwaggerError = &client.GenericSwaggerError{}
			if errors.As(err, swagErr) {
				fmt.Println("Error:", swagErr.Error(), swagErr.Body())
				return
			}
			fmt.Println(err.Error())
			return
		}

		tokenContent, err := json.Marshal(login)
		if err != nil {
			fmt.Println("Logged in but got another error:", err.Error())
			return
		}

		err = os.WriteFile("token.json", tokenContent, 0600)
		if err != nil {
			fmt.Println("Save token error:", err.Error())
			return
		}
	case LIST_USERS:
		var page, count int
		for {
			fmt.Print("Items per page: ")
			countStr, _, _ := cmdReader.ReadLine()
			if len(countStr) == 0 {
				count = 10
				break
			}
			count, err = strconv.Atoi(string(countStr))
			if err != nil {
				break
			}
		}
		for {
			fmt.Print("Page: ")
			pageStr, _, _ := cmdReader.ReadLine()
			if len(pageStr) == 0 {
				page = 0
				break
			}
			page, err = strconv.Atoi(string(pageStr))
			if err != nil {
				break
			}
		}
		req := &client.UsersApiListUsersOpts{
			Page: optional.NewInt32(int32(page)),
			Size: optional.NewInt32(int32(count))}
		apiClient.UsersApi.ListUsers(ctx, req)
	case GET_USER:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		apiClient.UsersApi.GetUserByName(ctx, string(username))
	case USERS_BY_ROLE:
		fmt.Print("Role search: ")
		role, _, _ := cmdReader.ReadLine()
		var page, count int
		for {
			fmt.Print("Items per page: ")
			countStr, _, _ := cmdReader.ReadLine()
			if len(countStr) == 0 {
				count = 10
				break
			}
			count, err = strconv.Atoi(string(countStr))
			if err != nil {
				break
			}
		}
		for {
			fmt.Print("Page: ")
			pageStr, _, _ := cmdReader.ReadLine()
			if len(pageStr) == 0 {
				page = 0
				break
			}
			page, err = strconv.Atoi(string(pageStr))
			if err != nil {
				break
			}
		}
		req := &client.UsersApiListUsersByRoleOpts{
			Page: optional.NewInt32(int32(page)),
			Size: optional.NewInt32(int32(count)),
		}
		apiClient.UsersApi.ListUsersByRole(ctx, string(role), req)
	case PASSWD:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Prev password: ")
		prevPass, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		fmt.Print("New password: ")
		fmt.Println()
		newPass, _ := term.ReadPassword(int(syscall.Stdin))
		body := client.LoginChangeReq{
			PrevPassword: string(prevPass),
			NewPassword:  string(newPass),
		}
		apiClient.UsersApi.UpdateUser(ctx, body, string(username))
	default:
		fmt.Println(USERAPI_HELP)
	}
}

func handleDocApi(taskId int) {
	var err error
	cfg := client.NewConfiguration()
	apiClient := client.NewAPIClient(cfg)

	cmdReader := bufio.NewReader(os.Stdin)

	switch DocOp(taskId) {
	case UPLOAD:
		fmt.Print("File path: ")
		fileMsg, _, err := cmdReader.ReadLine()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
		// fileMsg := []byte("go.mod")
		file, err := os.Open(string(fileMsg))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()
		hasher := sha512.New()
		if _, err := io.Copy(hasher, file); err != nil {
			log.Fatal(err)
		}

		hash := hasher.Sum(nil)
		hash64 := base64.StdEncoding.EncodeToString(hash)

		fmt.Print("Doc name: ")
		docName, _, _ := cmdReader.ReadLine()

		// docName := []byte("go.mod")
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		for i, v := range db.Roles {
			fmt.Printf("\t%v: %v\n", i, v)
		}
		fmt.Print("Select roles by index separated by commas: ")
		rolesStr, _, _ := cmdReader.ReadLine()

		rolesIdx := strings.Split(string(rolesStr), ",")
		roles := ""
		for _, v := range rolesIdx {
			idx, err := strconv.Atoi(strings.Trim(v, " "))
			if err != nil {
				log.Fatal(err)
			}
			roles += string(db.Roles[idx]) + ","
		}

		roles = roles[:len(roles)-1]

		_, err = apiClient.DocumentApi.UploadDocument(ctx, file, hash64, string(docName), roles)
		if err != nil {
			var swagErr *client.GenericSwaggerError = &client.GenericSwaggerError{}
			if errors.As(err, swagErr) {
				fmt.Println("Error:", swagErr.Error(), string(swagErr.Body()))
				return
			}
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Uploaded document successfully")
	case DOWNLOAD:
		var id string
		for {
			fmt.Print("Doc id: ")
			id, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				continue
			}
			_, err = strconv.Atoi(string(id))
			if err != nil {
				break
			}
		}
		apiClient.DocumentApi.DownloadDocument(ctx, id)
	case DOC_ROLES:
		for i, v := range db.Roles {
			fmt.Printf("\t%v: %v\n", i, v)
		}
		fmt.Print("Role: ")
		idxStr, _, _ := cmdReader.ReadLine()
		idx, err := strconv.Atoi(strings.Trim(string(idxStr), " "))
		if err != nil {
			log.Fatal(err)
		}
		role := db.Roles[idx]
		var page, count int
		for {
			fmt.Print("Items per page: ")
			countStr, _, _ := cmdReader.ReadLine()
			if len(countStr) == 0 {
				count = 10
				break
			}
			count, err = strconv.Atoi(string(countStr))
			if err == nil {
				break
			}
			fmt.Println(err)
		}
		for {
			fmt.Print("Page: ")
			pageStr, _, _ := cmdReader.ReadLine()
			if len(pageStr) == 0 {
				page = 0
				break
			}
			page, err = strconv.Atoi(string(pageStr))
			if err == nil {
				break
			}
			fmt.Println(err)
		}
		opts := &client.DocumentApiListRoleDocumentsOpts{
			Page: optional.NewInt32(int32(page)),
			Size: optional.NewInt32(int32(count)),
		}
		roleFiles, _, err := apiClient.DocumentApi.ListRoleDocuments(ctx, string(role), opts)
		if err != nil {
			var swagErr *client.GenericSwaggerError = &client.GenericSwaggerError{}
			if errors.As(err, swagErr) {
				fmt.Println("Error:", swagErr.Error(), string(swagErr.Body()))
				return
			}
			fmt.Println(err.Error())
			return
		}
		for _, v := range roleFiles {
			fmt.Printf("%v:\n\tname: %v\n\tdate:%v\n", v.DocId, v.DocName, v.CreationDate)
			fmt.Print("\troles: ")
			l := len(v.Permissions.Roles)
			roles := ""
			for i, role := range v.Permissions.Roles {
				roles += role
				if i+1 != l {
					roles += ", "
				}
			}
			fmt.Println(roles)
		}
	case DELETE:
		var id string
		for {
			fmt.Print("Doc id: ")
			id, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				continue
			}
			_, err = strconv.Atoi(string(id))
			if err != nil {
				break
			}
		}
		apiClient.DocumentApi.DeleteDocument(ctx, id)
	case CHMOD:
		var id string
		for {
			fmt.Print("Doc id: ")
			id, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				continue
			}
			_, err = strconv.Atoi(string(id))
			if err != nil {
				break
			}
		}

		roles := make([]string, 0)
		users := make([]string, 0)

		i := 0
		for {
			fmt.Printf("Role %v: ", i)
			role, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				break
			}
			roles = append(roles, string(role))
			i++
		}

		i = 0
		for {
			fmt.Printf("User %v: ", i)
			user, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				break
			}
			users = append(users, string(user))
			i++
		}
		req := client.SharePermissions{
			Users: users,
			Roles: roles,
		}
		apiClient.DocumentApi.ChangeDocPermissions(ctx, req, id)
	case GETMOD:
		var id string
		for {
			fmt.Print("Doc id: ")
			id, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				continue
			}
			_, err = strconv.Atoi(string(id))
			if err != nil {
				break
			}
		}
		perms, _, err := apiClient.DocumentApi.GetDocPermissions(ctx, id)
		if err != nil {
			var swagErr *client.GenericSwaggerError = &client.GenericSwaggerError{}
			if errors.As(err, swagErr) {
				fmt.Println("Error:", swagErr.Error(), swagErr.Body())
				return
			}
			fmt.Println(err.Error())
			return
		}
		jason, _ := json.Marshal(perms)
		fmt.Println(string(jason))
	case GETDOC:
		var id string
		for {
			fmt.Print("Doc id: ")
			id, _, _ := cmdReader.ReadLine()
			if len(id) == 0 {
				continue
			}
			_, err = strconv.Atoi(string(id))
			if err != nil {
				break
			}
		}
		apiClient.DocumentApi.GetDocument(ctx, id)
	default:
		fmt.Println(DOCAPI_HELP)
	}
}

func main() {
	ctx = context.Background()
	tokenBytes, err := os.ReadFile("token.json")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println("Error:", err)
			return
		}
	} else {
		var login *client.LoginResp = &client.LoginResp{}
		err = json.Unmarshal(tokenBytes, login)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		ctx = context.WithValue(ctx, client.ContextAccessToken, login.Token)
	}
	serviceInt := flag.Int("service", -1, "0: Users API\n1: Documents API")
	taskInt := flag.Int("task", -1, "Task id for selected API")

	flag.Parse()

	switch *serviceInt {
	case 0:
		handleUserApi(*taskInt)
	case 1:
		handleDocApi(*taskInt)
	default:
		fmt.Println("0: Users API\n1: Documents API")
		os.Exit(1)
	}

}
