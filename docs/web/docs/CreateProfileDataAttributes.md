# CreateProfileDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | Username | 
**Pseudonym** | Pointer to **string** | Pseudonym | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Avatar** | Pointer to **string** | Avatar URL | [optional] 
**Sex** | Pointer to **string** | sex prikin&#39; | [optional] 
**Birthdate** | Pointer to **time.Time** | Birthday | [optional] 

## Methods

### NewCreateProfileDataAttributes

`func NewCreateProfileDataAttributes(username string, ) *CreateProfileDataAttributes`

NewCreateProfileDataAttributes instantiates a new CreateProfileDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateProfileDataAttributesWithDefaults

`func NewCreateProfileDataAttributesWithDefaults() *CreateProfileDataAttributes`

NewCreateProfileDataAttributesWithDefaults instantiates a new CreateProfileDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *CreateProfileDataAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *CreateProfileDataAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *CreateProfileDataAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPseudonym

`func (o *CreateProfileDataAttributes) GetPseudonym() string`

GetPseudonym returns the Pseudonym field if non-nil, zero value otherwise.

### GetPseudonymOk

`func (o *CreateProfileDataAttributes) GetPseudonymOk() (*string, bool)`

GetPseudonymOk returns a tuple with the Pseudonym field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPseudonym

`func (o *CreateProfileDataAttributes) SetPseudonym(v string)`

SetPseudonym sets Pseudonym field to given value.

### HasPseudonym

`func (o *CreateProfileDataAttributes) HasPseudonym() bool`

HasPseudonym returns a boolean if a field has been set.

### GetDescription

`func (o *CreateProfileDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *CreateProfileDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *CreateProfileDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *CreateProfileDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAvatar

`func (o *CreateProfileDataAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *CreateProfileDataAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *CreateProfileDataAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *CreateProfileDataAttributes) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetSex

`func (o *CreateProfileDataAttributes) GetSex() string`

GetSex returns the Sex field if non-nil, zero value otherwise.

### GetSexOk

`func (o *CreateProfileDataAttributes) GetSexOk() (*string, bool)`

GetSexOk returns a tuple with the Sex field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSex

`func (o *CreateProfileDataAttributes) SetSex(v string)`

SetSex sets Sex field to given value.

### HasSex

`func (o *CreateProfileDataAttributes) HasSex() bool`

HasSex returns a boolean if a field has been set.

### GetBirthdate

`func (o *CreateProfileDataAttributes) GetBirthdate() time.Time`

GetBirthdate returns the Birthdate field if non-nil, zero value otherwise.

### GetBirthdateOk

`func (o *CreateProfileDataAttributes) GetBirthdateOk() (*time.Time, bool)`

GetBirthdateOk returns a tuple with the Birthdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBirthdate

`func (o *CreateProfileDataAttributes) SetBirthdate(v time.Time)`

SetBirthdate sets Birthdate field to given value.

### HasBirthdate

`func (o *CreateProfileDataAttributes) HasBirthdate() bool`

HasBirthdate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


