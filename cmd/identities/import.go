package identities

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/server"
	"github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type RootCmd struct {
	Import ImportCmd `cmd:"" help:"Import identities and permissions"`
}

type ImportCmd struct{}

func (s *ImportCmd) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)
	ketoClient, err := keto.NewClient()
	if err != nil {
		return err
	}

	customersFile, err := os.ReadFile("ory/kratos/imports/customers.json")
	if err != nil {
		return fmt.Errorf("failed to read customers.json file: %w", err)
	}
	developersFile, err := os.ReadFile("ory/kratos/imports/developers.json")
	if err != nil {
		return fmt.Errorf("failed to read developers.json file: %w", err)
	}

	var createCustomerIdentities []client.CreateIdentityBody
	if err := json.Unmarshal(customersFile, &createCustomerIdentities); err != nil {
		return fmt.Errorf("failed to unmarshal customers.json file: %w", err)
	}
	var createDeveloperIdentities []client.CreateIdentityBody
	if err := json.Unmarshal(developersFile, &createDeveloperIdentities); err != nil {
		return fmt.Errorf("failed to unmarshal developers.json file: %w", err)
	}

	var customerIdentitiesPatch []client.IdentityPatch
	for _, createBody := range createCustomerIdentities {
		customerIdentitiesPatch = append(customerIdentitiesPatch, client.IdentityPatch{Create: &createBody})
	}
	var developerIdentitiesPatch []client.IdentityPatch
	for _, createBody := range createDeveloperIdentities {
		developerIdentitiesPatch = append(developerIdentitiesPatch, client.IdentityPatch{Create: &createBody})
	}

	ctx := context.Background()
	reqCust := kratosClient.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	reqCust = reqCust.PatchIdentitiesBody(client.PatchIdentitiesBody{Identities: customerIdentitiesPatch})
	batchResultCust, identityCustResp, err := reqCust.Execute()
	if err != nil {
		return fmt.Errorf("failed to import customer identities: %w", err)
	}
	defer identityCustResp.Body.Close()

	reqDev := kratosClient.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	reqDev = reqDev.PatchIdentitiesBody(client.PatchIdentitiesBody{Identities: developerIdentitiesPatch})
	batchResultDev, identityDevResp, err := reqDev.Execute()
	if err != nil {
		return fmt.Errorf("failed to import developer identities: %w", err)
	}
	defer identityDevResp.Body.Close()

	var tuples []*rts.RelationTupleDelta
	for _, customerResult := range batchResultCust.Identities {
		if *customerResult.Action == "error" {
			fmt.Printf("error creating tuple for customer: %s", *customerResult.Identity)
			continue
		}
		group := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "customer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(server.AuthzIdentityID(*customerResult.Identity)),
			},
		}
		owner := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    server.AuthzIdentityID(*customerResult.Identity),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(server.AuthzIdentityID(*customerResult.Identity)),
			},
		}
		parents := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    server.AuthzIdentityID(*customerResult.Identity),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		}
		tuples = append(tuples, group, owner, parents)
	}

	for _, devResult := range batchResultDev.Identities {
		if *devResult.Action == "error" {
			fmt.Printf("error creating tuple for customer: %s", *devResult.Identity)
			continue
		}
		group := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "developer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(server.AuthzIdentityID(*devResult.Identity)),
			},
		}
		owner := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    server.AuthzIdentityID(*devResult.Identity),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(server.AuthzIdentityID(*devResult.Identity)),
			},
		}
		parents := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    server.AuthzIdentityID(*devResult.Identity),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		}
		tuples = append(tuples, group, owner, parents)
	}

	if _, err := ketoClient.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuples,
	}); err != nil {
		return fmt.Errorf("failed to create relation tuples for identities: %w", err)
	}

	return nil
}
