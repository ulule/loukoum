package loukoum

// Select start a SelectBuilder using given columns.
func Select(columns ...interface{}) SelectBuilder {
	return NewSelectBuilder().Columns(columns)
}

// SelectDistinct start a SelectBuilder using given columns and "DISTINCT" option.
func SelectDistinct(columns ...interface{}) SelectBuilder {
	return Select(columns...).Distinct()
}
