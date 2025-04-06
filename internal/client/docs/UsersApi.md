# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateUsers**](UsersApi.md#CreateUsers) | **Post** /users | 
[**GetUserByName**](UsersApi.md#GetUserByName) | **Get** /users/{username} | 
[**ListUsers**](UsersApi.md#ListUsers) | **Get** /users | 
[**ListUsersByRole**](UsersApi.md#ListUsersByRole) | **Get** /roles/{role}/users | 
[**LoginPost**](UsersApi.md#LoginPost) | **Post** /login | 
[**RefreshPost**](UsersApi.md#RefreshPost) | **Post** /refresh | 
[**UpdateUser**](UsersApi.md#UpdateUser) | **Put** /users/{username}/creds | 

# **CreateUsers**
> InlineResponse201 CreateUsers(ctx, body)


Create Users. Only an admin may create a user for a new employee

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RegisterReq**](RegisterReq.md)| Fields for creating a new user | 

### Return type

[**InlineResponse201**](inline_response_201.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUserByName**
> UserPublic GetUserByName(ctx, username)


Primarily used to obtain the public key and share a doc

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **username** | **string**| The name that needs to be fetched. Use Users1 for testing.  | 

### Return type

[**UserPublic**](UserPublic.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUsers**
> UserPublic ListUsers(ctx, optional)


List all users except the one making the request

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiListUsersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListUsersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int32**| The page number to retrieve. | [default to 1]
 **size** | **optional.Int32**| The number of users to retrieve per page. | [default to 10]

### Return type

[**UserPublic**](UserPublic.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUsersByRole**
> []string ListUsersByRole(ctx, role, optional)


List all users except the one making the request that have the specified role

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **role** | **string**| User role | 
 **optional** | ***UsersApiListUsersByRoleOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListUsersByRoleOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **optional.Int32**| The page number to retrieve. | [default to 1]
 **size** | **optional.Int32**| The number of users to retrieve per page. | [default to 10]

### Return type

**[]string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LoginPost**
> LoginResp LoginPost(ctx, body)


Login

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**LoginReq**](LoginReq.md)| Fields for logging in, returns session token | 

### Return type

[**LoginResp**](LoginResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RefreshPost**
> LoginResp RefreshPost(ctx, body)


Use refresh token to get new access token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RefreshReq**](RefreshReq.md)| Fields for logging in, returns session token | 

### Return type

[**LoginResp**](LoginResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateUser**
> UserPublic UpdateUser(ctx, body, username)


Change user credentials

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**LoginChangeReq**](LoginChangeReq.md)| Fields for creating a new user | 
  **username** | **string**| The name that needs to be fetched. Use Users1 for testing.  | 

### Return type

[**UserPublic**](UserPublic.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

