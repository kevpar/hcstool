/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type ModifySettingRequest struct {
	ResourcePath string `json:"ResourcePath,omitempty"`

	RequestType string `json:"RequestType,omitempty"` // NOTE: Swagger generated as string. Locally updated.

	Settings interface{} `json:"Settings,omitempty"` // NOTE: Swagger generated as *interface{}. Locally updated

	GuestRequest interface{} `json:"GuestRequest,omitempty"` // NOTE: Swagger generated as *interface{}. Locally updated
}