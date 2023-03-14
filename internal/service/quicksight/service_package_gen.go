// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package quicksight

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceDataSource,
			TypeName: "aws_quicksight_data_source",
		},
		{
			Factory:  ResourceGroup,
			TypeName: "aws_quicksight_group",
		},
		{
			Factory:  ResourceGroupMembership,
			TypeName: "aws_quicksight_group_membership",
		},
		{
			Factory:  ResourceUser,
			TypeName: "aws_quicksight_user",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.QuickSight
}

var ServicePackage = &servicePackage{}
