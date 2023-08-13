package domain

import "context"

// Dir represents the Directory data strict
type Dir struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name" validate:"required,ascii"`
	ParentDir string `json:"parent_dir" validate:"required,ascii"`
	Path      string `json:"path"`
	Nchild    int    `json:"nchild"`
}

// DirUsecase represents the dir's usecases
type DirUsecase interface {
	/* NOTE: On directory "in{active}" state
	* We are leaving it for later, when we recv some feedback, to decide if it
	* is necessary or we can implement it later
	 */
	// TODO: Add root dir migration
	GetAll(c context.Context) ([]Dir, RequestErr)
	// GetByUuid(c context.Context, uuid string) (Dir, RequestErr)
	// Store(c context.Context, d *Dir) RequestErr
	// Update(c context.Context) RequestErr
	// Delete(c context.Context) RequestErr
	// Move(c context.Context) RequestErr
}

// DirRepository represents the dir's repository contract
type DirRepository interface {
	GetAll(ctx context.Context) ([]Dir, error)
	// GetByUuid(ctx context.Context, uuid string) (Dir, error)
	// ExistsByUuid(ctx context.Context, uuid string) bool
	// ExistsByName(ctx context.Context, name string) bool
	// Store(ctx context.Context, d *Dir) error
	// Delete(ctx context.Context, uuid string) error
	// ChgName(ctx context.Context, uuid string, nName string) error
	// ChgParentDir(ctx context.Context, uuid string, nPUuid string) error
}
