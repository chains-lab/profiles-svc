# UpdateProfileDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Pseudonym** | Pointer to **string** | Pseudonym | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Avatar** | Pointer to **string** | Avatar URL | [optional] 
**Sex** | Pointer to **string** | sex prikin&#39; | [optional] 
**Birthdate** | Pointer to **time.Time** | Birthday | [optional] 

## Methods

### NewUpdateProfileDataAttributes

`func NewUpdateProfileDataAttributes() *UpdateProfileDataAttributes`

NewUpdateProfileDataAttributes instantiates a new UpdateProfileDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateProfileDataAttributesWithDefaults

`func NewUpdateProfileDataAttributesWithDefaults() *UpdateProfileDataAttributes`

NewUpdateProfileDataAttributesWithDefaults instantiates a new UpdateProfileDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPseudonym

`func (o *UpdateProfileDataAttributes) GetPseudonym() string`

GetPseudonym returns the Pseudonym field if non-nil, zero value otherwise.

### GetPseudonymOk

`func (o *UpdateProfileDataAttributes) GetPseudonymOk() (*string, bool)`

GetPseudonymOk returns a tuple with the Pseudonym field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPseudonym

`func (o *UpdateProfileDataAttributes) SetPseudonym(v string)`

SetPseudonym sets Pseudonym field to given value.

### HasPseudonym

`func (o *UpdateProfileDataAttributes) HasPseudonym() bool`

HasPseudonym returns a boolean if a field has been set.

### GetDescription

`func (o *UpdateProfileDataAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *UpdateProfileDataAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *UpdateProfileDataAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *UpdateProfileDataAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAvatar

`func (o *UpdateProfileDataAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *UpdateProfileDataAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *UpdateProfileDataAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *UpdateProfileDataAttributes) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetSex

`func (o *UpdateProfileDataAttributes) GetSex() string`

GetSex returns the Sex field if non-nil, zero value otherwise.

### GetSexOk

`func (o *UpdateProfileDataAttributes) GetSexOk() (*string, bool)`

GetSexOk returns a tuple with the Sex field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSex

`func (o *UpdateProfileDataAttributes) SetSex(v string)`

SetSex sets Sex field to given value.

### HasSex

`func (o *UpdateProfileDataAttributes) HasSex() bool`

HasSex returns a boolean if a field has been set.

### GetBirthdate

`func (o *UpdateProfileDataAttributes) GetBirthdate() time.Time`

GetBirthdate returns the Birthdate field if non-nil, zero value otherwise.

### GetBirthdateOk

`func (o *UpdateProfileDataAttributes) GetBirthdateOk() (*time.Time, bool)`

GetBirthdateOk returns a tuple with the Birthdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBirthdate

`func (o *UpdateProfileDataAttributes) SetBirthdate(v time.Time)`

SetBirthdate sets Birthdate field to given value.

### HasBirthdate

`func (o *UpdateProfileDataAttributes) HasBirthdate() bool`

HasBirthdate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


