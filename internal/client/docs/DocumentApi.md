# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ChangeDocPermissions**](DocumentApi.md#ChangeDocPermissions) | **Put** /docs/{docId}/perms | 
[**DeleteDocument**](DocumentApi.md#DeleteDocument) | **Delete** /docs/{docId}/delete | 
[**DownloadDocument**](DocumentApi.md#DownloadDocument) | **Get** /docs/{docId}/download | 
[**GetDocPermissions**](DocumentApi.md#GetDocPermissions) | **Get** /docs/{docId}/perms | 
[**GetDocument**](DocumentApi.md#GetDocument) | **Get** /docs/{docId} | 
[**ListRoleDocuments**](DocumentApi.md#ListRoleDocuments) | **Get** /roles/{role}/docs | 
[**UploadDocument**](DocumentApi.md#UploadDocument) | **Post** /docs/upload | 

# **ChangeDocPermissions**
> ChangeDocPermissions(ctx, body, docId)


Change doc permissions

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**SharePermissions**](SharePermissions.md)|  | 
  **docId** | **string**|  | 

### Return type

 (empty response body)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteDocument**
> DeleteDocument(ctx, docId)


Only the owner may perform this action on a document

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **docId** | **string**| Identifier of the doc to delete | 

### Return type

 (empty response body)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DownloadDocument**
> *os.File DownloadDocument(ctx, docId)


Download a doc if the user role matches the doc permissions. Must have appropriate roles or user log in to access.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **docId** | **string**| The doc that needs to be fetched. | 

### Return type

[***os.File**](*os.File.md)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDocPermissions**
> SharePermissions GetDocPermissions(ctx, docId)


Get doc permissions

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **docId** | **string**| Identifier of the doc | 

### Return type

[**SharePermissions**](SharePermissions.md)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDocument**
> Document GetDocument(ctx, docId)


Get doc info. Must have appropriate roles or user log in to access.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **docId** | **string**| The doc that needs to be fetched. | 

### Return type

[**Document**](Document.md)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListRoleDocuments**
> []Document ListRoleDocuments(ctx, role, optional)


List the docs of a specific role

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **role** | **string**| Docs accessible by a certain role | 
 **optional** | ***DocumentApiListRoleDocumentsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DocumentApiListRoleDocumentsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **optional.Int32**| The page number to retrieve. | [default to 1]
 **size** | **optional.Int32**| The number of users to retrieve per page. | [default to 10]

### Return type

[**[]Document**](Document.md)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UploadDocument**
> UploadDocument(ctx, doc, docName, xHash)


Upload a doc to your account

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **doc** | ***os.File*****os.File**|  | 
  **docName** | **string**| Name of the doc to display in the app | 
  **xHash** | **string**| Hash value of the original doc to verify integrity | 

### Return type

 (empty response body)

### Authorization

[token](../README.md#token)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

