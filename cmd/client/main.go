package main

import (
	"bufio"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"varuna-openapi/internal/client"

	"github.com/antihax/optional"
	"github.com/chzyer/readline"
	"golang.org/x/term"
)

const USERAPI_HELP = `Options:
	0 - Create users
	1 - Login
	2 - List users
	3 - Get users by name
	4 - List users by role
	5 - Update user`

const DOCAPI_HELP = `Options:
	0 - Upload
	1 - Download
	2 - Delete
	3 - List documents by role
	4 - Get doc permissions
	5 - Change doc permissions
	6 - Get doc meta-data`

func handleUserApi(taskId int) {
	var err error
	cfg := client.NewConfiguration()
	apiClient := client.NewAPIClient(cfg)
	ctx := context.Background()

	cmdReader := bufio.NewReader(os.Stdin)

	switch taskId {
	case 0:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Email: ")
		email, _, _ := cmdReader.ReadLine()
		fmt.Print("Password: ")
		pass, _ := term.ReadPassword(int(syscall.Stdin))
		req := client.RegisterReq{Email: string(email), Password: string(pass), Username: string(username)}
		apiClient.UsersApi.CreateUsers(ctx, req)
	case 1:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Password: ")
		pass, _ := term.ReadPassword(int(syscall.Stdin))
		req := client.LoginReq{Username: string(username), Password: string(pass)}
		apiClient.UsersApi.LoginPost(ctx, req)
	case 2:
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
	case 3:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		apiClient.UsersApi.GetUserByName(ctx, string(username))
	case 4:
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
	case 5:
		fmt.Print("Username: ")
		username, _, _ := cmdReader.ReadLine()
		fmt.Print("Prev password: ")
		prevPass, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Print("New password: ")
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
	ctx := context.Background()

	cmdReader := bufio.NewReader(os.Stdin)

	switch taskId {
	case 0:
		fileMsg, _ := readline.ReadMessage(cmdReader)
		file, err := os.Open(string(fileMsg.Data))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
		defer file.Close()
		hasher := sha512.New()
		if _, err := io.Copy(hasher, file); err != nil {
			panic(err)
		}

		fmt.Print("Doc name: ")
		docName, _, _ := cmdReader.ReadLine()

		hash := hasher.Sum(nil)
		hash64 := base64.StdEncoding.EncodeToString(hash)
		apiClient.DocumentApi.UploadDocument(ctx, file, hash64, string(docName))
	case 1:
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
	case 3:
		fmt.Print("Role: ")
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
		opts := &client.DocumentApiListRoleDocumentsOpts{
			Page: optional.NewInt32(int32(page)),
			Size: optional.NewInt32(int32(count)),
		}
		apiClient.DocumentApi.ListRoleDocuments(ctx, string(role), opts)
	case 2:
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
	case 5:
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
	case 4:
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
		apiClient.DocumentApi.GetDocPermissions(ctx, id)
	case 6:
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
