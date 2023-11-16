package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VirtualServiceResource{}
var _ resource.ResourceWithImportState = &VirtualServiceResource{}

func NewVirtualServiceResource() resource.Resource {
	return &VirtualServiceResource{}
}

// ExampleResource defines the resource implementation.
type VirtualServiceResource struct {
	client *Client
}

// ExampleResourceModel describes the resource data model.
type VirtualServiceResourceModel struct {
	// Uuid_count types.String `tfsdk:"uuid_count"`
	Rsinfo VirtualServiceParameter `tfsdk:"virtualservice"`
}

type VirtualServiceParameter struct {
	Name        types.String `tfsdk:"name"`
	State       types.String `tfsdk:"state"`
	Mode        types.String `tfsdk:"mode"`
	Ip          types.String `tfsdk:"ip"`
	Port        types.String `tfsdk:"port"`
	Protocol    types.String `tfsdk:"protocol"`
	SessionKeep types.String `tfsdk:"session_keep"`
	DefaultPool types.String `tfsdk:"default_pool"`
	TcpPolicy   types.String `tfsdk:"tcp_policy"` //引用tcp超时时间，不引用默认600s
	Snat        types.String `tfsdk:"snat"`
	SessionBkp  types.String `tfsdk:"session_bkp"` //必须配置集群模式
	Vrrp        types.String `tfsdk:"vrrp"`        //涉及普通双机热备场景，需要关联具体的vrrp组
}

func (r *VirtualServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_VirtualService"
}

func (r *VirtualServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtualservice": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"mode": schema.StringAttribute{
						Required: true,
					},
					"ip": schema.StringAttribute{
						Required: true,
					},
					"port": schema.StringAttribute{
						Required: true,
					},
					"state": schema.StringAttribute{
						Optional: true,
					},
					"protocol": schema.StringAttribute{
						Optional: true,
					},
					"session_keep": schema.StringAttribute{
						Optional: true,
					},
					"default_pool": schema.StringAttribute{
						Optional: true,
					},
					"tcp_policy": schema.StringAttribute{
						Optional: true,
					},
					"snat": schema.StringAttribute{
						Optional: true,
					},
					"session_bkp": schema.StringAttribute{
						Optional: true,
					},
					"vrrp": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *VirtualServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*Client)

	if req.ProviderData == nil {
		return
	}
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *VirtualServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VirtualServiceResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource")
	sendToweb_VirtualServiceRequest(ctx, "POST", r.client, data.Rsinfo)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VirtualServiceResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_VirtualServiceRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VirtualServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_VirtualServiceRequest(ctx, "PUT", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtualServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VirtualServiceResourceModel
	tflog.Info(ctx, " Delete Start")

	sendToweb_VirtualServiceRequest(ctx, "DELETE", r.client, data.Rsinfo)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *VirtualServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_VirtualServiceRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo VirtualServiceParameter) {
	sendData := VirtualServiceRequestModel{
		Name:        Rsinfo.Name.ValueString(),
		State:       Rsinfo.State.ValueString(),
		Mode:        Rsinfo.Mode.ValueString(),
		Ip:          Rsinfo.Ip.ValueString(),
		Port:        Rsinfo.Port.ValueString(),
		Protocol:    Rsinfo.Protocol.ValueString(),
		SessionKeep: Rsinfo.SessionKeep.ValueString(),
		DefaultPool: Rsinfo.DefaultPool.ValueString(),
		TcpPolicy:   Rsinfo.TcpPolicy.ValueString(),
		Snat:        Rsinfo.Snat.ValueString(),
		SessionBkp:  Rsinfo.SessionBkp.ValueString(),
		Vrrp:        Rsinfo.Vrrp.ValueString(),
	}

	requstData := VirtualServiceRequest{
		Virtualservice: sendData,
	}

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/slb/vs/virtual/virtualservice"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}
