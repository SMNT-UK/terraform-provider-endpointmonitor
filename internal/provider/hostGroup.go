package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type HostGroup struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
	Hosts       []CheckHost `json:"checkHosts"`
}

func mapToHostGroup(hostGroupModel HostGroupModel) HostGroup {
	hostGroup := HostGroup{
		Id:          int(hostGroupModel.Id.ValueInt32()),
		Name:        hostGroupModel.Name.ValueString(),
		Description: hostGroupModel.Description.ValueString(),
		Enabled:     hostGroupModel.Enabled.ValueBool(),
		Hosts:       make([]CheckHost, 0),
	}

	for _, hostId := range hostGroupModel.Hosts {
		hostGroup.Hosts = append(hostGroup.Hosts, CheckHost{Id: int(hostId.ValueInt32())})
	}

	return hostGroup
}

func mapToHostGroupModel(hostGroup HostGroup) HostGroupModel {
	hostGroupModel := HostGroupModel{
		Id:          types.Int32Value(int32(hostGroup.Id)),
		Name:        types.StringValue(hostGroup.Name),
		Description: types.StringValue(hostGroup.Description),
		Enabled:     types.BoolValue(hostGroup.Enabled),
		Hosts:       []types.Int32{},
	}

	for _, host := range hostGroup.Hosts {
		hostGroupModel.Hosts = append(hostGroupModel.Hosts, types.Int32Value(int32(host.Id)))
	}

	return hostGroupModel
}
