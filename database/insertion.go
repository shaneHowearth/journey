// Package database -
package database

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

const stmtInsertPost = "INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, status, metaDescription, image, author_id, createdAt, createdBy, updated_at, updated_by, published_at, published_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertUser = "INSERT INTO users (id, uuid, name, slug, password, email, image, cover, createdAt, createdBy, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertRoleUser = "INSERT INTO roles_users (id, roleID, userID) VALUES (?, ?, ?)"
const stmtInsertTag = "INSERT INTO tags (id, uuid, name, slug, createdAt, createdBy, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertPostTag = "INSERT INTO posts_tags (id, postID, tagID) VALUES (?, ?, ?)"
const stmtInsertSetting = "INSERT INTO settings (id, uuid, key, value, type, createdAt, createdBy, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

// InsertPost -
func InsertPost(title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription []byte, image []byte, createdAt time.Time, createdBy int64) (int64, error) {

	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	var result sql.Result
	if published {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.NewV4().String(), title, slug, markdown, html, featured, isPage, status, metaDescription, image, createdBy, createdAt, createdBy, createdAt, createdBy, createdAt, createdBy)
	} else {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.NewV4().String(), title, slug, markdown, html, featured, isPage, status, metaDescription, image, createdBy, createdAt, createdBy, createdAt, createdBy, nil, nil)
	}
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	return postID, writeDB.Commit()
}

// InsertUser -
func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, createdAt time.Time, createdBy int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertUser, nil, uuid.NewV4().String(), name, slug, password, email, image, cover, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	return userID, writeDB.Commit()
}

// InsertRoleUser -
func InsertRoleUser(roleID int, userID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtInsertRoleUser, nil, roleID, userID)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}

// InsertTag -
func InsertTag(name []byte, slug string, createdAt time.Time, createdBy int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertTag, nil, uuid.NewV4().String(), name, slug, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	tagID, err := result.LastInsertId()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return 0, fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return 0, err
	}
	return tagID, writeDB.Commit()
}

// InsertPostTag -
func InsertPostTag(postID int64, tagID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtInsertPostTag, nil, postID, tagID)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}

func insertSettingString(key string, value string, settingType string, createdAt time.Time, createdBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.NewV4().String(), key, value, settingType, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}

func insertSettingInt64(key string, value int64, settingType string, createdAt time.Time, createdBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.NewV4().String(), key, value, settingType, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeErr := writeDB.Rollback()
		if writeErr != nil {
			return fmt.Errorf("rollback generated error %v whilst handling %w", writeErr, err)
		}
		return err
	}
	return writeDB.Commit()
}
