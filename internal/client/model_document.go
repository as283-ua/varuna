/*
 * Varuna Docs
 *
 * PI Project
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

type Document struct {
	DocId        int64             `json:"docId,omitempty"`
	DocName      string            `json:"docName,omitempty"`
	Hash         string            `json:"hash,omitempty"`
	Description  string            `json:"description,omitempty"`
	CreationDate string            `json:"creationDate,omitempty"`
	Permissions  *SharePermissions `json:"permissions,omitempty"`
}
