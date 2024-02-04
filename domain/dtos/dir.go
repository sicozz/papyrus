package dtos

const (
	DirGetDtoType    = "carpeta"
	DirGetDtoVisible = false
	DirGetDtoOpen    = false
)

// DirGetDto represents the information served in the GetAll endpoint of DirHandler
type DirGetDto struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name" validate:"required,excludesall=!@#?"`
	ParentDir string `json:"parent_dir" validate:"required,uuid"`
	Path      string `json:"path"`
	Nchild    int    `json:"nchild"`
	Depth     int    `json:"depth"`
	Type      string `json:"type"`
	Visible   bool   `json:"visible"`
	Open      bool   `json:"open"`
	State     string `json:"state"`

	// file and/or plan attributes
	RespUser    string `json:"responsible_user"`
	CreatorUser string `json:"creator_user"`
	Subtype     string `json:"subtype"`
	Datecreate  string `json:"date_create"`
	DateClose   string `json:"date_close"`
	Term        int    `json:"term"`
}

type DirSizeGetDto struct {
	Size      string `json:"size"`
	FileCount int    `json:"file_count"`
}

type DirStoreDto struct {
	Name      string `json:"name" validate:"required,excludesall=!@#?"`
	ParentDir string `json:"parent_dir" validate:"required,uuid"`
}

type DirUpdateDto struct {
	Name string `json:"name" validate:"required,excludesall=!@#?"`
}

type DirDuplicateDto struct {
	Uuid      string `json:"uuid" validate:"uuid"`
	ParentDir string `json:"parent_dir" validate:"uuid"`
	Name      string `json:"name" validate:"excludesall=!@#?"`
}

type DirMoveDto struct {
	ParentDir string `json:"parent_dir" validate:"omitempty,uuid"`
}

type DocsNotDirGetDto struct {
	Uuid         string `json:"uuid"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	FsPath       string `json:"fs_path"`
	DateCreation string `json:"date_create"`
	DateInput    string `json:"date_input"`
	Type         string `json:"type"`
	State        string `json:"state"`
	ParentDir    string `json:"parent_dir"`
	RespUser     string `json:"responsible_user"`
	AppUser1     string `json:"approval_user1"`
	AppUser2     string `json:"approval_user2"`
	AppUser3     string `json:"approval_user3"`

	// TODO: make this a checks of boolean type
	Chk1    string `json:"user_check1"`
	Chk2    string `json:"user_check2"`
	Chk3    string `json:"user_check3"`
	Version string `json:"version"`
	Term    int    `json:"term"`
	Subtype string `json:"subtype"`

	Procedure   string `json:"procedure"`
	CreatorUser string `json:"creator_user"`
	RecvUser    string `json:"receiver_user"`
	Chk         bool   `json:"check"`

	Path  string `json:"path"`
	Depth int    `json:"depth"`
}
