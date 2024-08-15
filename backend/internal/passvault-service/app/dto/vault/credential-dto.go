package vault

type Credential struct {
	Name       string `json:"name" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Credential string `json:"credential"`
	Url        string `json:"url"`
	Note       string `json:"note"`
}
