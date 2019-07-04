package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type Filter struct{}

func NewFilter() map[string]interface{} {
	return map[string]interface{}{}
}

func ByAccountId(filter map[string]interface{}, id string) {
	filter["account_id"] = id
}
func ById(filter map[string]interface{}, id uint64) {
	filter["id"] = id
}
func ByBookId(filter map[string]interface{}, id uint64) {
	filter["book_id"] = id
}
func ByStatus(filter map[string]interface{}, status domain.ReadState) {
	filter["read_state"] = status
}
