package aom

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
	"io/ioutil"
	"time"
)

func ResourceApmAkSk() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceApmAkSkCreate,
		ReadContext:   ResourceApmAkSkRead,
		UpdateContext: ResourceApmAkSkUpdate,
		DeleteContext: ResourceApmAkSkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"descp": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ak": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sk": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceApmAkSkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diagErr := httpclient_go.NewHttpClientGo(cfg)
	if diagErr != nil {
		return diagErr
	}

	opts := entity.CreateAkSkParam{
		Descp: d.Get("descp").(string),
	}

	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys").WithBody(opts)

	response, err := client.Do()
	if err != nil {
		diag.Errorf("error creating aksk fields %s: client do error: %s", opts, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s to %v : %s", string(body), opts, err)
	}
	if response.StatusCode == 200 {
		rlt := &entity.AkSkResultVO{}

		err = json.Unmarshal(body, rlt)

		d.SetId("success")
		d.Set("ak", rlt.Ak)
		d.Set("sk", rlt.Sk)
		return nil
	}
	return diag.Errorf("error creating aksk fields %s: %s", opts, string(body))
}

func ResourceApmAkSkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diagErr := httpclient_go.NewHttpClientGo(cfg)
	if diagErr != nil {
		return diagErr
	}

	client.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys")

	response, err := client.Do()
	logp.Printf("[TEST]", err, response)

	body, diags := client.CheckDeletedDiag(d, err, response, "error to query akSks")
	if diags != nil {
		return diags
	}

	rlt := &entity.GetAkSkListVO{}
	err = json.Unmarshal(body, &rlt)
	if err != nil {
		return diag.Errorf("error covert data %s: %s", string(body), err)
	}
	for _, item := range rlt.AccessAkSkModels {
		if item.Ak == d.Get("ak").(string) {
			return nil
		}
	}

	resourceID := d.Id()
	d.SetId("")
	d.Set("ak", "")
	d.Set("sk", "")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource not found",
			Detail:   fmt.Sprintf("the resource %s is goneand will be removed in Teraform state.", resourceID),
		},
	}
}

func ResourceApmAkSkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceApmAkSkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diagErr := httpclient_go.NewHttpClientGo(cfg)
	if diagErr != nil {
		return diagErr
	}

	client.WithMethod(httpclient_go.MethodDelete).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys/"+d.Get("ak").(string))

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete aksk %s: %s", d.Get("application_id"), err)
	}

	if response.StatusCode == 200 {
		return nil
	}
	mErr := &multierror.Error{}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	rlt := &entity.ErrorResp{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error create A fields: %w", err)
	}

	return nil
}
