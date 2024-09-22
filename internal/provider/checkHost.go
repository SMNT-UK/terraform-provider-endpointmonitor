package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type CheckHost struct {
	Id                  int     `json:"id"`
	Hostname            string  `json:"hostname"`
	Description         string  `json:"description"`
	Type                *string `json:"type"`
	Enabled             bool    `json:"enabled"`
	MaxWebJourneyChecks int     `json:"maxWebJourneyChecks"`
	SendCheckFiles      bool    `json:"sendCheckFiles"`
}

func mapToCheckHost(checkHostModel CheckHostModel) CheckHost {
	return CheckHost{
		Id:                  int(checkHostModel.Id.ValueInt32()),
		Hostname:            checkHostModel.Hostname.ValueString(),
		Description:         checkHostModel.Description.ValueString(),
		Type:                checkHostModel.Type.ValueStringPointer(),
		Enabled:             checkHostModel.Enabled.ValueBool(),
		MaxWebJourneyChecks: int(checkHostModel.MaxWebJourneyChecks.ValueInt32()),
		SendCheckFiles:      checkHostModel.SendCheckFiles.ValueBool(),
	}
}

func mapToCheckHostModel(checkHost CheckHost) CheckHostModel {
	return CheckHostModel{
		Id:                  types.Int32Value(int32(checkHost.Id)),
		Hostname:            types.StringValue(checkHost.Hostname),
		Description:         types.StringValue(checkHost.Description),
		Type:                types.StringPointerValue(checkHost.Type),
		Enabled:             types.BoolValue(checkHost.Enabled),
		MaxWebJourneyChecks: types.Int32Value(int32(checkHost.MaxWebJourneyChecks)),
		SendCheckFiles:      types.BoolValue(checkHost.SendCheckFiles),
	}
}
