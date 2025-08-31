package postgres

type IntegrityViolationError interface {
	error
	IsIntegrityViolationError()
}

type UniqueViolationError interface {
	IntegrityViolationError
	IsUniqueViolationError()
}

type RestrictViolationError interface {
	IntegrityViolationError
	IsRestrictViolationError()
}

type ForeignKeyViolationError interface {
	IntegrityViolationError
	IsForeignKeyViolationError()
}

type NotNullViolationError interface {
	IntegrityViolationError
	IsNotNullViolationError()
}

type CheckViolationError interface {
	IntegrityViolationError
	IsCheckViolationError()
}

type ExclusionViolationError interface {
	IntegrityViolationError
	IsExclusionViolationError()
}

type errIntegrityViolation struct{}

func (euv errIntegrityViolation) Error() string {
	return "integrity violation"
}
func (euv errIntegrityViolation) IsIntegrityViolationError() {}

type errUniqueViolation struct {
	Name string
}

func (euv errUniqueViolation) Error() string {
	return "unique violation " + euv.Name
}
func (euv errUniqueViolation) IsIntegrityViolationError() {}
func (euv errUniqueViolation) IsUniqueViolationError()    {}

type errRestrictViolation struct {
	Name string
}

func (erv errRestrictViolation) Error() string {
	return "restrict violation " + erv.Name
}
func (erv errRestrictViolation) IsIntegrityViolationError() {}
func (erv errRestrictViolation) IsRestrictViolationError()  {}

type errForeignKeyViolation struct {
	Name string
}

func (efv errForeignKeyViolation) Error() string {
	return "foreign key violation " + efv.Name
}
func (efv errForeignKeyViolation) IsIntegrityViolationError()  {}
func (efv errForeignKeyViolation) IsForeignKeyViolationError() {}

type errNotNullViolation struct {
	Name string
}

func (env errNotNullViolation) Error() string {
	return "not null violation " + env.Name
}
func (env errNotNullViolation) IsIntegrityViolationError()  {}
func (env errNotNullViolation) IsForeignKeyViolationError() {}

type errCheckViolation struct {
	Name string
}

func (ecv errCheckViolation) Error() string {
	return "check violation " + ecv.Name
}
func (ecv errCheckViolation) IsIntegrityViolationError() {}
func (ecv errCheckViolation) IsCheckViolationError()     {}

type errExclusionViolation struct {
	Name string
}

func (eev errExclusionViolation) Error() string {
	return "exclusion violation " + eev.Name
}
func (eev errExclusionViolation) IsIntegrityViolationError() {}
func (eev errExclusionViolation) IsExclusionViolationError() {}
