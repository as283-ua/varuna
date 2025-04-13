/*
 * Varuna Docs
 *
 * PI Project
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"varuna-openapi/internal/server/db"
	"varuna-openapi/internal/server/util"
)

func ChangeDocPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DownloadDocument(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Authorization header required"}`, http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusBadRequest)
		return
	}

	token := parts[1]
	user, ok := db.DB.Users[token]
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/docs/")
	parts = strings.SplitN(path, "/", 2)
	docId, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, `{"error": "Bad id"}`, http.StatusBadRequest)
		return
	}
	if docId >= len(db.DB.Files) || docId < 0 {
		http.Error(w, `{"error": "Bad id"}`, http.StatusBadRequest)
		return
	}
	fileMeta := db.DB.Files[docId]
	var userRole *db.Role = nil
	for _, v := range user.Roles {
		if slices.Contains(fileMeta.Roles, v) {
			userRole = &v
			break
		}
	}

	if userRole == nil {
		http.Error(w, `{"error": "Bad role"}`, http.StatusUnauthorized)
		return
	}

	roleKey := db.DB.KMS[*userRole]
	filename := fmt.Sprintf("%v_%v", fileMeta.Owner, fileMeta.Name)

	encFile, err := os.ReadFile("files/" + filename)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	encFileKey := fileMeta.RoleKeys[*userRole]
	fileKey, err := util.Decrypt(encFileKey, roleKey)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	fileData, err := util.Decrypt(encFile, fileKey)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileMeta.Name+"\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fileData)
}

func GetDocPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListRoleDocuments(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Authorization header required"}`, http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusBadRequest)
		return
	}

	token := parts[1]
	user, ok := db.DB.Users[token]
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/roles/")
	parts = strings.SplitN(path, "/", 2)
	if len(parts) < 2 || parts[1] != "docs" {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}
	role := db.Role(strings.ToLower(parts[0]))

	hasRole := false
	for _, r := range user.Roles {
		if r == role {
			hasRole = true
			break
		}
	}
	if !hasRole {
		http.Error(w, `{"error": "User does not have the requested role"}`, http.StatusBadRequest)
		return
	}

	var filesIdx []int
	if role == db.RoleAdmin {
		filesIdx = make([]int, len(db.DB.Files))
		for i := 0; i < len(filesIdx); i++ {
			filesIdx[i] = i
		}
	} else {
		filesIdx, ok = db.DB.RoleFiles[role]
		if !ok {
			http.Error(w, `{"error": "Role not found"}`, http.StatusBadRequest)
			return
		}
	}

	page := 1
	size := 10
	if p := req.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val >= 1 {
			page = val
		}
	}
	if s := req.URL.Query().Get("size"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val >= 1 {
			size = val
		}
	}

	start := (page - 1) * size
	if start >= len(filesIdx) {
		w.Header().Set("X-Total-Pages", "1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	end := start + size
	if end > len(filesIdx) {
		end = len(filesIdx)
	}

	result := make([]Document, 0, end-start)
	for _, idx := range filesIdx[start:end] {
		f := db.DB.Files[idx]

		strRoles := make([]string, len(f.Roles))
		for i, r := range f.Roles {
			strRoles[i] = string(r)
		}

		doc := Document{
			DocId:        int64(idx),
			DocName:      f.Name,
			CreationDate: f.CreatedAt.Format(time.RFC3339),
			Permissions: &SharePermissions{
				Roles: strRoles,
			},
		}
		result = append(result, doc)
	}

	totalPages := (len(filesIdx) + size - 1) / size
	w.Header().Set("X-Total-Pages", strconv.Itoa(totalPages))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Authorization header required"}`, http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusBadRequest)
		return
	}
	token := parts[1]
	user, ok := db.DB.Users[token]
	if !ok {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // limit ~10MB
	if err != nil {
		http.Error(w, `{"error": "Failed to parse multipart form"}`, http.StatusBadRequest)
		return
	}

	docName := r.URL.Query().Get("docName")
	if docName == "" {
		http.Error(w, `{"error": "Missing required query parameter: docName"}`, http.StatusBadRequest)
		return
	}

	hashHeader64 := r.Header.Get("X-Hash")
	if hashHeader64 == "" {
		http.Error(w, `{"error": "Missing required header: X-Hash"}`, http.StatusBadRequest)
		return
	}

	rolestr := r.URL.Query().Get("roles")
	if rolestr == "" {
		http.Error(w, `{"error": "Missing required query parameter: roles"}`, http.StatusBadRequest)
		return
	}
	roles := strings.Split(rolestr, ",")

	for _, r := range roles {
		if !slices.Contains(user.Roles, db.Role(r)) {
			http.Error(w, `{"error": "Invalid requested role for file '`+r+`'. User is not part of said role"}`, http.StatusBadRequest)
			return
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error": "Missing or invalid file field"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	hasher := sha512.New()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error": "Failed to read file content"}`, http.StatusBadRequest)
		return
	}
	hasher.Write(data)
	calculatedHash := hasher.Sum(nil)
	var hashHeader []byte = make([]byte, 64)
	base64.StdEncoding.Decode(hashHeader, []byte(hashHeader64))

	if !bytes.Equal(calculatedHash, hashHeader) {
		http.Error(w, `{"error": "Hash mismatch: file integrity check failed"}`, http.StatusBadRequest)
		return
	}

	key := make([]byte, 32)
	rand.Read(key)

	encData, _ := util.Encrypt(data, key)

	filePath := fmt.Sprintf("files/%s_%s", user.Username, docName)
	err = os.WriteFile(filePath, encData, 0644)
	if err != nil {
		http.Error(w, `{"error": "Failed to store file on server"}`, http.StatusInternalServerError)
		return
	}

	rolesR := make([]db.Role, len(roles))
	for i, v := range roles {
		rolesR[i] = db.Role(v)
	}
	log.Printf("Storing document '%s' for user '%s' (%d bytes)\n", docName, user.Username, len(data))

	db.DB.AddFile(db.File{
		Name:      docName,
		Owner:     token,
		Roles:     rolesR,
		CreatedAt: time.Now(),
	}, key)

	w.WriteHeader(http.StatusOK)
}
