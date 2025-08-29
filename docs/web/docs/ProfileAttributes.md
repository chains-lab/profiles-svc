# ProfileAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | Username | 
**Pseudonym** | Pointer to **string** | Pseudonym | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Avatar** | Pointer to **string** | Avatar URL | [optional] 
**Official** | **bool** | Is Official Account | 
**Sex** | Pointer to **string** | sex prikin&#39; | [optional] 
**Birthdate** | Pointer to **time.Time** | Birthday | [optional] 
**UpdatedAt** | **time.Time** | Updated At | 
**CreatedAt** | **time.Time** | Created At | 

## Methods

### NewProfileAttributes

`func NewProfileAttributes(username string, official bool, updatedAt time.Time, createdAt time.Time, ) *ProfileAttributes`

NewProfileAttributes instantiates a new ProfileAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProfileAttributesWithDefaults

`func NewProfileAttributesWithDefaults() *ProfileAttributes`

NewProfileAttributesWithDefaults instantiates a new ProfileAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *ProfileAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ProfileAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ProfileAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPseudonym

`func (o *ProfileAttributes) GetPseudonym() string`

GetPseudonym returns the Pseudonym field if non-nil, zero value otherwise.

### GetPseudonymOk

`func (o *ProfileAttributes) GetPseudonymOk() (*string, bool)`

GetPseudonymOk returns a tuple with the Pseudonym field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPseudonym

`func (o *ProfileAttributes) SetPseudonym(v string)`

SetPseudonym sets Pseudonym field to given value.

### HasPseudonym

`func (o *ProfileAttributes) HasPseudonym() bool`

HasPseudonym returns a boolean if a field has been set.

### GetDescription

`func (o *ProfileAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ProfileAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ProfileAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ProfileAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAvatar

`func (o *ProfileAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *ProfileAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *ProfileAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *ProfileAttributes) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetOfficial

`func (o *ProfileAttributes) GetOfficial() bool`

GetOfficial returns the Official field if non-nil, zero value otherwise.

### GetOfficialOk

`func (o *ProfileAttributes) GetOfficialOk() (*bool, bool)`

GetOfficialOk returns a tuple with the Official field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfficial

`func (o *ProfileAttributes) SetOfficial(v bool)`

SetOfficial sets Official field to given value.


### GetSex

`func (o *ProfileAttributes) GetSex() string`

GetSex returns the Sex field if non-nil, zero value otherwise.

### GetSexOk

`func (o *ProfileAttributes) GetSexOk() (*string, bool)`

GetSexOk returns a tuple with the Sex field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSex

`func (o *ProfileAttributes) SetSex(v string)`

SetSex sets Sex field to given value.

### HasSex

`func (o *ProfileAttributes) HasSex() bool`

HasSex returns a boolean if a field has been set.

### GetBirthdate

`func (o *ProfileAttributes) GetBirthdate() time.Time`

GetBirthdate returns the Birthdate field if non-nil, zero value otherwise.

### GetBirthdateOk

`func (o *ProfileAttributes) GetBirthdateOk() (*time.Time, bool)`

GetBirthdateOk returns a tuple with the Birthdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBirthdate

`func (o *ProfileAttributes) SetBirthdate(v time.Time)`

SetBirthdate sets Birthdate field to given value.

### HasBirthdate

`func (o *ProfileAttributes) HasBirthdate() bool`

HasBirthdate returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ProfileAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ProfileAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ProfileAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetCreatedAt

`func (o *ProfileAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ProfileAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ProfileAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


