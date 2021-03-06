package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSnapshotRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			// Computed
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_option": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"disk_encryption_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_url": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"key_encryption_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_url": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceArmSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).snapshotsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Snapshot %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error loading Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if props := resp.SnapshotProperties; props != nil {
		d.Set("os_type", string(props.OsType))
		d.Set("time_created", props.TimeCreated.String())

		if props.DiskSizeGB != nil {
			d.Set("disk_size_gb", int(*props.DiskSizeGB))
		}

		if props.EncryptionSettings != nil {
			d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettings))
		}
	}

	if data := resp.CreationData; data != nil {
		d.Set("creation_option", string(data.CreateOption))
		if data.SourceURI != nil {
			d.Set("source_uri", *data.SourceURI)
		}
		if data.SourceResourceID != nil {
			d.Set("source_resource_id", *data.SourceResourceID)
		}
		if data.StorageAccountID != nil {
			d.Set("storage_account_id", *data.StorageAccountID)
		}
	}

	return nil
}
