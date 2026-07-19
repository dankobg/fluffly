package identities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aarondl/opt/omit"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	"github.com/google/uuid"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type ImportIdentitiesCmd struct{}

func (ic *ImportIdentitiesCmd) Run() error {
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

	pool, err := postgres.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pool.Close()

	pg := postgres.New(pool)

	customersFile, err := os.ReadFile("ory/kratos/imports/customers.json")
	if err != nil {
		return fmt.Errorf("failed to read customers.json file: %w", err)
	}

	developersFile, err := os.ReadFile("ory/kratos/imports/developers.json")
	if err != nil {
		return fmt.Errorf("failed to read developers.json file: %w", err)
	}

	var createCustomerIdentities []orykratos.CreateIdentityBody
	if err := json.Unmarshal(customersFile, &createCustomerIdentities); err != nil {
		return fmt.Errorf("failed to unmarshal customers.json file: %w", err)
	}

	var createDeveloperIdentities []orykratos.CreateIdentityBody
	if err := json.Unmarshal(developersFile, &createDeveloperIdentities); err != nil {
		return fmt.Errorf("failed to unmarshal developers.json file: %w", err)
	}

	var customerIdentitiesPatch []orykratos.IdentityPatch
	for _, createBody := range createCustomerIdentities {
		customerIdentitiesPatch = append(customerIdentitiesPatch, orykratos.IdentityPatch{Create: &createBody})
	}

	var developerIdentitiesPatch []orykratos.IdentityPatch
	for _, createBody := range createDeveloperIdentities {
		developerIdentitiesPatch = append(developerIdentitiesPatch, orykratos.IdentityPatch{Create: &createBody})
	}

	ctx := context.Background()
	reqCust := kratosClient.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	reqCust = reqCust.PatchIdentitiesBody(orykratos.PatchIdentitiesBody{Identities: customerIdentitiesPatch})

	batchResultCust, identityCustResp, err := reqCust.Execute()
	if err != nil {
		return fmt.Errorf("failed to import customer identities: %w", err)
	}

	defer func() { _ = identityCustResp.Body.Close() }()

	reqDev := kratosClient.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	reqDev = reqDev.PatchIdentitiesBody(orykratos.PatchIdentitiesBody{Identities: developerIdentitiesPatch})

	batchResultDev, identityDevResp, err := reqDev.Execute()
	if err != nil {
		return fmt.Errorf("failed to import developer identities: %w", err)
	}

	defer func() { _ = identityDevResp.Body.Close() }()

	var userIDs []uuid.UUID
	for _, customerResult := range batchResultCust.GetIdentities() {
		userIDs = append(userIDs, uuid.MustParse(customerResult.GetIdentity()))
	}

	for _, devResult := range batchResultDev.GetIdentities() {
		userIDs = append(userIDs, uuid.MustParse(devResult.GetIdentity()))
	}

	for _, uid := range userIDs {
		if _, err := pg.User().CreateUser(ctx, models.UserSetter{ID: omit.From(uid)}); err != nil {
			log.Fatalln("failed to insert users:", err)
		}
	}

	var tuples []*rts.RelationTupleDelta

	for _, customerResult := range batchResultCust.GetIdentities() {
		if customerResult.GetAction() == "error" {
			fmt.Printf("error creating tuple for customer: %s\n", customerResult.GetIdentity())
			continue
		}

		group := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "customer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(customerResult.GetIdentity())),
			},
		}
		owner := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(customerResult.GetIdentity()),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(customerResult.GetIdentity())),
			},
		}
		analytic := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Analytic",
				Object:    shared.AuthzAnalyticsID(customerResult.GetIdentity()),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(customerResult.GetIdentity())),
			},
		}
		parents := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(customerResult.GetIdentity()),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		}
		tuples = append(tuples, group, owner, analytic, parents)
	}

	for _, devResult := range batchResultDev.GetIdentities() {
		if devResult.GetAction() == "error" {
			fmt.Printf("error creating tuple for customer: %s\n", devResult.GetIdentity())
			continue
		}

		group := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "developer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(devResult.GetIdentity())),
			},
		}
		owner := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(devResult.GetIdentity()),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(devResult.GetIdentity())),
			},
		}
		analytic := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Analytic",
				Object:    shared.AuthzAnalyticsID(devResult.GetIdentity()),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(devResult.GetIdentity())),
			},
		}
		parents := &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_INSERT,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(devResult.GetIdentity()),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		}
		tuples = append(tuples, group, owner, analytic, parents)
	}

	if _, err := ketoClient.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuples,
	}); err != nil {
		return fmt.Errorf("failed to create relation tuples for identities: %w", err)
	}

	return nil
}
