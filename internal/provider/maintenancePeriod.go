package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type MaintenancePeriod struct {
	Id              int    `json:"id"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DayOfWeek       string `json:"dayOfWeek"`
	Checks          []int  `json:"checks"`
	CheckGroups     []int  `json:"checkGroups"`
	DashboardGroups []int  `json:"dashboardGroups"`
}

func mapToMaintenancePeriod(maintenancePeriodModel MaintenancePeriodModel) MaintenancePeriod {
	maintenancePeriod := MaintenancePeriod{
		Id:              int(maintenancePeriodModel.Id.ValueInt32()),
		Description:     maintenancePeriodModel.Description.ValueString(),
		Enabled:         maintenancePeriodModel.Enabled.ValueBool(),
		StartTime:       maintenancePeriodModel.StartTime.ValueString(),
		EndTime:         maintenancePeriodModel.EndTime.ValueString(),
		DayOfWeek:       maintenancePeriodModel.DayOfWeek.ValueString(),
		Checks:          make([]int, 0),
		CheckGroups:     make([]int, 0),
		DashboardGroups: make([]int, 0),
	}

	for _, checkId := range maintenancePeriodModel.Checks {
		maintenancePeriod.Checks = append(maintenancePeriod.Checks, int(checkId.ValueInt32()))
	}

	for _, checkGroupId := range maintenancePeriodModel.CheckGroups {
		maintenancePeriod.CheckGroups = append(maintenancePeriod.CheckGroups, int(checkGroupId.ValueInt32()))
	}

	for _, dashboardGroupId := range maintenancePeriodModel.DashboardGroups {
		maintenancePeriod.DashboardGroups = append(maintenancePeriod.DashboardGroups, int(dashboardGroupId.ValueInt32()))
	}

	return maintenancePeriod
}

func mapToMaintenancePeriodModel(maintenancePeriod MaintenancePeriod) MaintenancePeriodModel {
	maintenancePeriodModel := MaintenancePeriodModel{
		Id:          types.Int32Value(int32(maintenancePeriod.Id)),
		Description: types.StringValue(maintenancePeriod.Description),
		Enabled:     types.BoolValue(maintenancePeriod.Enabled),
		StartTime:   types.StringValue(maintenancePeriod.StartTime),
		EndTime:     types.StringValue(maintenancePeriod.EndTime),
		DayOfWeek:   types.StringValue(maintenancePeriod.DayOfWeek),
	}

	for _, checkId := range maintenancePeriod.Checks {
		maintenancePeriodModel.Checks = append(maintenancePeriodModel.Checks, types.Int32Value(int32(checkId)))
	}

	for _, checkGroupId := range maintenancePeriod.CheckGroups {
		maintenancePeriodModel.CheckGroups = append(maintenancePeriodModel.CheckGroups, types.Int32Value(int32(checkGroupId)))
	}

	for _, dashboardGroupId := range maintenancePeriod.DashboardGroups {
		maintenancePeriodModel.DashboardGroups = append(maintenancePeriodModel.DashboardGroups, types.Int32Value(int32(dashboardGroupId)))
	}

	return maintenancePeriodModel
}
