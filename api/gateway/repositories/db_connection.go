package repositories

type DBConnection interface {
	Bind(bind interface{}) DBConnection
	Select(filter interface{}) DBConnection
	Paginate(page, perPage uint64) DBConnection
	OrFilter(filter interface{}) DBConnection
	Create(data interface{}) DBConnection
	Delete(data interface{}) DBConnection
	Update(data interface{}) DBConnection
	SortDesc(key string) DBConnection
	SortAsc(key string) DBConnection
	Count(count *int64) DBConnection
	Table(table interface{}) DBConnection
	HasError() error
}
