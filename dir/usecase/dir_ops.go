package usecase

import (
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/dtos"
	"github.com/sicozz/papyrus/domain/mapper"
	"github.com/sicozz/papyrus/utils/constants"
)

func fillDetailsBFS(nodeUuid string, dirs []domain.Dir, nodePath string, nodeNchild, nodeDepth int) (res []dtos.DirGetDto, err error) {
	root, err := findDirByUuid(nodeUuid, dirs)
	if err != nil {
		return nil, err
	}
	rootDto := mapper.MapDirToDirGetDto(root, nodePath, nodeNchild, nodeDepth)
	rootDto.ParentDir = ""

	detailed := []dtos.DirGetDto{rootDto}
	q := getChildren(root.Uuid, dirs)

	for len(q) > 0 {
		d := q[0]
		q = q[1:]

		pDto, err := findDtoByUuid(d.ParentDir, detailed)
		if err != nil {
			return nil, err
		}
		nChild := getNchild(d.Uuid, dirs)
		path := strings.Join(
			[]string{pDto.Path, d.Name},
			string(os.PathSeparator),
		)
		depth := pDto.Depth + 1

		dDto := mapper.MapDirToDirGetDto(d, path, nChild, depth)
		detailed = append(detailed, dDto)

		q = append(q, getChildren(d.Uuid, dirs)...)
	}

	for i := range detailed {
		if detailed[i].Uuid == constants.RootDirUuid {
			detailed[i].ParentDir = constants.RootDirUuid
		}
	}
	res = detailed

	return
}

func getChildren(pUuid string, dirs []domain.Dir) (children []domain.Dir) {
	children = []domain.Dir{}
	for _, d := range dirs {
		if d.ParentDir == pUuid && d.ParentDir != d.Uuid {
			children = append(children, d)
		}
	}

	return
}

func findDirByUuid(uuid string, dirs []domain.Dir) (dir domain.Dir, err error) {
	for _, d := range dirs {
		if d.Uuid == uuid {
			return d, nil
		}
	}

	return domain.Dir{}, errors.New("Failed to find dir with uuid: " + uuid)
}

func findDtoByUuid(uuid string, dirs []dtos.DirGetDto) (dto dtos.DirGetDto, err error) {
	for _, d := range dirs {
		if d.Uuid == uuid {
			return d, nil
		}
	}

	return dtos.DirGetDto{}, errors.New("Failed to find dir with uuid: " + uuid)
}

func getNchild(uuid string, dirs []domain.Dir) (nChild int) {
	nChild = 0
	for _, d := range dirs {
		if d.ParentDir == uuid {
			nChild += 1
		}
	}

	return
}

func getBranchFromTree(nodeUuid string, dirs []domain.Dir) (branch []domain.Dir, err error) {
	node, err := findDirByUuid(nodeUuid, dirs)
	if err != nil {
		return nil, err
	}

	branch = []domain.Dir{}
	q := []domain.Dir{node}

	for len(q) > 0 {
		d := q[0]
		q = q[1:]

		branch = append(branch, d)

		q = append(q, getChildren(d.Uuid, dirs)...)
	}

	return
}

func flushTreeUuids(nodeUuid string, dirs []domain.Dir) (res []domain.Dir, err error) {
	node, err := findDirByUuid(nodeUuid, dirs)
	if err != nil {
		return nil, err
	}

	nBranch := []domain.Dir{}
	q := []domain.Dir{node}

	for len(q) > 0 {
		d := q[0]
		q = q[1:]

		nUuid := uuid.New().String()

		for i, dd := range dirs {
			if dd.ParentDir == d.Uuid {
				dirs[i].ParentDir = nUuid
			}
		}

		d.Uuid = nUuid

		nBranch = append(nBranch, d)

		q = append(q, getChildren(d.Uuid, dirs)...)
	}

	res = nBranch

	return
}
