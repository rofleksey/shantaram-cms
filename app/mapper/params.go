package mapper

import (
	"shantaram/app/api"
	"shantaram/pkg/database"
)

func MapParams(p database.Param) api.Params {
	return api.Params{
		HeaderText:     p.HeaderText,
		HeaderDeadline: p.HeaderDeadline,
	}
}
