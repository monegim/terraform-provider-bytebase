package provider

import (
	"testing"

	v1pb "buf.build/gen/go/bytebase/bytebase/protocolbuffers/go/v1"
)

// The API never returns the Vault AppRole secret, and a database read has no prior
// state to copy role_id/secret from. Flattening the instance resource must not panic.
func TestFlattenDatabaseInstanceResourceVaultAppRoleWithoutState(t *testing.T) {
	instance := &v1pb.InstanceResource{
		Name:   "instances/postgres-17",
		Title:  "postgres-17",
		Engine: v1pb.Engine_POSTGRES,
		DataSources: []*v1pb.DataSource{
			{
				Id:       "admin",
				Type:     v1pb.DataSourceType_ADMIN,
				Username: "bytebase",
				Host:     "postgres-17.example",
				Port:     "5432",
				ExternalSecret: &v1pb.DataSourceExternalSecret{
					SecretType:      v1pb.DataSourceExternalSecret_VAULT_KV_V2,
					Url:             "https://vault.example",
					EngineName:      "kv",
					SecretName:      "databases/postgres-17",
					PasswordKeyName: "password",
					AuthType:        v1pb.DataSourceExternalSecret_VAULT_APP_ROLE,
					AuthOption: &v1pb.DataSourceExternalSecret_AppRole{
						AppRole: &v1pb.DataSourceExternalSecret_AppRoleAuthOption{
							Type: v1pb.DataSourceExternalSecret_AppRoleAuthOption_PLAIN,
						},
					},
				},
			},
		},
	}

	got := flattenDatabaseInstanceResource(instance)
	if len(got) != 1 {
		t.Fatalf("flattenDatabaseInstanceResource returned %d elements, want 1", len(got))
	}
}

func TestConvertToV1DataSourceVaultAppRoleMissingIDs(t *testing.T) {
	ds, err := convertToV1DataSource(map[string]interface{}{
		"id":   "admin",
		"type": v1pb.DataSourceType_ADMIN.String(),
		"external_secret": []interface{}{map[string]interface{}{
			"vault": []interface{}{map[string]interface{}{
				"url":               "https://vault.example",
				"engine_name":       "kv",
				"secret_name":       "databases/postgres-17",
				"password_key_name": "password",
				"app_role": []interface{}{map[string]interface{}{
					"secret_type": v1pb.DataSourceExternalSecret_AppRoleAuthOption_PLAIN.String(),
				}},
			}},
		}},
	})
	if err != nil {
		t.Fatalf("convertToV1DataSource returned error: %v", err)
	}
	if ds.GetExternalSecret().GetAppRole().GetRoleId() != "" || ds.GetExternalSecret().GetAppRole().GetSecretId() != "" {
		t.Fatalf("expected empty role_id and secret_id")
	}
}
