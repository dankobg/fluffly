package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/bob/orm"
	"github.com/stephenafamo/scan"
)

var _ persistence.OrganizationPersistor = (*PgOrganizationPersistor)(nil)

type PgOrganizationPersistor struct {
	*PgPersistor
}

func NewPgOrganizationPersistor(ps *PgPersistor) *PgOrganizationPersistor {
	return &PgOrganizationPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrOrganizationIntegrityViolation  struct{ errIntegrityViolation }
	ErrOrganizationUniqueViolation     struct{ errUniqueViolation }
	ErrOrganizationForeignKeyViolation struct{ errForeignKeyViolation }
	ErrOrganizationCheckViolation      struct{ errCheckViolation }
)

var (
	ErrOrganizationNotFound                         = errors.New("organization not found")
	ErrOrganizationSocialNotFound                   = errors.New("organization social not found")
	ErrOrganizationPhotoNotFound                    = errors.New("organization photo not found")
	ErrOrganizationVideoNotFound                    = errors.New("organization video not found")
	errOrganizationIntegrity                        = ErrOrganizationIntegrityViolation{}
	errOrganizationUniqueName                       = ErrOrganizationUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errOrganizationContactForeignKeyOrganizationID  = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "contact.organization_id"}}
	errOrganizationContactForeignKeyAddressID       = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "contact.address_id"}}
	errOrganizationAddressForeignKeyCountryID       = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "address.country_id"}}
	errOrganizationWorkHourUniqueOrganizationID     = ErrOrganizationUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "work_hour.organization_id"}}
	errOrganizationWorkHourForeignKeyOrganizationID = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "work_hour.organization_id"}}
	errOrganizationWorkHourCheckDayProvided         = ErrOrganizationCheckViolation{errCheckViolation: errCheckViolation{Name: "work_hour.day_provided"}}
	errOrganizationPhotoForeignKeyOrganizationID    = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "photo.organization_id"}}
	errOrganizationSocialForeignKeyOrganizationID   = ErrOrganizationForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "social.organization_id"}}
)

func convertOrganizationPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_organization_name":
			return errOrganizationUniqueName
		case "fk_address_country_id":
			return errOrganizationAddressForeignKeyCountryID
		case "fk_organization_contact_organization_id":
			return errOrganizationContactForeignKeyOrganizationID
		case "fk_organization_contact_address_id":
			return errOrganizationContactForeignKeyAddressID
		case "uq_organization_work_hour_organization_id":
			return errOrganizationWorkHourUniqueOrganizationID
		case "fk_organization_work_hour_organization_id":
			return errOrganizationWorkHourForeignKeyOrganizationID
		case "ck_organization_work_hour_provided":
			return errOrganizationWorkHourCheckDayProvided
		case "fk_organization_photo_organization_id":
			return errOrganizationPhotoForeignKeyOrganizationID
		case "fk_organization_social_organization_id":
			return errOrganizationSocialForeignKeyOrganizationID
		}

		return errOrganizationIntegrity
	}

	return pgErr
}

func (pst *PgOrganizationPersistor) ListOrganizations(ctx context.Context, filters dbtype.ListOrganizationsFilters) (dbtype.PagedResult[dbtype.OrganizationWithJoinData], error) {
	q := psql.Select(
		sm.Columns(models.Organizations.Columns),
		sm.From(models.Organizations.Name()),
		sm.GroupBy(models.Organizations.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Organizations.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListOrganizationsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListOrganizationsParamsEmbedContact:
					q.Apply(sm.Columns(
						models.OrganizationContacts.Columns.WithPrefix("contact."),
						models.Addresses.Columns.Except("coords").WithPrefix("contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "contact_address.coords"`),
						models.Countries.Columns.WithPrefix("contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))

				case api.ListOrganizationsParamsEmbedWorkHour:
					q.Apply(sm.Columns(
						models.OrganizationWorkHours.Columns.WithPrefix("work_hour."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))

				case api.ListOrganizationsParamsEmbedPhotos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(photos_row))
  FROM
    (
      SELECT
        organization_photo.id AS "id",
        organization_photo.organization_id AS "organizationID",
        organization_photo.object_kind AS "objectKind",
        organization_photo.object_ref_small AS "objectRefSmall",
        organization_photo.object_ref_medium AS "objectRefMedium",
        organization_photo.object_ref_large AS "objectRefLarge",
        organization_photo.object_ref_full AS "objectRefFull",
        organization_photo.created_at AS "createdAt",
        organization_photo.updated_at AS "updatedAt"
      FROM
        public.organization_photo
      WHERE
        organization_photo.organization_id = organization.id
    ) AS photos_row
) AS "photos"`),
					))

				case api.ListOrganizationsParamsEmbedVideos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(videos_row))
  FROM
    (
      SELECT
        organization_video.id AS "id",
        organization_video.organization_id AS "organizationID",
        organization_video.object_kind AS "objectKind",
        organization_video.object_ref AS "objectRef",
        organization_video.created_at AS "createdAt",
        organization_video.updated_at AS "updatedAt"
      FROM
        public.organization_video
      WHERE
        organization_video.organization_id = organization.id
    ) AS videos_row
) AS "videos"`),
					))

				case api.ListOrganizationsParamsEmbedSocials:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(socials_row))
  FROM
    (
      SELECT
        organization_social.id AS "id",
        organization_social.organization_id AS "organizationID",
        organization_social.platform AS "platform",
        organization_social.url AS "url",
        organization_social.created_at AS "createdAt",
        organization_social.updated_at AS "updatedAt"
      FROM
        public.organization_social
      WHERE
        organization_social.organization_id = organization.id
    ) AS socials_row
) AS "socials"`),
					))

				}
			}
		}

		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Organizations.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.Organizations.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}

		if filters.Status != nil {
			statuses := make([]any, len(*filters.Status))
			for i, id := range *filters.Status {
				statuses[i] = id
			}

			q.Apply(sm.Where(models.Organizations.Columns.Status.In(psql.Arg(statuses...))))
		}
	}

	type ListOrganizationsRow struct {
		dbtype.OrganizationWithJoinData
		TotalCount int64
	}

	organizations, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListOrganizationsRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.PagedResult[dbtype.OrganizationWithJoinData]{}, fmt.Errorf("query organizations")
	}

	result := dbtype.PagedResult[dbtype.OrganizationWithJoinData]{
		Data: make([]dbtype.OrganizationWithJoinData, len(organizations)),
	}
	for i, row := range organizations {
		result.Data[i] = row.OrganizationWithJoinData
	}

	if len(organizations) > 0 {
		result.TotalCount = organizations[0].TotalCount
	}

	return result, nil
}

func (pst *PgOrganizationPersistor) GetOrganizationByID(ctx context.Context, organizationID int64, filters dbtype.GetOrganizationByIDFilters) (dbtype.OrganizationWithJoinData, error) {
	q := psql.Select(
		sm.Columns(models.Organizations.Columns),
		sm.From(models.Organizations.Name()),
		sm.Where(models.Organizations.Columns.ID.EQ(psql.Arg(organizationID))),
		sm.GroupBy(models.Organizations.Columns.ID),
	)

	if hasAnyLogicFilters(&filters.GetOrganizationParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.GetOrganizationParamsEmbedContact:
					q.Apply(sm.Columns(
						models.OrganizationContacts.Columns.WithPrefix("contact."),
						models.Addresses.Columns.Except("coords").WithPrefix("contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "contact_address.coords"`),
						models.Countries.Columns.WithPrefix("contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))

				case api.GetOrganizationParamsEmbedWorkHour:
					q.Apply(sm.Columns(
						models.OrganizationWorkHours.Columns.WithPrefix("work_hour."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))

				case api.GetOrganizationParamsEmbedPhotos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(photos_row))
  FROM
    (
      SELECT
        organization_photo.id AS "id",
        organization_photo.organization_id AS "organizationID",
        organization_photo.object_kind AS "objectKind",
        organization_photo.object_ref_small AS "objectRefSmall",
        organization_photo.object_ref_medium AS "objectRefMedium",
        organization_photo.object_ref_large AS "objectRefLarge",
        organization_photo.object_ref_full AS "objectRefFull",
        organization_photo.created_at AS "createdAt",
        organization_photo.updated_at AS "updatedAt"
      FROM
        public.organization_photo
      WHERE
        organization_photo.organization_id = organization.id
    ) AS photos_row
) AS "photos"`),
					))

				case api.GetOrganizationParamsEmbedVideos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(videos_row))
  FROM
    (
      SELECT
        organization_video.id AS "id",
        organization_video.organization_id AS "organizationID",
        organization_video.object_kind AS "objectKind",
        organization_video.object_ref AS "objectRef",
        organization_video.created_at AS "createdAt",
        organization_video.updated_at AS "updatedAt"
      FROM
        public.organization_video
      WHERE
        organization_video.organization_id = organization.id
    ) AS videos_row
) AS "videos"`),
					))

				case api.GetOrganizationParamsEmbedSocials:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(socials_row))
  FROM
    (
      SELECT
        organization_social.id AS "id",
        organization_social.organization_id AS "organizationID",
        organization_social.platform AS "platform",
        organization_social.url AS "url",
        organization_social.created_at AS "createdAt",
        organization_social.updated_at AS "updatedAt"
      FROM
        public.organization_social
      WHERE
        organization_social.organization_id = organization.id
    ) AS socials_row
) AS "socials"`),
					))

				}
			}
		}
	}

	orgWithJoinData, err := bob.One(ctx, pst.exec, q, scan.StructMapper[dbtype.OrganizationWithJoinData](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.OrganizationWithJoinData{}, fmt.Errorf("query country")
	}

	return orgWithJoinData, nil
}

func (pst *PgOrganizationPersistor) ApplyForOrganization(ctx context.Context, userID uuid.UUID, in dbtype.OrganizationApplyForSetter) (models.Organization, error) {
	var insertedOrganization models.Organization

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Organizations.Insert(&in.Organization, im.Returning(models.Organizations.Columns))

		org, err := bob.One(ctx, tx, q1, scan.StructMapper[models.Organization]())
		if err != nil {
			return fmt.Errorf("insert organization: %w", err)
		}

		insertedOrganization = org

		// q2 := models.Addresses.Insert(&in.Address, im.Returning(models.Addresses.Columns.ID))
		q2 := psql.Insert(
			im.Into(models.Addresses.Name(), "country_id", "unit_number", "street_number", "street_address", "city", "region", "postal_code", "coords", "note"),
			im.Values(
				psql.Arg(in.Address.CountryID.MustGet()),
				psql.Arg(valOrNil(in.Address.UnitNumber)),
				psql.Arg(valOrNil(in.Address.StreetNumber)),
				psql.Arg(in.Address.StreetAddress.MustGet()),
				psql.Arg(in.Address.City.MustGet()),
				psql.Arg(valOrNil(in.Address.Region)),
				psql.Arg(valOrNil(in.Address.PostalCode)),
				psql.Raw("ST_GeomFromEWKB(?)", valOrNil(in.Address.Coords)),
				psql.Arg(valOrNil(in.Address.Note)),
			),
			im.Returning(models.Addresses.Columns.ID),
		)

		insertedAddress, err := bob.One(ctx, tx, q2, scan.StructMapper[models.Address]())
		if err != nil {
			return fmt.Errorf("insert organization address: %w", err)
		}

		in.Contact.OrganizationID.Set(insertedOrganization.ID)
		in.Contact.AddressID.Set(insertedAddress.ID)

		q3 := models.OrganizationContacts.Insert(&in.Contact, im.Returning(models.OrganizationContacts.Columns.ID))
		if _, err := bob.One(ctx, tx, q3, scan.StructMapper[models.OrganizationContact]()); err != nil {
			return fmt.Errorf("insert organization contact: %w", err)
		}

		if !in.WorkHour.IsUnset() && !in.WorkHour.IsNull() {
			inWorkHour := in.WorkHour.MustGet()
			inWorkHour.OrganizationID.Set(insertedOrganization.ID)

			q4 := models.OrganizationWorkHours.Insert(&inWorkHour, im.Returning(models.OrganizationWorkHours.Columns.ID))
			if _, err := bob.One(ctx, tx, q4, scan.StructMapper[models.OrganizationWorkHour]()); err != nil {
				return fmt.Errorf("insert organization work hour: %w", err)
			}
		}

		if !in.Photos.IsUnset() && !in.Photos.IsNull() {
			inPhotos := in.Photos.MustGet()
			if len(inPhotos) > 0 {
				setters := make([]*models.OrganizationPhotoSetter, len(inPhotos))
				for i, x := range inPhotos {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q5 := models.OrganizationPhotos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationPhotos.Columns))
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("insert organization photos: %w", err)
				}
			}
		}

		if !in.Videos.IsUnset() && !in.Videos.IsNull() {
			inVideos := in.Videos.MustGet()
			if len(inVideos) > 0 {
				setters := make([]*models.OrganizationVideoSetter, len(inVideos))
				for i, x := range inVideos {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q6 := models.OrganizationVideos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationVideos.Columns))
				if _, err := bob.Exec(ctx, tx, q6); err != nil {
					return fmt.Errorf("insert organization videos: %w", err)
				}
			}
		}

		if !in.Socials.IsUnset() && !in.Socials.IsNull() {
			inSocials := in.Socials.MustGet()
			if len(inSocials) > 0 {
				setters := make([]*models.OrganizationSocialSetter, len(inSocials))
				for i, x := range inSocials {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q7 := models.OrganizationSocials.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationSocials.Columns))
				if _, err := bob.Exec(ctx, tx, q7); err != nil {
					return fmt.Errorf("insert organization socials: %w", err)
				}
			}
		}

		var membershipSetter models.OrganizationMembershipSetter
		membershipSetter.OrganizationID.Set(insertedOrganization.ID)
		membershipSetter.UserID.Set(userID)

		q8 := models.OrganizationMemberships.Insert(&membershipSetter)
		if _, err := bob.Exec(ctx, tx, q8); err != nil {
			return fmt.Errorf("insert organization membership: %w", err)
		}

		return nil
	})

	return insertedOrganization, txErr
}

func (pst *PgOrganizationPersistor) ApproveOrganization(ctx context.Context, organizationID int64) error {
	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Organizations.Update(
			models.OrganizationSetter{Status: omit.From("approved")}.UpdateMod(),
			um.Where(models.Organizations.Columns.ID.EQ(psql.Arg(organizationID))),
			um.Returning(models.Organizations.Columns),
		)
		if _, err := bob.Exec(ctx, tx, q1); err != nil {
			return fmt.Errorf("update organization status to approved: %w", err)
		}

		q2 := models.OrganizationMemberships.Update(
			models.OrganizationMembershipSetter{Status: omit.From("approved")}.UpdateMod(),
			um.Where(models.OrganizationMemberships.Columns.OrganizationID.EQ(psql.Arg(organizationID))),
			um.Returning(models.OrganizationMemberships.Columns),
		)
		if _, err := bob.Exec(ctx, tx, q2); err != nil {
			return fmt.Errorf("update organization membership status to approved: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

func (pst *PgOrganizationPersistor) RejectOrganization(ctx context.Context, organizationID int64) error {
	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Organizations.Update(
			models.OrganizationSetter{Status: omit.From("rejected")}.UpdateMod(),
			um.Where(models.Organizations.Columns.ID.EQ(psql.Arg(organizationID))),
			um.Returning(models.Organizations.Columns),
		)
		if _, err := bob.Exec(ctx, tx, q1); err != nil {
			return fmt.Errorf("update organization status to rejected: %w", err)
		}

		q2 := models.OrganizationMemberships.Update(
			models.OrganizationMembershipSetter{Status: omit.From("rejected")}.UpdateMod(),
			um.Where(models.OrganizationMemberships.Columns.OrganizationID.EQ(psql.Arg(organizationID))),
			um.Returning(models.OrganizationMemberships.Columns),
		)
		if _, err := bob.Exec(ctx, tx, q2); err != nil {
			return fmt.Errorf("update organization membership status to rejected: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

func (pst *PgOrganizationPersistor) CreateOrganization(ctx context.Context, in dbtype.OrganizationCreateSetter) (models.Organization, error) {
	var insertedOrganization models.Organization

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Organizations.Insert(&in.Organization, im.Returning(models.Organizations.Columns))

		org, err := bob.One(ctx, tx, q1, scan.StructMapper[models.Organization]())
		if err != nil {
			return fmt.Errorf("insert organization: %w", err)
		}

		insertedOrganization = org

		// q2 := models.Addresses.Insert(&in.Address, im.Returning(models.Addresses.Columns.ID))
		q2 := psql.Insert(
			im.Into(models.Addresses.Name(), "country_id", "unit_number", "street_number", "street_address", "city", "region", "postal_code", "coords", "note"),
			im.Values(
				psql.Arg(in.Address.CountryID.MustGet()),
				psql.Arg(valOrNil(in.Address.UnitNumber)),
				psql.Arg(valOrNil(in.Address.StreetNumber)),
				psql.Arg(in.Address.StreetAddress.MustGet()),
				psql.Arg(in.Address.City.MustGet()),
				psql.Arg(valOrNil(in.Address.Region)),
				psql.Arg(valOrNil(in.Address.PostalCode)),
				psql.Raw("ST_GeomFromEWKB(?)", valOrNil(in.Address.Coords)),
				psql.Arg(valOrNil(in.Address.Note)),
			),
			im.Returning(models.Addresses.Columns.ID),
		)

		insertedAddress, err := bob.One(ctx, tx, q2, scan.StructMapper[models.Address]())
		if err != nil {
			return fmt.Errorf("insert organization address: %w", err)
		}

		in.Contact.OrganizationID.Set(insertedOrganization.ID)
		in.Contact.AddressID.Set(insertedAddress.ID)

		q3 := models.OrganizationContacts.Insert(&in.Contact, im.Returning(models.OrganizationContacts.Columns.ID))
		if _, err := bob.One(ctx, tx, q3, scan.StructMapper[models.OrganizationContact]()); err != nil {
			return fmt.Errorf("insert organization contact: %w", err)
		}

		if !in.WorkHour.IsUnset() && !in.WorkHour.IsNull() {
			inWorkHour := in.WorkHour.MustGet()
			inWorkHour.OrganizationID.Set(insertedOrganization.ID)

			q4 := models.OrganizationWorkHours.Insert(&inWorkHour, im.Returning(models.OrganizationWorkHours.Columns.ID))
			if _, err := bob.One(ctx, tx, q4, scan.StructMapper[models.OrganizationWorkHour]()); err != nil {
				return fmt.Errorf("insert organization work hour: %w", err)
			}
		}

		if !in.Photos.IsUnset() && !in.Photos.IsNull() {
			inPhotos := in.Photos.MustGet()
			if len(inPhotos) > 0 {
				setters := make([]*models.OrganizationPhotoSetter, len(inPhotos))
				for i, x := range inPhotos {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q5 := models.OrganizationPhotos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationPhotos.Columns))
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("insert organization photos: %w", err)
				}
			}
		}

		if !in.Videos.IsUnset() && !in.Videos.IsNull() {
			inVideos := in.Videos.MustGet()
			if len(inVideos) > 0 {
				setters := make([]*models.OrganizationVideoSetter, len(inVideos))
				for i, x := range inVideos {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q6 := models.OrganizationVideos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationVideos.Columns))
				if _, err := bob.Exec(ctx, tx, q6); err != nil {
					return fmt.Errorf("insert organization videos: %w", err)
				}
			}
		}

		if !in.Socials.IsUnset() && !in.Socials.IsNull() {
			inSocials := in.Socials.MustGet()
			if len(inSocials) > 0 {
				setters := make([]*models.OrganizationSocialSetter, len(inSocials))
				for i, x := range inSocials {
					x.OrganizationID.Set(insertedOrganization.ID)
					setters[i] = &x
				}

				q7 := models.OrganizationSocials.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationSocials.Columns))
				if _, err := bob.Exec(ctx, tx, q7); err != nil {
					return fmt.Errorf("insert organization socials: %w", err)
				}
			}
		}

		return nil
	})

	return insertedOrganization, txErr
}

func (pst *PgOrganizationPersistor) UpdateOrganization(ctx context.Context, organizationID int64, in dbtype.OrganizationUpdateSetter) (models.Organization, error) {
	q := psql.Select(
		sm.Columns(
			models.OrganizationContacts.Columns.ID.As("contact_id"),
			models.OrganizationWorkHours.Columns.ID.As("work_hour_id"),
			models.Addresses.Columns.ID.As("address_id"),
		),
		sm.From(models.Organizations.Name()),
		sm.LeftJoin(models.OrganizationContacts.Name()).
			On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
		sm.LeftJoin(models.Addresses.Name()).
			On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
		sm.LeftJoin(models.OrganizationWorkHours.Name()).
			On(models.OrganizationWorkHours.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
		sm.Where(models.Organizations.Columns.ID.EQ(psql.Arg(organizationID))),
	)

	type idDestResult struct {
		ContactID  *int64 `db:"contact_id"`
		AddressID  *int64 `db:"address_id"`
		WorkHourID *int64 `db:"work_hour_id"`
	}

	idDest, err := bob.One(ctx, pst.exec, q, scan.StructMapper[idDestResult]())
	if err != nil {
		return models.Organization{}, fmt.Errorf("query organization related ids upfront: %w", err)
	}

	var updatedOrganization models.Organization

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		if !in.Organization.IsUnset() && !in.Organization.IsNull() {
			inOrganization := in.Organization.MustGet()
			q1 := models.Organizations.Update(inOrganization.UpdateMod(), um.Returning(models.Organizations.Columns))

			org, err := bob.One(ctx, tx, q1, scan.StructMapper[models.Organization]())
			if err != nil {
				return fmt.Errorf("update organization: %w", err)
			}

			updatedOrganization = org
		}

		if !in.Address.IsUnset() && !in.Address.IsNull() {
			inAddress := in.Address.MustGet()
			// 	q2 := models.Addresses.Update(
			// 		inAddress.UpdateMod(),
			// 		um.Where(models.Addresses.Columns.ID.EQ(psql.Arg(organizationID))),
			// 		um.Returning(models.Addresses.Columns),
			// 	)

			// @TODO: hack untill i can directly pass the postgis `geography` type
			q2 := psql.Update(
				um.Table(models.Addresses.Name()),
				um.Where(models.Addresses.Columns.ID.EQ(psql.Arg(*idDest.AddressID))),
			)
			if inAddress.CountryID.IsValue() {
				q2.Apply(um.SetCol("country_id").ToArg(inAddress.CountryID.MustGet()))
			}

			if !inAddress.UnitNumber.IsUnset() {
				q2.Apply(um.SetCol("unit_number").ToArg(valOrNil(inAddress.UnitNumber)))
			}

			if !inAddress.StreetNumber.IsUnset() {
				q2.Apply(um.SetCol("street_number").ToArg(valOrNil(inAddress.StreetNumber)))
			}

			if inAddress.StreetAddress.IsValue() {
				q2.Apply(um.SetCol("street_address").ToArg(inAddress.StreetAddress.MustGet()))
			}

			if inAddress.City.IsValue() {
				q2.Apply(um.SetCol("city").ToArg(inAddress.City.MustGet()))
			}

			if !inAddress.Region.IsUnset() {
				q2.Apply(um.SetCol("region").ToArg(valOrNil(inAddress.Region)))
			}

			if !inAddress.PostalCode.IsUnset() {
				q2.Apply(um.SetCol("postal_code").ToArg(valOrNil(inAddress.PostalCode)))
			}

			if !inAddress.Coords.IsUnset() {
				q2.Apply(um.SetCol("coords").To(psql.Raw("ST_GeomFromEWKB(?)", valOrNil(inAddress.Coords))))
			}

			if !inAddress.Note.IsUnset() {
				q2.Apply(um.SetCol("note").ToArg(valOrNil(inAddress.Note)))
			}

			if _, err := bob.Exec(ctx, tx, q2); err != nil {
				return fmt.Errorf("update organization address: %w", err)
			}

			if _, err := bob.Exec(ctx, tx, q2); err != nil {
				return fmt.Errorf("update organization address: %w", err)
			}
		}

		if !in.Contact.IsUnset() && !in.Contact.IsNull() {
			inContact := in.Contact.MustGet()

			q3 := models.OrganizationContacts.Update(
				inContact.UpdateMod(),
				um.Where(models.OrganizationContacts.Columns.ID.EQ(psql.Arg(*idDest.ContactID))),
				um.Returning(models.OrganizationContacts.Columns),
			)
			if _, err := bob.Exec(ctx, tx, q3); err != nil {
				return fmt.Errorf("update organization contact: %w", err)
			}
		}

		if !in.WorkHour.IsUnset() {
			if in.WorkHour.IsNull() {
				if idDest.WorkHourID != nil {
					q4 := models.OrganizationWorkHours.Delete(dm.Where(models.OrganizationWorkHours.Columns.ID.EQ(psql.Arg(*idDest.WorkHourID))))
					if _, err := bob.Exec(ctx, pst.exec, q4); err != nil {
						return fmt.Errorf("update organization work hour: %w", err)
					}
				}
			} else {
				inWorkHour := in.WorkHour.MustGet()

				q5 := models.OrganizationWorkHours.Insert(
					&inWorkHour,
					im.Returning(models.OrganizationWorkHours.Columns),
					im.OnConflict(models.OrganizationWorkHours.Columns.OrganizationID).
						DoUpdate(im.SetExcluded(models.OrganizationWorkHours.Columns.Except(
							models.OrganizationWorkHours.Columns.ID.String(),
							models.OrganizationWorkHours.Columns.OrganizationID.String(),
							models.OrganizationWorkHours.Columns.CreatedAt.String(),
						).Names()...)),
				)
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("update organization work hour: %w", err)
				}
			}
		}

		if !in.Photos.IsUnset() && !in.Photos.IsNull() {
			inPhotos := in.Photos.MustGet()
			if len(inPhotos) > 0 {
				setters := make([]*models.OrganizationPhotoSetter, len(inPhotos))
				for i, x := range inPhotos {
					x.OrganizationID.Set(updatedOrganization.ID)
					setters[i] = &x
				}

				q5 := models.OrganizationPhotos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationPhotos.Columns))
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("insert organization photos: %w", err)
				}
			}
		}

		if !in.Videos.IsUnset() && !in.Videos.IsNull() {
			inVideos := in.Videos.MustGet()
			if len(inVideos) > 0 {
				setters := make([]*models.OrganizationVideoSetter, len(inVideos))
				for i, x := range inVideos {
					x.OrganizationID.Set(updatedOrganization.ID)
					setters[i] = &x
				}

				q6 := models.OrganizationVideos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationVideos.Columns))
				if _, err := bob.Exec(ctx, tx, q6); err != nil {
					return fmt.Errorf("insert organization videos: %w", err)
				}
			}
		}

		if !in.Socials.IsUnset() {
			deleteOld := in.Socials.IsNull() || (!in.Socials.IsNull() && len(in.Socials.MustGet()) > 0)
			if deleteOld {
				q7 := models.OrganizationSocials.Delete(dm.Where(models.OrganizationSocials.Columns.OrganizationID.EQ(psql.Arg(organizationID))))
				if _, err := bob.Exec(ctx, pst.exec, q7); err != nil {
					return fmt.Errorf("delete old organization social platforms: %w", err)
				}
			}

			if !in.Socials.IsNull() && len(in.Socials.MustGet()) > 0 {
				inSocials := in.Socials.MustGet()

				setters := make([]*models.OrganizationSocialSetter, len(inSocials))
				for i, x := range inSocials {
					x.OrganizationID.Set(updatedOrganization.ID)
					setters[i] = &x
				}
				// @TODO: should probably do upsert but i am lazy,
				// i just delete old and add new, no need to fetch old, diff, upsert, delete etc...
				q8 := models.OrganizationSocials.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationSocials.Columns))
				if _, err := bob.Exec(ctx, tx, q8); err != nil {
					return fmt.Errorf("insert organization social platforms: %w", err)
				}
			}
		}

		return nil
	})
	if txErr != nil {
		return models.Organization{}, nil
	}

	return updatedOrganization, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationByID(ctx context.Context, organizationID int64) (int64, error) {
	q := models.Organizations.Delete(dm.Where(models.Organizations.Columns.ID.EQ(psql.Arg(organizationID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete organization: %w", err)
	}

	return organizationID, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizations(ctx context.Context, ids []int64) error {
	organizationIDs := make([]any, len(ids))
	for i, id := range ids {
		organizationIDs[i] = id
	}

	q := models.Organizations.Delete(dm.Where(models.Organizations.Columns.ID.In(psql.Arg(organizationIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete organizations: %w", err)
	}

	return nil
}
