package aom

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io/ioutil"
	"time"
)

func ResourceAomApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAomApplicationCreate,
		ReadContext:   ResourceAomApplicationRead,
		UpdateContext: ResourceAomApplicationUpdate,
		DeleteContext: ResourceAomApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aom_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"register_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceAomApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity.BizAppParam{
		Description:  d.Get("description").(string),
		DisplayName:  d.Get("display_name").(string),
		EpsId:        d.Get("enterprise_project_id").(string),
		Name:         d.Get("name").(string),
		RegisterType: d.Get("register_type").(string),
	}

	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/applications").WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error create Application fields %s: %s", opts, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 200 {
		rlt := &entity.CreateModelVo{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		d.SetId(rlt.Id)
		return ResourceAomApplicationRead(ctx, d, meta)
	}
	return diag.Errorf("error create Application %v. error: %s", opts, string(body))
}

func ResourceAomApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}
	client.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/applications/"+d.Id())
	response, err := client.Do()

	body, diags := client.CheckDeletedDiag(d, err, response, "error retrieving Application")
	if body == nil {
		return diags
	}

	rlt := &entity.BizAppVo{}
	err = json.Unmarshal(body, rlt)
	if err != nil {
		return diag.Errorf("error retrieving Application %s", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("aom_id", rlt.AomId),
		d.Set("app_id", rlt.AppId),
		d.Set("create_time", rlt.CreateTime),
		d.Set("creator", rlt.Creator),
		d.Set("description", rlt.Description),
		d.Set("display_name", rlt.DisplayName),
		d.Set("enterprise_project_id", rlt.EpsId),
		d.Set("modified_time", rlt.ModifiedTime),
		d.Set("modifier", rlt.Modifier),
		d.Set("name", rlt.Name),
		d.Set("register_type", rlt.RegisterType),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Application fields: %s", err)
	}

	return nil
}

func ResourceAomApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}
	opts := entity.BizAppParam{
		Description:  d.Get("description").(string),
		DisplayName:  d.Get("display_name").(string),
		EpsId:        d.Get("enterprise_project_id").(string),
		Name:         d.Get("name").(string),
		RegisterType: d.Get("register_type").(string),
	}
	client.WithMethod(httpclient_go.MethodPut).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/applications/"+d.Id()).WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update Application fields %s: %s", opts, err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update Application %s: %s", string(body), err)
	}
	return diag.Errorf("error update Application %s:  %s", opts, string(body))
}

func ResourceAomApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodDelete).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/applications/"+d.Id())

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete Application %s: %s", d.Id(), err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete Application %s: %s", d.Id(), err)
	}
	return diag.Errorf("error delete Application %s:  %s", d.Id(), string(body))
}
