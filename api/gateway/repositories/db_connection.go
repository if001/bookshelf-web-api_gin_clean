package repositories

type DBConnection interface {
	Bind(bind interface{}) DBConnection
	Where(filter interface{}) DBConnection
	Like(key, filter string) DBConnection
	OrLike(key, filter string) DBConnection
	Paginate(page, perPage uint64) DBConnection
	OrFilter(filter interface{}) DBConnection
	Create(data interface{}) DBConnection
	Delete(data interface{}) DBConnection
	Update(data interface{}) DBConnection
	SortDesc(key string) DBConnection
	SortAsc(key string) DBConnection
	GroupBy(key string) DBConnection
	Limit(num int) DBConnection
	Count(count *int64) DBConnection
	Table(table interface{}) DBConnection
	TX() DBConnection
	TxRollback() error
	TxExec() error
	CountedAuthorQuery(bind interface{}) error
	CountedPublisherQuery(bind interface{}) error
	HasError() error
	SelectBookWith() DBConnection
	SelectBookWithAuthorName() DBConnection
	SelectBookWithPublisherName() DBConnection
	GroupByDate(key, format string) DBConnection
	SearchBook(value string) DBConnection
}
