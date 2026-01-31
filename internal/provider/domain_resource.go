package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j6s/sweego-provider/internal/sweego"
)

var _ resource.Resource = &SweegoDomainResource{}
var _ resource.ResourceWithImportState = &SweegoDomainResource{}

func NewSweegoDomainResource() resource.Resource {
	return &SweegoDomainResource{}
}

// SweegoDomainResource defines the resource implementation.
type SweegoDomainResource struct {
	api *sweego.SweegoApi
}

// SweegoDomainResourceModel describes the resource data model.
type SweegoDomainResourceModel struct {
	Uuid                 types.String `tfsdk:"uuid"`
	IsVerified           types.Bool   `tfsdk:"is_verified"`
	TrackingOpenEnabled  types.Bool   `tfsdk:"tracking_open_enabled"`
	TrackingClickEnabled types.Bool   `tfsdk:"tracking_click_enabled"`
	Domain               types.String `tfsdk:"domain"`
	DomainRecord         types.Object `tfsdk:"domain_record"`
	DkimRecord           types.Object `tfsdk:"dkim_record"`
	DmarcRecord          types.Object `tfsdk:"dmarc_record"`
	InboundRecordList    types.List   `tfsdk:"inbound_record_list"`
	TrackingRecord       types.Object `tfsdk:"tracking_record"`
}

func (r *SweegoDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

var dnsRecordAttributes = map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
	"data": types.StringType,
}

func (r *SweegoDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Sweego domain",

		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "Domain name of the domain that should be managed (e.g. my-domain.eu)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "UUID of the domain in sweego's system.",
				Computed:    true,
			},
			"tracking_click_enabled": schema.BoolAttribute{
				Description: "Whether or not click tracking is enabled",
				Computed:    true,
			},
			"tracking_open_enabled": schema.BoolAttribute{
				Description: "Whether or not open tracking is enabled",
				Computed:    true,
			},
			"is_verified": schema.BoolAttribute{
				Description: "Whether or not the domain is verified",
				Computed:    true,
			},
			"domain_record": schema.ObjectAttribute{
				Description:    "CNAME DNS Record that needs to be set in order to verify the domain",
				Computed:       true,
				AttributeTypes: dnsRecordAttributes,
			},
			"dkim_record": schema.ObjectAttribute{
				Description:    "DKIM DNS Record that needs to be set in order to send E-Mails",
				Computed:       true,
				AttributeTypes: dnsRecordAttributes,
			},
			"dmarc_record": schema.ObjectAttribute{
				Description:    "DMARC DNS Record that needs to be set in order to send E-Mails",
				Computed:       true,
				AttributeTypes: dnsRecordAttributes,
			},
			"inbound_record_list": schema.ListAttribute{
				Description: "List of DNS Records that need to be set, if sweego should accept E-Mails",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: dnsRecordAttributes,
				},
			},
			"tracking_record": schema.ObjectAttribute{
				Description:    "CNAME DNS Record that needs to be set in order to use tracking",
				Computed:       true,
				AttributeTypes: dnsRecordAttributes,
			},
		},
	}
}

func (r *SweegoDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	api, ok := req.ProviderData.(*sweego.SweegoApi)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sweego.SweegoApi, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.api = api
}

func (r *SweegoDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SweegoDomainResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.api.WithLogger(NewLoggerAdapter(ctx)).CreateDomain(data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating domain", err.Error())
		return
	}

	data = r.fillStateFromResponse(domain, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SweegoDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SweegoDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.api.WithLogger(NewLoggerAdapter(ctx)).GetDomain(data.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain", fmt.Sprintf("Error reading domain: %s", err.Error()))
		return
	}

	data = r.fillStateFromResponse(domain, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SweegoDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Error updating domain", "Updating existing domains is not supported")
}

func (r *SweegoDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SweegoDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.api.WithLogger(NewLoggerAdapter(ctx)).DeleteDomain(data.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting domain", fmt.Sprintf("Error deleting domain: %s", err.Error()))
	}
}

func (r *SweegoDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	domain, err := r.api.WithLogger(NewLoggerAdapter(ctx)).GetDomain(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain", fmt.Sprintf("Error reading domain: %s", err.Error()))
	}

	data := r.fillStateFromResponse(domain, SweegoDomainResourceModel{})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SweegoDomainResource) fillStateFromResponse(response sweego.SweegoDomainDetails, state SweegoDomainResourceModel) SweegoDomainResourceModel {
	if response.Uuid != "" {
		state.Uuid = types.StringValue(response.Uuid)
	}
	state.Domain = types.StringValue(response.Domain)
	state.IsVerified = types.BoolValue(response.IsVerified)
	state.TrackingOpenEnabled = types.BoolValue(response.TrackingOpenEnabled)
	state.TrackingClickEnabled = types.BoolValue(response.TrackingClickEnabled)
	state.DomainRecord = recordToObject(response.DomainRecord)
	state.DkimRecord = recordToObject(response.DkimRecord)
	state.DmarcRecord = recordToObject(response.DmarcRecord)
	state.TrackingRecord = recordToObject(response.TrackingRecord)

	recordList := make([]attr.Value, len(response.InboundRecordList))
	for i, record := range response.InboundRecordList {
		recordList[i] = recordToObject(record)
	}
	state.InboundRecordList = types.ListValueMust(types.ObjectType{
		AttrTypes: dnsRecordAttributes,
	}, recordList)

	return state
}

func recordToObject(record sweego.SweegoDomainRecord) types.Object {
	return types.ObjectValueMust(dnsRecordAttributes, map[string]attr.Value{
		"type": types.StringValue(record.Type),
		"name": types.StringValue(record.Name),
		"data": types.StringValue(record.Data),
	})
}
