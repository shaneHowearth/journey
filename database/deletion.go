// Package database -
package database

import "fmt"

const stmtDeletePostTagsByPostID = "DELETE FROM posts_tags WHERE post_id = ?"
const stmtDeletePostByID = "DELETE FROM posts WHERE id = ?"

// DeletePostTagsForPostID -
func DeletePostTagsForPostID(postID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst dealing with %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtDeletePostTagsByPostID, postID)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst dealing with %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}

// DeletePostByID -
func DeletePostByID(ID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst dealing with %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtDeletePostByID, ID)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst dealing with %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}
