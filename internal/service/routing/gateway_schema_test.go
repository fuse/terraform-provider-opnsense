package routing

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/routing"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvertGatewayRoundTrip(t *testing.T) {
	model := &gatewayResourceModel{
		Enabled:        types.BoolValue(true),
		Name:           types.StringValue("WAN_GW"),
		Description:    types.StringValue("Interface WAN Gateway"),
		Interface:      types.StringValue("wan"),
		IpProtocol:     types.StringValue("inet"),
		Gateway:        types.StringValue("192.0.2.1"),
		DefaultGateway: types.BoolValue(true),
		FarGateway:     types.BoolValue(true),
		MonitorDisable: types.BoolValue(true),
		MonitorNoRoute: types.BoolValue(false),
		Monitor:        types.StringNull(),
		ForceDown:      types.BoolValue(false),
		Priority:       types.StringValue("255"),
		Weight:         types.StringValue("1"),
	}

	apiStruct, err := convertGatewaySchemaToStruct(model)
	if err != nil {
		t.Fatalf("schema to struct: %v", err)
	}
	if apiStruct.Disabled != "0" {
		t.Fatalf("disabled: got %q", apiStruct.Disabled)
	}
	if apiStruct.DefaultGateway != "1" || apiStruct.FarGateway != "1" {
		t.Fatalf("flags: default=%q far=%q", apiStruct.DefaultGateway, apiStruct.FarGateway)
	}
	if apiStruct.Interface.String() != "wan" {
		t.Fatalf("interface: got %q", apiStruct.Interface.String())
	}

	back, err := convertGatewayStructToSchema(apiStruct)
	if err != nil {
		t.Fatalf("struct to schema: %v", err)
	}
	if !back.Enabled.ValueBool() || !back.DefaultGateway.ValueBool() || !back.FarGateway.ValueBool() {
		t.Fatalf("bool round-trip failed: %+v", back)
	}
	if back.Name.ValueString() != "WAN_GW" || back.Gateway.ValueString() != "192.0.2.1" {
		t.Fatalf("string round-trip failed: %+v", back)
	}
}

func TestConvertGatewayEmptyBoolStrings(t *testing.T) {
	apiStruct := &routing.Gateway{
		Disabled:       "0",
		Name:           "gw-peer",
		Interface:      api.SelectedMap("opt2"),
		IpProtocol:     api.SelectedMap("inet"),
		Gateway:        "10.0.0.10",
		DefaultGateway: "0",
		FarGateway:     "0",
		MonitorDisable: "1",
		MonitorNoRoute: "",
		ForceDown:      "",
		Priority:       "255",
		Weight:         "1",
	}

	model, err := convertGatewayStructToSchema(apiStruct)
	if err != nil {
		t.Fatalf("struct to schema: %v", err)
	}
	if model.MonitorNoRoute.ValueBool() || model.ForceDown.ValueBool() {
		t.Fatalf("empty API bool strings should read as false")
	}
}
