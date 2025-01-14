package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prefixDataSource{}

type prefixDataSource struct {
	providerData string
}

type prefixDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

func (d *prefixDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "dx_prefix_value"
}

func (d *prefixDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a prefix value based on the provider configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Description: "The prefix value from provider configuration",
			},
		},
	}
}

func (d *prefixDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(string)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected string, got: %T", req.ProviderData),
		)
		return
	}

	d.providerData = data
}

func (d *prefixDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state prefixDataSourceModel

	state.Value = types.StringValue(d.providerData)
	state.Id = types.StringValue(d.providerData)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
