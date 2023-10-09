package commons

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/organization"
)

const (
	OrgUserGroupResourceName ResourceName = "spotinst_organization_user_group"
)

var OrgUserGroupResource *OrgUserGroupTerraformResource

type OrgUserGroupTerraformResource struct {
	GenericResource
}

type OrgUserGroupWrapper struct {
	OrgUserGroup *organization.UserGroup
}

func NewOrgUserGroupResource(fieldsMap map[FieldName]*GenericField) *OrgUserGroupTerraformResource {
	return &OrgUserGroupTerraformResource{
		GenericResource: GenericResource{
			resourceName: OrgUserGroupResourceName,
			fields:       NewGenericFields(fieldsMap),
		},
	}
}

func (res *OrgUserGroupTerraformResource) OnCreate(
	resourceData *schema.ResourceData,
	meta interface{}) (*organization.UserGroup, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return nil, fmt.Errorf("resource fields are nil or empty, cannot create")
	}

	OrgUserGroupWrapper := NewOrgUserGroupWrapper()

	for _, field := range res.fields.fieldsMap {
		if field.onCreate == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnCreate), field.resourceAffinity, field.fieldNameStr)
		if err := field.onCreate(OrgUserGroupWrapper, resourceData, meta); err != nil {
			return nil, err
		}
	}
	return OrgUserGroupWrapper.GetOrgUserGroup(), nil
}

func (res *OrgUserGroupTerraformResource) OnRead(
	OrgUserGroup *organization.UserGroup,
	resourceData *schema.ResourceData,
	meta interface{}) error {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return fmt.Errorf("resource fields are nil or empty, cannot read")
	}

	OrgUserGroupWrapper := NewOrgUserGroupWrapper()
	OrgUserGroupWrapper.SetOrgUserGroup(OrgUserGroup)

	for _, field := range res.fields.fieldsMap {
		if field.onRead == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnRead), field.resourceAffinity, field.fieldNameStr)
		if err := field.onRead(OrgUserGroupWrapper, resourceData, meta); err != nil {
			return err
		}
	}
	return nil
}

func (res *OrgUserGroupTerraformResource) OnUpdate(
	resourceData *schema.ResourceData,
	meta interface{}) (bool, *organization.UserGroup, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return false, nil, fmt.Errorf("resource fields are nil or empty, cannot update")
	}

	OrgUserGroupWrapper := NewOrgUserGroupWrapper()
	hasChanged := false
	for _, field := range res.fields.fieldsMap {
		if field.onUpdate == nil {
			continue
		}
		if field.hasFieldChange(resourceData, meta) {
			log.Printf(string(ResourceFieldOnUpdate), field.resourceAffinity, field.fieldNameStr)
			if err := field.onUpdate(OrgUserGroupWrapper, resourceData, meta); err != nil {
				return false, nil, err
			}
			hasChanged = true
		}
	}

	return hasChanged, OrgUserGroupWrapper.GetOrgUserGroup(), nil
}

func NewOrgUserGroupWrapper() *OrgUserGroupWrapper {
	return &OrgUserGroupWrapper{
		OrgUserGroup: &organization.UserGroup{},
	}
}

func (OrgUserGroupWrapper *OrgUserGroupWrapper) GetOrgUserGroup() *organization.UserGroup {
	return OrgUserGroupWrapper.OrgUserGroup
}

func (OrgUserGroupWrapper *OrgUserGroupWrapper) SetOrgUserGroup(OrgUserGroup *organization.UserGroup) {
	OrgUserGroupWrapper.OrgUserGroup = OrgUserGroup
}
