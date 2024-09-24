package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	VAULT_TOKEN = os.Getenv("VAULT_TOKEN")
)

func getVaultClient(comm Common) (*vault.Client, error) {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(comm.VaultAddr),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// authenticate with a root token (insecure)
	if err = client.SetToken(VAULT_TOKEN); err != nil {
		return nil, err
	}

	return client, nil
}

func ExistsVault(comm Common, project string, entity string) (bool, error) {
	ctx := context.Background()

	// prepare a client with the given base address
	client, err := getVaultClient(comm)
	if err != nil {
		return false, err
	}

	m, err := client.Secrets.KvV2Read(ctx, fmt.Sprintf("openai/%s/%s", project, entity),
		vault.WithMountPath("secret"))
	if err != nil {
		return false, err
	}

	_, ok := m.Data.Data[entity]
	return ok, nil
}

func InsertVault(comm Common, project string, entity string, api_key string) error {
	ctx := context.Background()

	// prepare a client with the given base address
	client, err := getVaultClient(comm)
	if err != nil {
		return err
	}

	m, err := client.Secrets.KvV2Read(ctx, fmt.Sprintf("openai/%s/%s", project, entity), vault.WithMountPath("secret"))
	// check if err is InvalidPath
	if vault.IsErrorStatus(err, 404) {
		// write a secret
		_, err = client.Secrets.KvV2Write(
			ctx,
			fmt.Sprintf("openai/%s/%s", project, entity),
			schema.KvV2WriteRequest{
				Data: map[string]any{
					entity: api_key,
				},
				Options: map[string]any{
					"cas": 0,
				},
			},
			vault.WithMountPath("secret"),
		)
		if err != nil {
			return fmt.Errorf("failed to write: %s", err)
		}
	} else {
		// patch the secret
		v := m.Data.Metadata["version"].(json.Number)
		ver, err := v.Int64()
		if err != nil {
			return fmt.Errorf("failed to get version: %s", err)
		}
		err = KvV2Patch(fmt.Sprintf("openai/%s/%s", project, entity), uint(ver), map[string]any{
			entity: api_key,
		})
		if err != nil {
			return fmt.Errorf("failed to patch %s", err)
		}
	}

	return nil
}

type PatchOption struct {
	Cas uint `json:"cas"`
}

type PatchItem struct {
	Options PatchOption    `json:"options"`
	Data    map[string]any `json:"data"`
}

func KvV2Patch(path string, cas uint, data map[string]any) error {
	b, err := json.Marshal(&PatchItem{
		Options: PatchOption{Cas: cas},
		Data:    data,
	})

	cli := &http.Client{}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("https://vault.in.dalpha.so/v1/secret/data/%s", path),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", VAULT_TOKEN)
	req.Header.Add("Content-Type", "application/merge-patch+json")
	if res, err := cli.Do(req); err != nil {
		return fmt.Errorf("failed to patch secret: %w", err)
	} else if res.StatusCode != 200 {
		return fmt.Errorf("failed to patch secret with code: %s", res.Status)
	}

	return nil
}
