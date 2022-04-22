package entity

type (
	UploadFileModel struct {
		Id     int     `db:"id"`
		Format *string `db:"formats"`
		Url    string  `db:"url"`
	}

	UploadFileModels []UploadFileModel

	ModifyUploadFileModel struct {
		Id     int    `db:"id"`
		Format string `db:"formats"`
	}

	ModifyUploadFileModels []ModifyUploadFileModel

	Image struct {
		Ext    string  `json:"ext"`
		Url    string  `json:"url"`
		Hash   string  `json:"hash"`
		Mime   string  `json:"mime"`
		Name   string  `json:"name"`
		Path   *string `json:"path"`
		Size   float64 `json:"size"`
		Width  int     `json:"width"`
		Height int     `json:"height"`
	}

	UploadFileFormat struct {
		Large     *Image `json:"large,omitempty"`
		Small     *Image `json:"small,omitempty"`
		Medium    *Image `json:"medium,omitempty"`
		Thumbnail *Image `json:"thumbnail,omitempty"`
	}

	ImageUrl struct {
		Url string `json:"url,omitempty"`
	}

	ModifyUploadFileModelUrl struct {
		Id  int    `db:"id"`
		Url string `db:"url"`
	}

	ModifyUploadFileModelUrls []ModifyUploadFileModelUrl
)
