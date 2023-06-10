// Cart table definition
package schema

// Describe cart table in postgres db
type Cart struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
}
