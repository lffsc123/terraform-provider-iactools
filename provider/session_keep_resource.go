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

// 会话保持
var _ resource.Resource = &SessionKeepResource{}
var _ resource.ResourceWithImportState = &SessionKeepResource{}

func NewSessionKeepResource() resource.Resource {
	return &SessionKeepResource{}
}

// ExampleResource defines the resource implementation.
type SessionKeepResource struct {
	client *Client
}

// ExampleResourceModel describes the resource data model.
type SessionKeepResourceModel struct {
	Rsinfo SessionKeepParameter `tfsdk:"sessionkeep"`
}

type SessionKeepParameter struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	//ActiveTime        types.String `tfsdk:"activeTime"`
	//OverrideLimit     types.String `tfsdk:"overrideLimit"`
	//MatchAcrossVs     types.String `tfsdk:"matchAcrossVs"`
	//MatchFailAction   types.String `tfsdk:"matchFailAction"`
	//PrefixType        types.String `tfsdk:"prefixType"`
	//PrefixLen         types.String `tfsdk:"prefixLen"`
	//Sessioncookie     types.String `tfsdk:"sessioncookie"`
	//Cookiename        types.String `tfsdk:"cookiename"`
	//Cookiemode        types.String `tfsdk:"cookiemode"`
	//CookieEncryMode   types.String `tfsdk:"cookieEncryMode"`
	//CookieEncryPasswd types.String `tfsdk:"cookieEncryPasswd"`
	//RadiusAttribute   types.String `tfsdk:"radiusAttribute"`
	//RelatePolicy      types.String `tfsdk:"relatePolicy"`
	//RadiusStopProc    types.String `tfsdk:"radiusStopProc"`
	//HashCondition     types.String `tfsdk:"hashCondition"`
	//SipHeadType       types.String `tfsdk:"sipHeadType"`
	//SipHeadName       types.String `tfsdk:"sipHeadName"`
	//RequestContent    types.String `tfsdk:"requestContent"`
	//ReplyContent      types.String `tfsdk:"replyContent"`
	//VirtualSystem     types.String `tfsdk:"virtualSystem"`
	//HttponlyAttribute types.String `tfsdk:"httponlyAttribute"`
	//SecureAttribute   types.String `tfsdk:"secureAttribute"`
}

func (r *SessionKeepResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_SessionKeep"
}

func (r *SessionKeepResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sessionkeep": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
					//"activeTime": schema.StringAttribute{
					//	Required: true,
					//},
					//"overrideLimit": schema.StringAttribute{
					//	Required: true,
					//},
					//"matchAcrossVs": schema.StringAttribute{
					//	Optional: true,
					//},
					//"matchFailAction": schema.StringAttribute{
					//	Optional: true,
					//},
					//"prefixType": schema.StringAttribute{
					//	Optional: true,
					//},
					//"prefixLen": schema.StringAttribute{
					//	Optional: true,
					//},
					//"sessioncookie": schema.StringAttribute{
					//	Optional: true,
					//},
					//"cookiename": schema.StringAttribute{
					//	Optional: true,
					//},
					//"cookiemode": schema.StringAttribute{
					//	Optional: true,
					//},
					//"cookieEncryMode": schema.StringAttribute{
					//	Optional: true,
					//},
					//"cookieEncryPasswd": schema.StringAttribute{
					//	Optional: true,
					//},
					//"radiusAttribute": schema.StringAttribute{
					//	Optional: true,
					//},
					//"relatePolicy": schema.StringAttribute{
					//	Optional: true,
					//},
					//"radiusStopProc": schema.StringAttribute{
					//	Optional: true,
					//},
					//"hashCondition": schema.StringAttribute{
					//	Optional: true,
					//},
					//"sipHeadType": schema.StringAttribute{
					//	Optional: true,
					//},
					//"sipHeadName": schema.StringAttribute{
					//	Optional: true,
					//},
					//"requestContent": schema.StringAttribute{
					//	Optional: true,
					//},
					//"replyContent": schema.StringAttribute{
					//	Optional: true,
					//},
					//"virtualSystem": schema.StringAttribute{
					//	Optional: true,
					//},
					//"httponlyAttribute": schema.StringAttribute{
					//	Optional: true,
					//},
					//"secureAttribute": schema.StringAttribute{
					//	Optional: true,
					//},
				},
			},
		},
	}
}

func (r *SessionKeepResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SessionKeepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SessionKeepResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource")
	sendToweb_SessionKeepRequest(ctx, "POST", r.client, data.Rsinfo)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionKeepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SessionKeepResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_SessionKeepRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionKeepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SessionKeepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_SessionKeepRequest(ctx, "PUT", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionKeepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SessionKeepResourceModel
	tflog.Info(ctx, " Delete Start")

	//sendToweb_SessionKeepRequest(ctx, "DELETE", r.client, data.Rsinfo)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *SessionKeepResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_SessionKeepRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo SessionKeepParameter) {
	sendData := SessionKeepRequestModel{
		Name: Rsinfo.Name.ValueString(),
		Type: Rsinfo.Type.ValueString(),
		//ActiveTime:        Rsinfo.ActiveTime.ValueString(),
		//OverrideLimit:     Rsinfo.OverrideLimit.ValueString(),
		//MatchAcrossVs:     Rsinfo.MatchAcrossVs.ValueString(),
		//MatchFailAction:   Rsinfo.MatchFailAction.ValueString(),
		//PrefixType:        Rsinfo.PrefixType.ValueString(),
		//PrefixLen:         Rsinfo.PrefixLen.ValueString(),
		//Sessioncookie:     Rsinfo.Sessioncookie.ValueString(),
		//Cookiename:        Rsinfo.Cookiename.ValueString(),
		//Cookiemode:        Rsinfo.Cookiemode.ValueString(),
		//CookieEncryMode:   Rsinfo.CookieEncryMode.ValueString(),
		//CookieEncryPasswd: Rsinfo.CookieEncryPasswd.ValueString(),
		//RadiusAttribute:   Rsinfo.RadiusAttribute.ValueString(),
		//RelatePolicy:      Rsinfo.RelatePolicy.ValueString(),
		//RadiusStopProc:    Rsinfo.RadiusStopProc.ValueString(),
		//HashCondition:     Rsinfo.HashCondition.ValueString(),
		//SipHeadType:       Rsinfo.SipHeadType.ValueString(),
		//SipHeadName:       Rsinfo.SipHeadName.ValueString(),
		//RequestContent:    Rsinfo.RequestContent.ValueString(),
		//ReplyContent:      Rsinfo.ReplyContent.ValueString(),
		//VirtualSystem:     Rsinfo.VirtualSystem.ValueString(),
		//HttponlyAttribute: Rsinfo.HttponlyAttribute.ValueString(),
		//SecureAttribute:   Rsinfo.SecureAttribute.ValueString(),
	}

	requstData := SessionKeepRequest{
		sessionkeep: sendData,
	}

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/slb/session_keep/adx_slb_session_keep/sessionkeep"

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
