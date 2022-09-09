package helper

import "dysn/auth/internal/model/consts"

var LangList = []interface{}{
	consts.RuLang,
	consts.EnLang,
}

func GetDefaultLang() string {
	return consts.RuLang
}
