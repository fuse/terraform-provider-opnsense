package acmeclient

import (
	"testing"

	opnacme "github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertValidationSchemaToStruct(t *testing.T) {
	model := &validationResourceModel{
		Enabled:        types.BoolValue(false),
		Name:           types.StringValue("cf-dns"),
		Description:    types.StringValue("cloudflare dns-01"),
		Method:         types.StringValue("dns01"),
		DnsService:     types.StringValue("dns_cf"),
		DnsSleep:       types.StringValue("30"),
		DnsCfEmail:     types.StringNull(),
		DnsCfKey:       types.StringNull(),
		DnsCfToken:     types.StringValue("dummy-token"),
		DnsCfAccountId: types.StringNull(),
		DnsCfZoneId:    types.StringValue("zone123"),
	}

	validation, err := convertValidationSchemaToStruct(model)
	require.NoError(t, err)
	require.Equal(t, "0", validation.Enabled)
	require.Equal(t, "cf-dns", validation.Name)
	require.Equal(t, api.SelectedMap("dns01"), validation.Method)
	require.Equal(t, api.SelectedMap("opnsense"), validation.HttpService)
	require.Equal(t, api.SelectedMap("acme"), validation.TlsalpnService)
	require.Equal(t, api.SelectedMap("dns_cf"), validation.DnsService)
	require.Equal(t, "dummy-token", validation.DnsCfToken)
	require.Equal(t, "zone123", validation.DnsCfZoneId)
	require.Empty(t, validation.DnsAwsId)
}

func TestOverlayValidationManagedFields(t *testing.T) {
	existing := &opnacme.Validation{
		Enabled:        "1",
		Name:           "old",
		Method:         api.SelectedMap("dns01"),
		HttpService:    api.SelectedMap("custom-http"),
		TlsalpnService: api.SelectedMap("custom-tls"),
		DnsService:     api.SelectedMap("dns_cf"),
		DnsAwsId:       "keep-me",
		DnsCfToken:     "old-token",
	}
	planned := &opnacme.Validation{
		Enabled:        "0",
		Name:           "new",
		Method:         api.SelectedMap("dns01"),
		HttpService:    api.SelectedMap("opnsense"),
		TlsalpnService: api.SelectedMap("acme"),
		DnsService:     api.SelectedMap("dns_cf"),
		DnsCfToken:     "new-token",
	}

	overlayValidationManagedFields(existing, planned)
	require.Equal(t, "0", existing.Enabled)
	require.Equal(t, "new", existing.Name)
	require.Equal(t, "new-token", existing.DnsCfToken)
	require.Equal(t, "keep-me", existing.DnsAwsId)
	require.Equal(t, api.SelectedMap("custom-http"), existing.HttpService)
	require.Equal(t, api.SelectedMap("custom-tls"), existing.TlsalpnService)
}

func TestConvertValidationStructToSchema(t *testing.T) {
	validation := &opnacme.Validation{
		Enabled:    "1",
		Name:       "cf-dns",
		Method:     api.SelectedMap("dns01"),
		DnsService: api.SelectedMap("dns_cf"),
		DnsCfToken: "secret",
	}

	model, err := convertValidationStructToSchema(validation)
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.Equal(t, "cf-dns", model.Name.ValueString())
	require.Equal(t, "secret", model.DnsCfToken.ValueString())
	require.Equal(t, "0", model.DnsSleep.ValueString())
}
