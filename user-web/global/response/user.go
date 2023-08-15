package response

type UserResponse struct {
	Id       int32  `json:"id,omitempty"`
	NickName string `json:"name,omitempty"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
}
