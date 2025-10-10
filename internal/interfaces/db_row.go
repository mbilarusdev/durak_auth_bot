package interfaces

type DBRow interface {
	Scan(dest ...any) error
}
