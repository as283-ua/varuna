/*
 * Varuna Docs
 *
 * PI Project
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

type LoginChangeReq struct {
	PrevPassword string `json:"prev-password,omitempty"`
	NewPassword  string `json:"new-password,omitempty"`
}
