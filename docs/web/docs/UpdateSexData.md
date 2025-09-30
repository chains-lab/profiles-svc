# UpdateSexData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | user id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateSexDataAttributes**](UpdateSexDataAttributes.md) |  | 

## Methods

### NewUpdateSexData

`func NewUpdateSexData(id uuid.UUID, type_ string, attributes UpdateSexDataAttributes, ) *UpdateSexData`

NewUpdateSexData instantiates a new UpdateSexData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateSexDataWithDefaults

`func NewUpdateSexDataWithDefaults() *UpdateSexData`

NewUpdateSexDataWithDefaults instantiates a new UpdateSexData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateSexData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateSexData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateSexData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateSexData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateSexData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateSexData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateSexData) GetAttributes() UpdateSexDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateSexData) GetAttributesOk() (*UpdateSexDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateSexData) SetAttributes(v UpdateSexDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


