package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type ProxyHost struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
}

func mapToProxyHost(proxyHostModel ProxyHostModel) ProxyHost {
	return ProxyHost{
		Id:          int(proxyHostModel.Id.ValueInt32()),
		Name:        proxyHostModel.Name.ValueString(),
		Description: proxyHostModel.Description.ValueString(),
		Hostname:    proxyHostModel.Hostname.ValueString(),
		Port:        int(proxyHostModel.Port.ValueInt32()),
	}
}

func mapToProxyHostModel(proxyHost ProxyHost) ProxyHostModel {
	return ProxyHostModel{
		Id:          types.Int32Value(int32(proxyHost.Id)),
		Name:        types.StringValue(proxyHost.Name),
		Description: types.StringValue(proxyHost.Description),
		Hostname:    types.StringValue(proxyHost.Hostname),
		Port:        types.Int32Value(int32(proxyHost.Port)),
	}
}
