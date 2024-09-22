package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type CheckGroup struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	DashboardGroup DashboardGroup `json:"dashboardGroup"`
}

func mapToCheckGroup(checkGroupModel CheckGroupModel) CheckGroup {
	return CheckGroup{
		Id:          int(checkGroupModel.Id.ValueInt32()),
		Name:        checkGroupModel.Name.ValueString(),
		Description: checkGroupModel.Description.ValueString(),
		DashboardGroup: DashboardGroup{
			Id: int(checkGroupModel.DashboardGroup.ValueInt32()),
		},
	}
}

func mapToCheckGroupModel(checkGroup CheckGroup) CheckGroupModel {
	return CheckGroupModel{
		Id:             types.Int32Value(int32(checkGroup.Id)),
		Name:           types.StringValue(checkGroup.Name),
		Description:    types.StringValue(checkGroup.Description),
		DashboardGroup: types.Int32Value(int32(checkGroup.DashboardGroup.Id)),
	}
}
