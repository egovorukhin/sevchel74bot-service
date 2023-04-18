package telegram

type Me struct {
	Id                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

// GetMe Инфа обо мне
func GetMe() (*Me, error) {
	me := &Me{}
	err := ExecuteGet("getMe", me)
	if err != nil {
		return nil, err
	}
	return me, nil
}
