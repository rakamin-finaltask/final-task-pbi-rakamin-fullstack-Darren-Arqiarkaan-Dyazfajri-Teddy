package dto

type PhotosRequest struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type PhotosResponse struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   string `json:"userId"`
}
