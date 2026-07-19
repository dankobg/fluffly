package petfinder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type ImportOrgsCmd struct {
	Dir     string `required:"" type:"path" help:"Downloaded seed data dir"`
	Min     int    `default:"1" help:"Starting page (inclusive)"`
	Max     int    `default:"1000" help:"Ending page (inclusive)"`
	Workers int    `default:"8" help:"Number of workers"`
}

type petfinderOrg struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Distance *string `json:"distance"`
	Address  struct {
		Address1 *string `json:"address1"`
		Address2 *string `json:"address2"`
		City     string  `json:"city"`
		State    *string `json:"state"`
		Postcode *string `json:"postcode"`
		Country  string  `json:"country"`
	} `json:"address"`
	Hours struct {
		Monday    *string `json:"monday"`
		Tuesday   *string `json:"tuesday"`
		Wednesday *string `json:"wednesday"`
		Thursday  *string `json:"thursday"`
		Friday    *string `json:"friday"`
		Saturday  *string `json:"saturday"`
		Sunday    *string `json:"sunday"`
	} `json:"hours"`
	URL              *string `json:"url"`
	Website          *string `json:"website"`
	MissionStatement *string `json:"mission_statement"`
	Adoption         struct {
		Policy *string `json:"policy"`
		URL    *string `json:"url"`
	} `json:"adoption"`
	SocialMedia struct {
		Facebook  *string `json:"facebook"`
		Twitter   *string `json:"twitter"`
		Youtube   *string `json:"youtube"`
		Instagram *string `json:"instagram"`
		Pinterest *string `json:"pinterest"`
	} `json:"social_media"`
	Photos []struct {
		Small  *string `json:"small"`
		Medium *string `json:"medium"`
		Large  *string `json:"large"`
		Full   *string `json:"full"`
	} `json:"photos"`
}

type petfinderOrgData struct {
	Organizations []petfinderOrg
}

type importOrgResult struct {
	Name string
	Err  error
}

var (
	seenCount    = make(map[string]int)
	mu           sync.Mutex
	orgIDSeq     int64
	addressIDSeq atomic.Int64
)

func (ic *ImportOrgsCmd) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	ketoClient, err := keto.NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	pages, err := filepath.Glob(filepath.Join(ic.Dir, "orgs", "page_*.json"))
	if err != nil {
		return fmt.Errorf("failed to glob files: %w", err)
	}

	if len(pages) == 0 {
		return fmt.Errorf("no page_*.json files found")
	}

	type copyResult struct {
		count int64
		err   error
	}

	var (
		jobs            = make(chan string)
		orgSetters      = make(chan models.OrganizationSetter, 1_000)
		contactSetters  = make(chan models.OrganizationContactSetter, 1_000)
		addressSetters  = make(chan models.AddressSetter, 1_000)
		workHourSetters = make(chan models.OrganizationWorkHourSetter, 1_000)
		photoSetters    = make(chan models.OrganizationPhotoSetter, 1_000)
		socialSetters   = make(chan models.OrganizationSocialSetter, 1_000)

		orgResults      = make(chan copyResult)
		contactResults  = make(chan copyResult)
		addressResults  = make(chan copyResult)
		workHourResults = make(chan copyResult)
		photoResults    = make(chan copyResult)
		socialResults   = make(chan copyResult)

		wg sync.WaitGroup
	)

	if err := disableConstraintsForOrganizationImport(ctx, pool); err != nil {
		return fmt.Errorf("disableConstraints: %w", err)
	}

	for range ic.Workers {
		wg.Go(func() {
			for file := range jobs {
				_ = processOrg(file, orgSetters, contactSetters, addressSetters, workHourSetters, photoSetters, socialSetters)
			}
		})
	}

	go func() {
		for _, page := range pages {
			n, err := pageNumber(page)
			if err != nil {
				continue
			}

			if n >= ic.Min && n <= ic.Max {
				jobs <- page
			}
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(orgSetters)
		close(contactSetters)
		close(addressSetters)
		close(workHourSetters)
		close(photoSetters)
		close(socialSetters)
	}()

	var wg2 sync.WaitGroup
	wg2.Add(6)

	go func() {
		defer func() {
			close(orgResults)
			wg2.Done()
		}()

		c, err := copyOrganizationData(ctx, pool, orgSetters)
		orgResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(contactResults)
			wg2.Done()
		}()

		c, err := copyOrganizationContactData(ctx, pool, contactSetters)
		contactResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(addressResults)
			wg2.Done()
		}()

		c, err := copyOrganizationContactAddressData(ctx, pool, addressSetters)
		addressResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(workHourResults)
			wg2.Done()
		}()

		c, err := copyOrganizationWorkHourData(ctx, pool, workHourSetters)
		workHourResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(photoResults)
			wg2.Done()
		}()

		c, err := copyOrganizationPhotoData(ctx, pool, photoSetters)
		photoResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(socialResults)
			wg2.Done()
		}()

		c, err := copyOrganizationSocialData(ctx, pool, socialSetters)
		socialResults <- copyResult{count: c, err: err}
	}()

	var totalErrors []error

	for r := range orgResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("org: %w", r.err))
		}
	}

	for r := range contactResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("contact: %w", r.err))
		}
	}

	for r := range addressResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("address: %w", r.err))
		}
	}

	for r := range workHourResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("work_hour: %w", r.err))
		}
	}

	for r := range photoResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("photo: %w", r.err))
		}
	}

	for r := range socialResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("social: %w", r.err))
		}
	}

	if len(totalErrors) > 0 {
		fmt.Println("total errors", len(totalErrors))
		return errors.Join(totalErrors...)
	}

	wg2.Wait()

	if err := enableConstraintsForOrganizationImport(ctx, pool); err != nil {
		return fmt.Errorf("enableConstraints: %w", err)
	}

	if err := setSequenceIDsAfterOrgImport(ctx, pool); err != nil {
		return fmt.Errorf("setSequenceIDs: %w", err)
	}

	if err := createOrganizationRelationTuples(ctx, ketoClient); err != nil {
		return fmt.Errorf("create organization tuples: %w", err)
	}

	return nil
}

func disableConstraintsForOrganizationImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("constraints dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "ALTER TABLE organization DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_contact DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE address DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_work_hour DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_photo DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_social DISABLE TRIGGER ALL")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("disable constraints tx.Commit: %w", err)
	}

	log.Println("disabled constraints before orgs bulk import")

	return nil
}

func enableConstraintsForOrganizationImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("constraints dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "ALTER TABLE organization ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_contact ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE address ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_work_hour ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_photo ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE organization_social ENABLE TRIGGER ALL")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("enable constraints tx.Commit: %w", err)
	}

	log.Println("enabled constraints after orgs bulk import")

	return nil
}

func setSequenceIDsAfterOrgImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("seqIDs dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "select setval('organization_id_seq', (SELECT MAX(id) FROM organization))")
	_, _ = tx.Exec(ctx, "select setval('organization_contact_id_seq', (SELECT MAX(id) FROM organization_contact))")
	_, _ = tx.Exec(ctx, "select setval('address_id_seq', (SELECT MAX(id) FROM address))")
	_, _ = tx.Exec(ctx, "select setval('organization_work_hour_id_seq', (SELECT MAX(id) FROM organization_work_hour))")
	_, _ = tx.Exec(ctx, "select setval('organization_photo_id_seq', (SELECT MAX(id) FROM organization_photo))")
	_, _ = tx.Exec(ctx, "select setval('organization_social_id_seq', (SELECT MAX(id) FROM organization_social))")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("seqIDs tx.Commit: %w", err)
	}

	return nil
}

func processOrg(
	filename string,
	orgSetters chan<- models.OrganizationSetter,
	contactSetters chan<- models.OrganizationContactSetter,
	addressSetters chan<- models.AddressSetter,
	workHourSetters chan<- models.OrganizationWorkHourSetter,
	photoSetters chan<- models.OrganizationPhotoSetter,
	socialSetters chan<- models.OrganizationSocialSetter,
) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("faied to read file: %w", err)
	}

	var data petfinderOrgData
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("faied to unmarshal petfinder org: %w", err)
	}

	for _, org := range data.Organizations {
		ocs := makeOrganizationCreateSetter(org)

		mu.Lock()
		seenCount[org.Name]++
		n := seenCount[org.Name]
		mu.Unlock()

		if n > 1 {
			ocs.Organization.Name = omit.From(fmt.Sprintf("%s - %d", org.Name, n))
		}

		orgSetters <- ocs.Organization

		contactSetters <- ocs.Contact

		addressSetters <- ocs.Address

		if !ocs.WorkHour.IsUnset() {
			workHourSetters <- ocs.WorkHour.MustGet()
		}

		if !ocs.Photos.IsUnset() {
			for _, x := range ocs.Photos.MustGet() {
				photoSetters <- x
			}
		}

		if !ocs.Socials.IsUnset() {
			for _, x := range ocs.Socials.MustGet() {
				socialSetters <- x
			}
		}
	}

	return nil
}

func copyOrganizationData(ctx context.Context, dbpool *pgxpool.Pool, orgSetters <-chan models.OrganizationSetter) (int64, error) {
	cols := []string{
		"name",
		"website",
		"mission_statement",
		"adoption_policy",
		"adoption_url",
		"distance",
		"status",
		"original_id",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"organization"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-orgSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.Name.MustGet(),
				valOrNil(v.Website),
				valOrNil(v.MissionStatement),
				valOrNil(v.AdoptionPolicy),
				valOrNil(v.AdoptionURL),
				valOrNil(v.Distance),
				"approved",
				valOrNil(v.OriginalID),
			}

			return row, nil
		}),
	)
}

func copyOrganizationContactData(ctx context.Context, dbpool *pgxpool.Pool, contactSetters <-chan models.OrganizationContactSetter) (int64, error) {
	cols := []string{
		"organization_id",
		"address_id",
		"phone",
		"email",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"organization_contact"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-contactSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.OrganizationID.MustGet(),
				v.AddressID.MustGet(),
				v.Phone.MustGet(),
				v.Email.MustGet(),
			}

			return row, nil
		}),
	)
}

func copyOrganizationContactAddressData(ctx context.Context, dbpool *pgxpool.Pool, addressSetters <-chan models.AddressSetter) (int64, error) {
	cols := []string{
		"country_id",
		"unit_number",
		"street_number",
		"street_address",
		"city",
		"region",
		"postal_code",
		"coords",
		"note",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"address"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-addressSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.CountryID.MustGet(),
				valOrNil(v.UnitNumber),
				valOrNil(v.StreetNumber),
				v.StreetAddress.MustGet(),
				v.City.MustGet(),
				valOrNil(v.Region),
				valOrNil(v.PostalCode),
				valOrNil(v.Coords),
				valOrNil(v.Note),
			}

			return row, nil
		}),
	)
}

func copyOrganizationWorkHourData(ctx context.Context, dbpool *pgxpool.Pool, workHourSetters <-chan models.OrganizationWorkHourSetter) (int64, error) {
	cols := []string{
		"organization_id",
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"organization_work_hour"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-workHourSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.OrganizationID),
				valOrNil(v.Monday),
				valOrNil(v.Tuesday),
				valOrNil(v.Wednesday),
				valOrNil(v.Thursday),
				valOrNil(v.Friday),
				valOrNil(v.Saturday),
				valOrNil(v.Sunday),
			}

			return row, nil
		}),
	)
}

func copyOrganizationPhotoData(ctx context.Context, dbpool *pgxpool.Pool, photoSetters <-chan models.OrganizationPhotoSetter) (int64, error) {
	cols := []string{
		"organization_id",
		"object_kind",
		"object_ref_small",
		"object_ref_medium",
		"object_ref_large",
		"object_ref_full",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"organization_photo"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-photoSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.OrganizationID),
				v.ObjectKind.MustGet(),
				valOrNil(v.ObjectRefSmall),
				valOrNil(v.ObjectRefMedium),
				valOrNil(v.ObjectRefLarge),
				valOrNil(v.ObjectRefFull),
			}

			return row, nil
		}),
	)
}

func copyOrganizationSocialData(ctx context.Context, dbpool *pgxpool.Pool, socialSetters <-chan models.OrganizationSocialSetter) (int64, error) {
	cols := []string{"organization_id", "platform", "url"}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"organization_social"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-socialSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.OrganizationID.MustGet(),
				v.Platform.MustGet(),
				v.URL.MustGet(),
			}

			return row, nil
		}),
	)
}

func pageNumber(path string) (int, error) {
	base := filepath.Base(path)

	var n int

	_, err := fmt.Sscanf(base, "page_%d.json", &n)

	return n, err
}

func makeOrganizationCreateSetter(org petfinderOrg) dbtype.OrganizationCreateSetter {
	var organizationCreateSetter dbtype.OrganizationCreateSetter

	orgID := atomic.AddInt64(&orgIDSeq, 1)
	nullableOrgID := omit.From(orgID)

	organizationCreateSetter.Organization = models.OrganizationSetter{
		ID:               nullableOrgID,
		Name:             omit.From(org.Name),
		Website:          omitnull.FromPtr(org.Website),
		MissionStatement: omitnull.FromPtr(org.MissionStatement),
		AdoptionPolicy:   omitnull.FromPtr(org.Adoption.Policy),
		AdoptionURL:      omitnull.FromPtr(org.Adoption.URL),
		Distance:         omitnull.FromPtr(org.Distance),
		OriginalID:       omitnull.From(org.ID),
	}

	addressID := addressIDSeq.Add(1)
	nullableAddressID := omit.From(addressID)

	randomFakeLatLon := pickRandomItem(fakeCoordsUSA[:])

	point := geom.NewPoint(geom.XY).SetSRID(dbcustom.SRID).MustSetCoords(geom.Coord{randomFakeLatLon[1], randomFakeLatLon[0]})
	coords := &ewkb.Point{Point: point}

	addressSetter := models.AddressSetter{
		ID:         nullableAddressID,
		City:       omit.From(org.Address.City),
		UnitNumber: omitnull.FromPtr(org.Address.Address2),
		Coords:     omitnull.From(coords),
		Region:     omitnull.FromPtr(org.Address.State),
		PostalCode: omitnull.FromPtr(org.Address.Postcode),
	}
	switch org.Address.Country {
	case "US":
		addressSetter.CountryID = omit.From[int64](236)
	case "CA":
		addressSetter.CountryID = omit.From[int64](41)
	case "MX":
		addressSetter.CountryID = omit.From[int64](144)
	default:
		addressSetter.CountryID = omit.From[int64](236)
	}

	if org.Address.Address1 != nil {
		addressSetter.StreetAddress = omit.From(*org.Address.Address1)
	} else {
		addressSetter.StreetAddress = omit.From(faker.GetRealAddress().Address)
	}

	organizationCreateSetter.Address = addressSetter

	contactSetter := models.OrganizationContactSetter{OrganizationID: nullableOrgID, AddressID: nullableAddressID}
	if strings.TrimSpace(org.Email) != "" {
		contactSetter.Email = omit.From(org.Email)
	} else {
		contactSetter.Email = omit.From(faker.Email())
	}

	if org.Phone != nil {
		contactSetter.Phone = omit.From(*org.Phone)
	} else {
		contactSetter.Phone = omit.From(faker.Phonenumber())
	}

	organizationCreateSetter.Contact = contactSetter

	var hasWorkHours bool

	workHourSetter := models.OrganizationWorkHourSetter{OrganizationID: omitnull.From(orgID)}

	if org.Hours.Monday != nil && strings.TrimSpace(*org.Hours.Monday) != "" {
		hasWorkHours = true
		workHourSetter.Monday = omitnull.From(*org.Hours.Monday)
	}

	if org.Hours.Tuesday != nil && strings.TrimSpace(*org.Hours.Tuesday) != "" {
		hasWorkHours = true
		workHourSetter.Tuesday = omitnull.From(*org.Hours.Tuesday)
	}

	if org.Hours.Wednesday != nil && strings.TrimSpace(*org.Hours.Wednesday) != "" {
		hasWorkHours = true
		workHourSetter.Wednesday = omitnull.From(*org.Hours.Wednesday)
	}

	if org.Hours.Thursday != nil && strings.TrimSpace(*org.Hours.Thursday) != "" {
		hasWorkHours = true
		workHourSetter.Thursday = omitnull.From(*org.Hours.Thursday)
	}

	if org.Hours.Friday != nil && strings.TrimSpace(*org.Hours.Friday) != "" {
		hasWorkHours = true
		workHourSetter.Friday = omitnull.From(*org.Hours.Friday)
	}

	if org.Hours.Saturday != nil && strings.TrimSpace(*org.Hours.Saturday) != "" {
		hasWorkHours = true
		workHourSetter.Saturday = omitnull.From(*org.Hours.Saturday)
	}

	if org.Hours.Sunday != nil && strings.TrimSpace(*org.Hours.Sunday) != "" {
		hasWorkHours = true
		workHourSetter.Sunday = omitnull.From(*org.Hours.Sunday)
	}

	if hasWorkHours {
		organizationCreateSetter.WorkHour = omitnull.From(workHourSetter)
	}

	if len(org.Photos) > 0 {
		organizationPhotoSetters := make([]models.OrganizationPhotoSetter, len(org.Photos))
		for i, photo := range org.Photos {
			organizationPhotoSetters[i] = models.OrganizationPhotoSetter{
				OrganizationID:  omitnull.From(orgID),
				ObjectKind:      omit.From(media.StorageKindExternal),
				ObjectRefSmall:  omitnull.FromPtr(photo.Small),
				ObjectRefMedium: omitnull.FromPtr(photo.Medium),
				ObjectRefLarge:  omitnull.FromPtr(photo.Large),
				ObjectRefFull:   omitnull.FromPtr(photo.Full),
			}
		}

		organizationCreateSetter.Photos = omitnull.From(organizationPhotoSetters)
	}

	// their seed data has no videos for organizations

	organizationSocialsSetters := make([]models.OrganizationSocialSetter, 0)

	var hasSocials bool
	if org.SocialMedia.Facebook != nil {
		hasSocials = true

		organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
			OrganizationID: nullableOrgID,
			Platform:       omit.From("Facebook"),
			URL:            omit.From(*org.SocialMedia.Facebook),
		})
	}

	if org.SocialMedia.Instagram != nil {
		hasSocials = true

		organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
			OrganizationID: nullableOrgID,
			Platform:       omit.From("Instagram"),
			URL:            omit.From(*org.SocialMedia.Instagram),
		})
	}

	if org.SocialMedia.Pinterest != nil {
		hasSocials = true

		organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
			OrganizationID: nullableOrgID,
			Platform:       omit.From("Pinterest"),
			URL:            omit.From(*org.SocialMedia.Pinterest),
		})
	}

	if org.SocialMedia.Twitter != nil {
		hasSocials = true

		organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
			OrganizationID: nullableOrgID,
			Platform:       omit.From("Twitter"),
			URL:            omit.From(*org.SocialMedia.Twitter),
		})
	}

	if org.SocialMedia.Youtube != nil {
		hasSocials = true

		organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
			OrganizationID: nullableOrgID,
			Platform:       omit.From("Youtube"),
			URL:            omit.From(*org.SocialMedia.Youtube),
		})
	}

	if hasSocials {
		organizationCreateSetter.Socials = omitnull.From(organizationSocialsSetters)
	}

	return organizationCreateSetter
}

func createOrganizationRelationTuples(ctx context.Context, c *keto.Client) error {
	const batchSize = 5_000

	for start := int64(0); start < orgIDSeq; start += batchSize {
		end := min(start+batchSize, orgIDSeq)

		relationTupleDeltas := make([]*rts.RelationTupleDelta, 0, end-start)

		for i := start; i < end; i++ {
			relationTupleDeltas = append(relationTupleDeltas, &rts.RelationTupleDelta{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Organization",
					Object:    shared.AuthzOrganizationID(i + 1),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Organizations", "organizations", ""),
				},
			})
		}

		if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
			RelationTupleDeltas: relationTupleDeltas,
		}); err != nil {
			return fmt.Errorf("create tuple batch [%d:%d]: %w", start, end, err)
		}
	}

	return nil
}
