package acmeclient

import (
	"testing"

	opnacme "github.com/browningluke/opnsense-go/pkg/acmeclient"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertAccountSchemaToStruct(t *testing.T) {
	model := &accountResourceModel{
		Enabled:     types.BoolValue(true),
		Name:        types.StringValue("staging"),
		Description: types.StringValue("staging account"),
		Email:       types.StringValue("admin@example.com"),
		Ca:          types.StringValue("letsencrypt_test"),
		CustomCa:    types.StringNull(),
		EabKid:      types.StringNull(),
		EabHmac:     types.StringNull(),
		AccountKey:  types.StringNull(),
	}

	account, err := convertAccountSchemaToStruct(model)
	require.NoError(t, err)
	require.Equal(t, "1", account.Enabled)
	require.Equal(t, "staging", account.Name)
	require.Equal(t, "admin@example.com", account.Email)
	require.Equal(t, api.SelectedMap("letsencrypt_test"), account.Ca)
}

func TestConvertAccountStructToSchema(t *testing.T) {
	account := &opnacme.Account{
		Enabled:          "1",
		Name:             "staging",
		Description:      "staging account",
		Email:            "admin@example.com",
		Ca:               api.SelectedMap("letsencrypt_test"),
		Key:              "-----BEGIN PRIVATE KEY-----\nTEST\n-----END PRIVATE KEY-----",
		StatusCode:       "200",
		StatusLastUpdate: "1234567890",
	}

	model, err := convertAccountStructToSchema(account)
	require.NoError(t, err)
	require.True(t, model.Enabled.ValueBool())
	require.Equal(t, "staging", model.Name.ValueString())
	require.Equal(t, "letsencrypt_test", model.Ca.ValueString())
	require.Equal(t, "-----BEGIN PRIVATE KEY-----\nTEST\n-----END PRIVATE KEY-----", model.AccountKey.ValueString())
	require.Equal(t, "200", model.StatusCode.ValueString())
}
