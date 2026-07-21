package petfinder

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/db/gen/enums"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/random"
	"github.com/dankobg/fluffly/shared"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"github.com/stephenafamo/bob/types"
)

type ImportAnimalsCmd struct {
	Dir     string `required:"" type:"path" help:"Downloaded seed data dir"`
	Min     int    `default:"1" help:"Starting page (inclusive)"`
	Max     int    `default:"1000" help:"Ending page (inclusive)"`
	Workers int    `default:"8" help:"Number of workers"`
}

type petfinderAnimal struct {
	ID             int     `json:"id"`
	OrganizationID *string `json:"organization_id"`
	URL            *string `json:"url"`
	Type           string  `json:"type"`
	Species        string  `json:"species"`
	Breeds         struct {
		Primary   *string `json:"primary"`
		Secondary *string `json:"secondary"`
		Mixed     bool    `json:"mixed"`
		Unknown   bool    `json:"unknown"`
	} `json:"breeds"`
	Colors struct {
		Primary   *string `json:"primary"`
		Secondary *string `json:"secondary"`
		Tertiary  *string `json:"tertiary"`
	} `json:"colors"`
	Age        *string `json:"age"`
	Gender     *string `json:"gender"`
	Size       *string `json:"size"`
	Coat       *string `json:"coat"`
	Attributes struct {
		SpayedNeutered bool `json:"spayed_neutered"`
		HouseTrained   bool `json:"house_trained"`
		Declawed       bool `json:"declawed"`
		SpecialNeeds   bool `json:"special_needs"`
		ShotsCurrent   bool `json:"shots_current"`
	} `json:"attributes"`
	Environment struct {
		Children bool `json:"children"`
		Dogs     bool `json:"dogs"`
		Cats     bool `json:"cats"`
	} `json:"environment"`
	Tags                 []string `json:"tags"`
	Name                 string   `json:"name"`
	Description          *string  `json:"description"`
	OrganizationAnimalID *string  `json:"organization_animal_id"`
	Photos               []struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
		Full   string `json:"full"`
	} `json:"photos"`
	PrimaryPhotoCropped *struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
		Full   string `json:"full"`
	} `json:"primary_photo_cropped"`
	Videos []struct {
		Embed string `json:"embed"`
	} `json:"videos"`
	Status          *string `json:"status"`
	StatusChangedAt *string `json:"status_changed_at"`
	PublishedAt     *string `json:"published_at"`
	Distance        *string `json:"distance"`
	Contact         struct {
		Email   *string `json:"email"`
		Phone   *string `json:"phone"`
		Address struct {
			Address1 *string `json:"address1"`
			Address2 *string `json:"address2"`
			City     *string `json:"city"`
			State    *string `json:"state"`
			Postcode *string `json:"postcode"`
			Country  *string `json:"country"`
		} `json:"address"`
	} `json:"contact"`
}

type petfinderAnimalData struct {
	Animals []petfinderAnimal
}

type importAnimalResult struct {
	Name string
	Err  error
}

var (
	typeToID         = make(map[string]int64)
	specieToID       = make(map[string]int64)
	breedToID        = make(map[string]int64)
	specieIDToSchema = make(map[int64]map[string]any)

	microchipDescriptions = []string{"to find if lost", "to track if lost", "for tracking position"}
	microchipLocations    = []string{"in left ear", "in right ear", "neck", "left shoulder", "right shoulder", "upper chest"}

	animalGenders = []string{"m", "f"}
	animalSizes   = []string{"small", "medium", "large"}
	animalAges    = []string{"baby", "youg", "adult", "senior"}

	animalIDSeq int64
)

func (ic *ImportAnimalsCmd) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	kratosClient := kratos.NewClient(cfg.App.KratosPublicURL, cfg.App.KratosAdminURL)

	ketoClient, err := keto.NewClient(cfg.App.KetoReadURL, cfg.App.KetoWriteURL)
	if err != nil {
		return err
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()

	if err := cacheAnimalLookups(ctx, pool); err != nil {
		return fmt.Errorf("failed to cache animal lookups: %w", err)
	}

	userID, err := fetchUserID(ctx, kratosClient)
	if err != nil {
		return fmt.Errorf("failed to fetch user id: %w", err)
	}

	pages, err := filepath.Glob(filepath.Join(ic.Dir, "animals", "page_*.json"))
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
		jobs               = make(chan string)
		animalSetters      = make(chan models.AnimalSetter, 1_000)
		microchipSetters   = make(chan models.MicrochipSetter, 1_000)
		animalBreedSetters = make(chan models.AnimalBreedSetter, 1_000)
		animalTagSetters   = make(chan models.AnimalTagSetter, 1_000)
		animalPhotoSetters = make(chan models.AnimalPhotoSetter, 1_000)
		// animalVideoSetters = make(chan dbtype.AnimalVideoSetter, 1_000)

		animalResults          = make(chan copyResult)
		animalMicrochipResults = make(chan copyResult)
		animalBreedResults     = make(chan copyResult)
		animalTagResults       = make(chan copyResult)
		animalPhotoResults     = make(chan copyResult)
		// animalVideoResults     = make(chan copyResult)

		wg sync.WaitGroup
	)

	if err := disableConstraintsForAnimalImport(ctx, pool); err != nil {
		return fmt.Errorf("disableConstraints: %w", err)
	}

	for range ic.Workers {
		wg.Go(func() {
			for file := range jobs {
				_ = processAnimal(file, userID, cfg.App.BaseURL, animalSetters, microchipSetters, animalBreedSetters, animalTagSetters, animalPhotoSetters)
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
		close(animalSetters)
		close(microchipSetters)
		close(animalBreedSetters)
		close(animalTagSetters)
		close(animalPhotoSetters)
		// close(animalVideoSetters)
	}()

	var wg2 sync.WaitGroup
	wg2.Add(5)

	go func() {
		defer func() {
			close(animalResults)
			wg2.Done()
		}()

		c, err := copyAnimalData(ctx, pool, animalSetters)
		animalResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(animalMicrochipResults)
			wg2.Done()
		}()

		c, err := copyAnimalMicrochipData(ctx, pool, microchipSetters)
		animalMicrochipResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(animalBreedResults)
			wg2.Done()
		}()

		c, err := copyAnimalBreedData(ctx, pool, animalBreedSetters)
		animalBreedResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(animalTagResults)
			wg2.Done()
		}()

		c, err := copyAnimalTagData(ctx, pool, animalTagSetters)
		animalTagResults <- copyResult{count: c, err: err}
	}()

	go func() {
		defer func() {
			close(animalPhotoResults)
			wg2.Done()
		}()

		c, err := copyAnimalPhotoData(ctx, pool, animalPhotoSetters)
		animalPhotoResults <- copyResult{count: c, err: err}
	}()

	// go func() {
	// 	defer func() {
	// 		close(animalVideoResults)
	// 		wg2.Done()
	// 	}()
	// 	c, err := copyAnimalVideoData(ctx, dbpool, animalVideoSetters)
	// 	animalVideoResults <- copyResult{count: c, err: err}
	// }()

	var totalErrors []error

	for r := range animalResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("animal: %w", r.err))
		}
	}

	for r := range animalMicrochipResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("microchip: %w", r.err))
		}
	}

	for r := range animalBreedResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("animal_breed: %w", r.err))
		}
	}

	for r := range animalTagResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("animal_tag: %w", r.err))
		}
	}

	for r := range animalPhotoResults {
		if r.err != nil {
			totalErrors = append(totalErrors, fmt.Errorf("animal_photo: %w", r.err))
		}
	}
	// for r := range animalVideoResults {
	// 	if r.err != nil {
	// 		totalErrors = append(totalErrors, fmt.Errorf("animal_video: %w", r.err))
	// 	}
	// }

	if len(totalErrors) > 0 {
		fmt.Println("total errors", len(totalErrors))
		return errors.Join(totalErrors...)
	}

	wg2.Wait()

	if err := enableConstraintsForAnimalImport(ctx, pool); err != nil {
		return fmt.Errorf("enableConstraints: %w", err)
	}

	if err := setSequenceIDsAfterAnimalImport(ctx, pool); err != nil {
		return fmt.Errorf("setSequenceIDs: %w", err)
	}

	if err := postImportAssignOrganizationIDToAnimal(ctx, pool); err != nil {
		return fmt.Errorf("assignOrgIDs: %w", err)
	}

	if err := createAnimalRelationTuples(ctx, ketoClient); err != nil {
		return fmt.Errorf("create animal tuples: %w", err)
	}

	return nil
}

func disableConstraintsForAnimalImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "ALTER TABLE animal DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE microchip DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_breed DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_tag DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_photo DISABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_video DISABLE TRIGGER ALL")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	log.Println("disabled constraints before animals bulk import")

	return nil
}

func enableConstraintsForAnimalImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "ALTER TABLE animal ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE microchip ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_breed ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_tag ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_photo ENABLE TRIGGER ALL")
	_, _ = tx.Exec(ctx, "ALTER TABLE animal_video ENABLE TRIGGER ALL")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	log.Println("enabled constraints after animals bulk import")

	return nil
}

func setSequenceIDsAfterAnimalImport(ctx context.Context, dbpool *pgxpool.Pool) error {
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("seqIDs dbpool.BeginTx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback: %v", err)
		}
	}()

	_, _ = tx.Exec(ctx, "select setval('animal_id_seq', (SELECT MAX(id) FROM animal))")
	_, _ = tx.Exec(ctx, "select setval('microchip_id_seq', (SELECT MAX(id) FROM microchip))")
	_, _ = tx.Exec(ctx, "select setval('animal_tag_id_seq', (SELECT MAX(id) FROM animal_tag))")
	_, _ = tx.Exec(ctx, "select setval('animal_photo_id_seq', (SELECT MAX(id) FROM animal_photo))")
	// _, _ = tx.Exec(ctx, "select setval('animal_video_id_seq', (SELECT MAX(id) FROM animal_video))")

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("seqIDs tx.Commit: %w", err)
	}

	return nil
}

func processAnimal(
	filename string,
	userID uuid.UUID,
	baseURL string,
	animalSetters chan<- models.AnimalSetter,
	microchipSetters chan<- models.MicrochipSetter,
	animalBreedSetters chan<- models.AnimalBreedSetter,
	animalTagSetters chan<- models.AnimalTagSetter,
	animalPhotoSetters chan<- models.AnimalPhotoSetter,
	// animalVideoSetters chan<- dbtype.AnimalVideoSetter,
) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("faied to read file: %w", err)
	}

	var data petfinderAnimalData
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("faied to unmarshal petfinder animal: %w", err)
	}

	for _, animal := range data.Animals {
		acs := makeAnimalCreateSetter(animal, userID, baseURL)

		animalSetters <- acs.Animal

		if !acs.Microchip.IsUnset() {
			microchipSetters <- acs.Microchip.GetOrZero()
		}

		if !acs.Breeds.IsUnset() {
			for _, x := range acs.Breeds.GetOrZero() {
				animalBreedSetters <- x
			}
		}

		if !acs.Tags.IsUnset() {
			for _, x := range acs.Tags.GetOrZero() {
				animalTagSetters <- x
			}
		}

		if !acs.Photos.IsUnset() {
			for _, x := range acs.Photos.GetOrZero() {
				animalPhotoSetters <- x
			}
		}
		// if acs.Videos.IsSpecified() {
		// 	for _, x := range acs.Videos.GetOrZero() {
		// 		animalVideoSetters <- x
		// 	}
		// }
	}

	return nil
}

func copyAnimalData(ctx context.Context, dbpool *pgxpool.Pool, animalSetters <-chan models.AnimalSetter) (int64, error) {
	cols := []string{
		"user_id",
		"organization_id",
		"type_id",
		"specie_id",
		"name",
		"gender",
		"hermaphrodite",
		"age",
		"size",
		"image_object_kind",
		"image_object_ref_small",
		"image_object_ref_medium",
		"image_object_ref_large",
		"image_object_ref_full",
		"description",
		"distance",
		"properties",
		"status",
		"original_organization_id",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"animal"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-animalSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.UserID.GetOrZero(),
				valOrNil(v.OrganizationID),
				v.TypeID.GetOrZero(),
				v.SpecieID.GetOrZero(),
				v.Name.GetOrZero(),
				valOrNil(v.Gender),
				v.Hermaphrodite.GetOrZero(),
				v.Age.GetOrZero(),
				v.Size.GetOrZero(),
				v.ImageObjectKind.GetOrZero(),
				v.ImageObjectRefSmall.GetOrZero(),
				v.ImageObjectRefMedium.GetOrZero(),
				v.ImageObjectRefLarge.GetOrZero(),
				v.ImageObjectRefFull.GetOrZero(),
				valOrNil(v.Description),
				valOrNil(v.Distance),
				valOrNil(v.Properties),
				"adoptable",
				valOrNil(v.OriginalOrganizationID),
			}

			return row, nil
		}),
	)
}

func copyAnimalMicrochipData(ctx context.Context, dbpool *pgxpool.Pool, microchipSetters <-chan models.MicrochipSetter) (int64, error) {
	cols := []string{
		"animal_id",
		"number",
		"brand",
		"description",
		"location",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"microchip"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-microchipSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.AnimalID),
				v.Number.GetOrZero(),
				valOrNil(v.Brand),
				valOrNil(v.Description),
				valOrNil(v.Location),
			}

			return row, nil
		}),
	)
}

func copyAnimalBreedData(ctx context.Context, dbpool *pgxpool.Pool, animalBreedSetters <-chan models.AnimalBreedSetter) (int64, error) {
	cols := []string{"animal_id", "breed_id", "primary"}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"animal_breed"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-animalBreedSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				v.AnimalID.GetOrZero(),
				v.BreedID.GetOrZero(),
				v.Primary.GetOrZero(),
			}

			return row, nil
		}),
	)
}

func copyAnimalTagData(ctx context.Context, dbpool *pgxpool.Pool, tagSetters <-chan models.AnimalTagSetter) (int64, error) {
	cols := []string{"animal_id", "name"}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"animal_tag"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-tagSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.AnimalID),
				v.Name.GetOrZero(),
			}

			return row, nil
		}),
	)
}

func copyAnimalPhotoData(ctx context.Context, dbpool *pgxpool.Pool, photoSetters <-chan models.AnimalPhotoSetter) (int64, error) {
	cols := []string{
		"animal_id",
		"object_kind",
		"object_ref_small",
		"object_ref_medium",
		"object_ref_large",
		"object_ref_full",
	}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"animal_photo"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-photoSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.AnimalID),
				v.ObjectKind.GetOrZero(),
				valOrNil(v.ObjectRefSmall),
				valOrNil(v.ObjectRefMedium),
				valOrNil(v.ObjectRefLarge),
				valOrNil(v.ObjectRefFull),
			}

			return row, nil
		}),
	)
}

func copyAnimalVideoData(ctx context.Context, dbpool *pgxpool.Pool, videoSetters <-chan models.AnimalVideoSetter) (int64, error) {
	cols := []string{"animal_id", "object_kind", "object_ref"}

	return dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"animal_video"},
		cols,
		pgx.CopyFromFunc(func() ([]any, error) {
			v, ok := <-videoSetters
			if !ok {
				return nil, nil
			}

			row := []any{
				valOrNil(v.AnimalID),
				v.ObjectKind.GetOrZero(),
				v.ObjectRef.GetOrZero(),
			}

			return row, nil
		}),
	)
}

func postImportAssignOrganizationIDToAnimal(ctx context.Context, dbpool *pgxpool.Pool) error {
	q := `update animal
	set organization_id = organization.id
	from organization
	where animal.original_organization_id = organization.original_id`
	if _, err := dbpool.Exec(ctx, q); err != nil {
		return fmt.Errorf("assign organization ids to those original import ids: %w", err)
	}

	return nil
}

func fetchUserID(ctx context.Context, kratosClient *kratos.Client) (uuid.UUID, error) {
	req := kratosClient.Admin.IdentityAPI.ListIdentities(ctx)
	req = req.CredentialsIdentifier("fluffly@test.com")

	identities, identityResp, err := req.Execute()
	if err != nil {
		return uuid.Nil, fmt.Errorf("fetch users with credential identifiers: %w", err)
	}

	defer func() { _ = identityResp.Body.Close() }()

	if len(identities) != 1 {
		return uuid.Nil, fmt.Errorf("users len != 1")
	}

	id := identities[0].GetId()
	if id == "" {
		return uuid.Nil, fmt.Errorf("empty id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse uuid")
	}

	return uid, nil
}

func cacheAnimalLookups(ctx context.Context, dbpool *pgxpool.Pool) error {
	type typeRow struct {
		ID   int64
		Name string
	}

	type speciesRow struct {
		ID               int64
		Name             string
		PropertiesSchema map[string]any
	}

	type breedRow struct {
		ID   int64
		Name string
	}

	typeRows, err := dbpool.Query(ctx, "select id, name from animal_type")
	if err != nil {
		return fmt.Errorf("dbpool.Query animal_type: %w", err)
	}

	types, err := pgx.CollectRows(typeRows, pgx.RowToStructByName[typeRow])
	if err != nil {
		return fmt.Errorf("pgx.CollectRows animal_type: %w", err)
	}

	specieRows, err := dbpool.Query(ctx, "select id, name, properties_schema from animal_specie")
	if err != nil {
		return fmt.Errorf("dbpool.Query animal_specie: %w", err)
	}

	species, err := pgx.CollectRows(specieRows, pgx.RowToStructByName[speciesRow])
	if err != nil {
		return fmt.Errorf("pgx.CollectRows animal_specie: %w", err)
	}

	breedRows, err := dbpool.Query(ctx, "select id, name from breed")
	if err != nil {
		return fmt.Errorf("dbpool.Query breed: %w", err)
	}

	breeds, err := pgx.CollectRows(breedRows, pgx.RowToStructByName[breedRow])
	if err != nil {
		return fmt.Errorf("pgx.CollectRows breed: %w", err)
	}

	for _, x := range types {
		typeToID[x.Name] = x.ID
	}

	for _, x := range species {
		specieToID[x.Name] = x.ID
		specieIDToSchema[x.ID] = x.PropertiesSchema
	}

	for _, x := range breeds {
		breedToID[x.Name] = x.ID
	}

	return nil
}

func makeAnimalCreateSetter(animal petfinderAnimal, userID uuid.UUID, baseURL string) dbtype.AnimalCreateSetter {
	var animalCreateSetter dbtype.AnimalCreateSetter

	fallbackImgMax := 5

	fallbackImgName := strings.ToLower(animal.Type)
	switch animal.Type {
	case "Small & Furry":
		fallbackImgName = "rabbit"
	case "Scales, Fins & Other":
		fallbackImgName = "aquatic"
	case "Barnyard":
	case "Horse":
		fallbackImgMax = 3
		fallbackImgName = "farm"
	}

	fallbackImg := fmt.Sprintf("%s/public/images/placeholder/%s-%d.svg", baseURL, fallbackImgName, fallbackImgMax)

	animalTypeID := typeToID[animal.Type]
	animalSpecieID := specieToID[animal.Species]

	animalID := atomic.AddInt64(&animalIDSeq, 1)

	animalSetter := models.AnimalSetter{
		ID:                     omit.From(animalID),
		UserID:                 omit.From(userID),
		TypeID:                 omit.From(animalTypeID),
		SpecieID:               omit.From(animalSpecieID),
		Name:                   omit.From(animal.Name),
		Description:            omitnull.FromPtr(animal.Description),
		OriginalOrganizationID: omitnull.FromPtr(animal.OrganizationID),
	}
	if animal.Gender != nil {
		gender := enums.GenderM
		if *animal.Gender == "Female" {
			gender = enums.GenderF
		}

		animalSetter.Gender = omitnull.From(gender)
	}

	var hermaphrodite bool
	if random.PercentChance(1) {
		hermaphrodite = true
	}

	animalSetter.Hermaphrodite = omit.From(hermaphrodite)
	if animal.Age != nil {
		animalSetter.Age = omit.From(strings.ToLower(*animal.Age))
	} else {
		animalSetter.Age = omit.From(pickRandomItem(animalAges))
	}

	if animal.Size != nil {
		size := strings.ToLower(*animal.Size)
		if size == "extra large" {
			size = "large"
		}

		animalSetter.Size = omit.From(size)
	} else {
		animalSetter.Size = omit.From(pickRandomItem(animalSizes))
	}

	if animal.PrimaryPhotoCropped != nil {
		animalSetter.ImageObjectKind = omit.From(media.StorageKindExternal)
		if animal.PrimaryPhotoCropped.Small != "" {
			animalSetter.ImageObjectRefSmall = omit.From(animal.PrimaryPhotoCropped.Small)
		}

		if animal.PrimaryPhotoCropped.Medium != "" {
			animalSetter.ImageObjectRefMedium = omit.From(animal.PrimaryPhotoCropped.Medium)
		}

		if animal.PrimaryPhotoCropped.Large != "" {
			animalSetter.ImageObjectRefLarge = omit.From(animal.PrimaryPhotoCropped.Large)
		}

		if animal.PrimaryPhotoCropped.Full != "" {
			animalSetter.ImageObjectRefFull = omit.From(animal.PrimaryPhotoCropped.Full)
		}
	} else {
		if len(animal.Photos) > 0 {
			var small, medium, large, full string

			for _, photo := range animal.Photos {
				if photo.Small != "" {
					small = photo.Small
				}

				if photo.Medium != "" {
					medium = photo.Medium
				}

				if photo.Large != "" {
					large = photo.Large
				}

				if photo.Full != "" {
					full = photo.Full
				}
			}

			animalSetter.ImageObjectKind = omit.From(media.StorageKindExternal)
			animalSetter.ImageObjectRefSmall = omit.From(small)
			animalSetter.ImageObjectRefMedium = omit.From(medium)
			animalSetter.ImageObjectRefLarge = omit.From(large)
			animalSetter.ImageObjectRefFull = omit.From(full)
		} else {
			animalSetter.ImageObjectKind = omit.From(media.StorageKindExternal)
			animalSetter.ImageObjectRefSmall = omit.From(fallbackImg)
			animalSetter.ImageObjectRefMedium = omit.From(fallbackImg)
			animalSetter.ImageObjectRefLarge = omit.From(fallbackImg)
			animalSetter.ImageObjectRefFull = omit.From(fallbackImg)
		}
	}

	properties := dbcustom.Properties{}
	schema := specieIDToSchema[animalSpecieID]
	schemaProperties := schema["properties"].(map[string]any)

	for prop, value := range schemaProperties {
		val := value.(map[string]any)
		kind := val["type"].(string)

		switch kind {
		case "boolean":
			if random.PercentChance(70) {
				properties[prop] = pickRandomItem([]bool{true, false})
			}
		case "string":
			if enum, ok := val["enum"].([]any); ok {
				// @TODO: support multi colors per animal later
				if prop == "color" {
					randColor := pickRandomItem(enum).(string)
					properties[prop] = cmp.Or(animal.Colors.Primary, animal.Colors.Secondary, animal.Colors.Tertiary, &randColor)
				} else {
					properties[prop] = pickRandomItem(enum)
				}
			} else if random.PercentChance(30) {
				properties[prop] = faker.Word()
			}
		}
	}

	animalSetter.Properties = omitnull.From(types.NewJSON(properties))
	animalCreateSetter.Animal = animalSetter

	breeds := make([]string, 0)
	if animal.Breeds.Primary != nil {
		breeds = append(breeds, *animal.Breeds.Primary)
	}

	if animal.Breeds.Secondary != nil {
		breeds = append(breeds, *animal.Breeds.Secondary)
	}

	uniqueBreedNames := uniqueStrings(breeds)
	if len(uniqueBreedNames) > 0 {
		animalBreedsSetters := make([]models.AnimalBreedSetter, 0)

		for i, breed := range uniqueBreedNames {
			breedID, ok := breedToID[breed]
			if ok {
				breedSetter := models.AnimalBreedSetter{AnimalID: omit.From(animalID), BreedID: omit.From(breedID)}
				if i == 0 {
					breedSetter.Primary = omit.From(true)
				}

				animalBreedsSetters = append(animalBreedsSetters, breedSetter)
			}
		}

		animalCreateSetter.Breeds = omitnull.From(animalBreedsSetters)
	}

	if random.PercentChance(30) {
		microchipSetter := models.MicrochipSetter{
			AnimalID: omitnull.From(animalID),
			Number:   omit.From(random.AlphaNumeric(6)),
			Brand:    omitnull.From(faker.Word()),
		}
		if random.PercentChance(50) {
			microchipSetter.Description = omitnull.From(pickRandomItem(microchipDescriptions))
			if animal.Type == "Cat" || animal.Type == "Dog" || animal.Type == "Rabbit" || animal.Type == "Small & Furry" {
				microchipSetter.Location = omitnull.From(pickRandomItem(microchipLocations))
			}
		}

		animalCreateSetter.Microchip = omitnull.From(microchipSetter)
	}

	if len(animal.Photos) > 0 {
		animalPhotoSetters := make([]models.AnimalPhotoSetter, len(animal.Photos))
		for i, photo := range animal.Photos {
			photoSetter := models.AnimalPhotoSetter{
				AnimalID:        omitnull.From(animalID),
				ObjectKind:      omit.From(string(media.StorageKindExternal)),
				ObjectRefSmall:  omitnull.From(photo.Small),
				ObjectRefMedium: omitnull.From(photo.Medium),
				ObjectRefLarge:  omitnull.From(photo.Large),
				ObjectRefFull:   omitnull.From(photo.Full),
			}
			animalPhotoSetters[i] = photoSetter
		}

		animalCreateSetter.Photos = omitnull.From(animalPhotoSetters)
	}

	if len(animal.Tags) > 0 {
		animalTagSetters := make([]models.AnimalTagSetter, 0)
		for _, tag := range uniqueStrings(animal.Tags) {
			animalTagSetters = append(animalTagSetters, models.AnimalTagSetter{
				AnimalID: omitnull.From(animalID),
				Name:     omit.From(tag),
			})
		}

		animalCreateSetter.Tags = omitnull.From(animalTagSetters)
	}

	// can't bother with their moronic embed video with html tags and stuff
	// they can't store a normal url like normal people but store html tag like retards

	return animalCreateSetter
}

func pickRandomItem[T any](xs []T) T {
	return xs[rand.IntN(len(xs))]
}

func uniqueStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	return out
}

func createAnimalRelationTuples(ctx context.Context, c *keto.Client) error {
	const batchSize = 5_000

	for start := int64(0); start < animalIDSeq; start += batchSize {
		end := min(start+batchSize, animalIDSeq)

		relationTupleDeltas := make([]*rts.RelationTupleDelta, 0, end-start)

		for i := start; i < end; i++ {
			relationTupleDeltas = append(relationTupleDeltas, &rts.RelationTupleDelta{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Animal",
					Object:    shared.AuthzAnimalID(i + 1),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Animals", "animals", ""),
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

func valOrNil[T any](v omitnull.Val[T]) *T {
	var out *T
	if !v.IsUnset() && !v.IsNull() {
		out = v.MustPtr()
	}

	return out
}
