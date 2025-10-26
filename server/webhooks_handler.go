package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) registrationAfterPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_password"), slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	a.Log.Debug("kratos webhook registration_after_password", slog.String("identity_id", payload.Identity.Id))
	v := r.Header.Get("Authorization")
	if v != a.Cfg.KratosAPIKey {
		a.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_password"), slog.String("authorization", v))
		http.Error(w, "unauthorized", http.StatusBadRequest)
		return
	}
	if err := a.createUserRelationTuples(r.Context(), payload.Identity.Id); err != nil {
		a.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to insert user relation-tuple", http.StatusBadRequest)
		return
	}
	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		a.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if _, err := a.persistor.User().CreateUser(r.Context(), dbtype.UserSetter{ID: nullable.NewNullableWithValue(identityID)}); err != nil {
		a.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to create user", http.StatusBadRequest)
		return
	}

	// send welcome email...

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

func (a *ApiHandler) registrationAfterOidc(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Identity *orykratos.Identity `json:"identity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.Log.Error("invalid webhook payload", slog.String("webhook", "registration_after_oidc"), slog.Any("error", err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	a.Log.Debug("kratos webhook registration_after_oidc", slog.String("identity_id", payload.Identity.Id))
	v := r.Header.Get("Authorization")
	if v != a.Cfg.KratosAPIKey {
		a.Log.Error("webhook auth failed", slog.String("webhook", "registration_after_oidc"), slog.String("authorization", v))
		http.Error(w, "unauthorized", http.StatusBadRequest)
		return
	}
	if err := a.createUserRelationTuples(r.Context(), payload.Identity.Id); err != nil {
		a.Log.Error("failed to insert user relation-tuple", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to insert user relation-tuple", http.StatusBadRequest)
		return
	}
	identityID, err := uuid.Parse(payload.Identity.Id)
	if err != nil {
		a.Log.Error("failed to parse identity id", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if _, err := a.persistor.User().CreateUser(r.Context(), dbtype.UserSetter{ID: nullable.NewNullableWithValue(identityID)}); err != nil {
		a.Log.Error("failed to create new user", slog.String("identity_id", payload.Identity.Id), slog.Any("error", err))
		http.Error(w, "failed to create user", http.StatusBadRequest)
		return
	}

	// send welcome email...

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

func (a *ApiHandler) createUserRelationTuples(ctx context.Context, identityID string) error {
	_, err := a.Keto.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customer",
					Relation:  "members",
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    AuthzIdentityID(identityID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    AuthzIdentityID(identityID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities", ""),
				},
			},
		},
	})
	return err
}
