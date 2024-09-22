package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type DashboardGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func mapToDashboardGroup(dashboardGroupModel DashboardGroupModel) DashboardGroup {
	return DashboardGroup{
		Id:          int(dashboardGroupModel.Id.ValueInt32()),
		Name:        dashboardGroupModel.Name.ValueString(),
		Description: dashboardGroupModel.Description.ValueString(),
	}
}

func mapToDashboardGroupModel(dashboardGroup DashboardGroup) DashboardGroupModel {
	return DashboardGroupModel{
		Id:          types.Int32Value(int32(dashboardGroup.Id)),
		Name:        types.StringValue(dashboardGroup.Name),
		Description: types.StringValue(dashboardGroup.Description),
	}
}
