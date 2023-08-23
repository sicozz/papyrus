package domain

import (
	"context"
	"errors"
	"os"
	"strings"

	gUuid "github.com/google/uuid"
	"github.com/sicozz/papyrus/domain/dtos"
)

// Dir represents the Directory data strict
type Dir struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name" validate:"required,ascii"`
	ParentDir string `json:"parent_dir" validate:"required,ascii,uuid"`
	Path      string `json:"path"`
	Nchild    int    `json:"nchild"`
	Depth     int    `json:"depth"`
}

// DirUsecase represents the dir's usecases
type DirUsecase interface {
	/* NOTE: On directory "in{active}" state
	* We are leaving it for later, when we recv some feedback, to decide if it
	* is necessary or we can implement it later
	 */
	GetAll(c context.Context) ([]dtos.DirGetDto, RequestErr)
	GetByUuid(c context.Context, uuid string) (dtos.DirGetDto, RequestErr)
	Store(c context.Context, d dtos.DirStoreDto) (dtos.DirGetDto, RequestErr)
	Update(c context.Context, uuid string, p dtos.DirUpdateDto) RequestErr
	Delete(c context.Context, uuid string) RequestErr
	Move(c context.Context, uuid string, nPUuid string) RequestErr
	Duplicate(c context.Context, uuid string, nName string, destUuid string) (dtos.DirGetDto, RequestErr)
}

// DirRepository represents the dir's repository contract
type DirRepository interface {
	GetAll(ctx context.Context) ([]Dir, error)
	GetByUuid(ctx context.Context, uuid string) (Dir, error)
	ExistByUuid(ctx context.Context, uuid string) bool
	// ExistsByName(ctx context.Context, name string) bool
	Store(ctx context.Context, d *Dir) error
	Delete(ctx context.Context, uuid string) error
	ChgName(ctx context.Context, uuid string, nName string) error

	// Refactor this group of functions
	ChgParentDir(ctx context.Context, uuid string, nPUuid string) error
	IncNchild(ctx context.Context, uuid string, nNchild int) error
	DecNchild(ctx context.Context, uuid string, nNchild int) error
	ChgPath(ctx context.Context, uuid string, nPath string) error
	ChgDepth(ctx context.Context, uuid string, nDepth int) error
	Insert(ctx context.Context, dir Dir) error
}

func FindDirByUuid(uuid string, dirs []*Dir) (*Dir, error) {
	for _, d := range dirs {
		if d.Uuid == uuid {
			return d, nil
		}
	}

	return nil, errors.New("Failed to find dir with uuid: " + uuid)
}

func GetByParent(pUuid string, dirs []*Dir) (children []*Dir) {
	children = []*Dir{}
	for _, d := range dirs {
		if d.ParentDir == pUuid {
			children = append(children, d)
		}
	}

	return
}

func RefreshChildren(uuid string, dirs []*Dir) (err error) {
	src, err := FindDirByUuid(uuid, dirs)
	if err != nil {
		return err
	}

	q := []*Dir{src}

	for len(q) > 0 {
		currDir := q[0]
		q = q[1:]

		currPDir, err := FindDirByUuid(currDir.ParentDir, dirs)
		if err != nil {
			return err
		}

		currDir.Path = strings.Join(
			[]string{currPDir.Path, currDir.Name},
			string(os.PathSeparator),
		)
		currDir.Depth = currPDir.Depth + 1

		q = append(q, GetByParent(currDir.Uuid, dirs)...)
	}

	return nil
}

func duplicateTopDown(srcUuid string, nName string, dstUuid string, dirs []Dir) (head *Dir, neoDirs []*Dir, err error) {
	dirSet := []*Dir{}
	for _, d := range dirs {
		tD := d
		dirSet = append(dirSet, &tD)
	}

	srcDir, err := FindDirByUuid(srcUuid, dirSet)
	srcDir.Name = nName
	if err != nil {
		return nil, nil, err
	}
	dstDir, err := FindDirByUuid(dstUuid, dirSet)
	if err != nil {
		return nil, nil, err
	}

	srcDir.ParentDir = dstUuid
	dstDir.Nchild += 1

	q := []*Dir{srcDir}

	neoDirs = []*Dir{}
	for len(q) > 0 {
		currDir := q[0]
		q = q[1:]

		currPDir, err := FindDirByUuid(currDir.ParentDir, dirSet)
		if err != nil {
			return nil, nil, err
		}

		currDir.Path = strings.Join(
			[]string{currPDir.Path, currDir.Name},
			string(os.PathSeparator),
		)
		currDir.Depth = currPDir.Depth + 1
		currChildren := GetByParent(currDir.Uuid, dirSet)
		currDir.Uuid = gUuid.New().String()
		for _, ch := range currChildren {
			ch.ParentDir = currDir.Uuid
		}

		q = append(q, currChildren...)
		neoDirs = append(neoDirs, currDir)
	}

	head = neoDirs[0]

	return
}

// TODO: This function works as a wrapper but is no necessary. Delete it later
func Duplicate(srcUuid string, nName string, dstUuid string, dirs []Dir) (neoDir *Dir, dupDirs []*Dir, err error) {
	neoDir, neoDirs, err := duplicateTopDown(srcUuid, nName, dstUuid, dirs)
	if err != nil {
		return nil, nil, err
	}

	// WARN: THIS GOES INTO REPOSITORY
	// dstDir.Nchild += 1

	return neoDir, neoDirs, nil
}
