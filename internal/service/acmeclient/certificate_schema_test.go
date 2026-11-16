package acmeclient

import (
	"testing"

	opnacme "github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertCertificateSchemaToStruct(t *testing.T) {
	model := &certificateResourceModel{
		Enabled:          types.BoolValue(true),
		Name:             types.StringValue("example.com"),
		Description:      types.StringValue("web cert"),
		AltNames:         tools.StringSliceToSet([]string{"www.example.com", "api.example.com"}),
		Account:          types.StringValue("account-uuid"),
		ValidationMethod: types.StringValue("validation-uuid"),
		KeyLength:        types.StringValue("key_2048"),
		Ocsp:             types.BoolValue(false),
		AutoRenewal:      types.BoolValue(true),
		RenewInterval:    types.StringValue("60"),
		Aliasmode:        types.StringValue("none"),
		Domainalias:      types.StringNull(),
		Challengealias:   types.StringNull(),
		RestartActions:   tools.StringSliceToSet([]string{"action-uuid"}),
	}

	cert, err := convertCertificateSchemaToStruct(model)
	require.NoError(t, err)
	require.Equal(t, "1", cert.Enabled)
	require.Equal(t, "example.com", cert.Name)
	require.Equal(t, api.SelectedMap("account-uuid"), cert.Account)
	require.Equal(t, api.SelectedMap("validation-uuid"), cert.ValidationMethod)
	require.Equal(t, api.SelectedMap("key_2048"), cert.KeyLength)
	require.Equal(t, "1", cert.AutoRenewal)
	require.Equal(t, "60", cert.RenewInterval)
	require.ElementsMatch(t, []string{"www.example.com", "api.example.com"}, []string(cert.AltNames))
	require.ElementsMatch(t, []string{"action-uuid"}, []string(cert.RestartActions))
}

func TestConvertCertificateStructToSchema(t *testing.T) {
	cert := &opnacme.Certificate{
		Enabled:          "1",
		Name:             "example.com",
		AltNames:         api.SelectedMapList{"www.example.com"},
		Account:          api.SelectedMap("account-uuid"),
		ValidationMethod: api.SelectedMap("validation-uuid"),
		KeyLength:        api.SelectedMap("key_2048"),
		Ocsp:             "0",
		AutoRenewal:      "1",
		RenewInterval:    "60",
		Aliasmode:        api.SelectedMap("none"),
		RestartActions:   api.SelectedMapList{"action-uuid"},
		CertRefId:        "abc123",
		StatusCode:       "200",
	}

	model, err := convertCertificateStructToSchema(cert)
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.Equal(t, "example.com", model.Name.ValueString())
	require.Equal(t, "account-uuid", model.Account.ValueString())
	require.Equal(t, "abc123", model.CertRefId.ValueString())
	require.ElementsMatch(t, []string{"www.example.com"}, tools.SetToStringSlice(model.AltNames))
	require.ElementsMatch(t, []string{"action-uuid"}, tools.SetToStringSlice(model.RestartActions))
}
