package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postmarkSDK "github.com/keighl/postmark"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"apitokens": &schema.Schema{
				Type:      schema.TypeList,
				Computed:  true,
				Sensitive: true,
				MinItems:  1,
				Optional:  true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/servers", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Account-Token", c.AccountToken)

	server := postmarkSDK.Server{}

	server.Name = d.Get("name").(string)
	server.Color = d.Get("color").(string)
	body, err := json.Marshal(server)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	if err != nil {
		return diag.FromErr(err)
	}
	err = json.NewDecoder(res.Body).Decode(&server)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(server.ID, 10))
	d.Set("apitokens", flattenStringList(server.ApiTokens))

	return diags
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}