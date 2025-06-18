package repositories

import (
	"context"
	"fmt"
	"github.com/hotrungnhan/surl/utils/injects"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type BaseRepository struct {
	db *injects.DB
}

func newBaseRepository(db *injects.DB) BaseRepository {
	return BaseRepository{
		db: db,
	}
}

func (r *BaseRepository) BuildQuery(db *gorm.DB, filter interface{}) (*gorm.DB, error) {
	val := reflect.ValueOf(filter)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("filter must be struct or pointer to struct")
	}

	typ := val.Type()
	query := db

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)

		// Skip if pointer is nil (ignore this filter)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}

		if (fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array) && fieldVal.IsNil() {
			continue
		}

		// Dereference pointer to get the actual value
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		fieldType := typ.Field(i)
		filterTag := fieldType.Tag.Get("filter")
		if filterTag == "" {
			continue
		}

		parts := strings.Split(filterTag, ";")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid filter tag format on field %s", fieldType.Name)
		}

		operator := strings.ToLower(parts[0])
		column := parts[1]

		table := lo.TernaryF(len(parts) > 2, func() string { return parts[2] }, func() string { return "" })
		method := lo.TernaryF(len(parts) > 3, func() string { return parts[3] }, func() string { return "" })

		full_column_name := lo.Ternary(table == "", fmt.Sprintf("\"%s\"", column), fmt.Sprintf("\"%s\".\"%s\"", table, column))

		query_indentify := lo.TernaryF(method != "", func() string {
			return fmt.Sprintf("%s(%s)", full_column_name, method)
		}, func() string { return full_column_name })
		switch operator {
		case "in":
			// Expect fieldVal to be a slice or string with comma-separated values
			query = query.Where(fmt.Sprintf("%s IN ?", query_indentify), fieldVal.Interface())

		case "like":
			// Add % wildcard for LIKE
			strVal := fmt.Sprintf("%%%v%%", fieldVal.Interface())
			query = query.Where(fmt.Sprintf("%s LIKE ?", query_indentify), strVal)

		case "ilike":
			// Add % wildcard for LIKE
			strVal := fmt.Sprintf("%%%v%%", fieldVal.Interface())
			query = query.Where(fmt.Sprintf("%s LIKE ?", query_indentify), strVal)

		case "nlike":
			// Add % wildcard for LIKE
			strVal := fmt.Sprintf("%%%v%%", fieldVal.Interface())
			query = query.Where(fmt.Sprintf("%s LIKE ?", query_indentify), strVal)

		case "nin":
			query = query.Where(fmt.Sprintf("%s NOT IN ?", query_indentify), fieldVal.Interface())

		case "neq":
			query = query.Where(fmt.Sprintf("%s != ?", query_indentify), fieldVal.Interface())

		case "eq":
			query = query.Where(fmt.Sprintf("%s = ?", query_indentify), fieldVal.Interface())

		case "gte":
			query = query.Where(fmt.Sprintf("%s = ?", query_indentify), fieldVal.Interface())

		case "gt":
			query = query.Where(fmt.Sprintf("%s = ?", query_indentify), fieldVal.Interface())

		case "lte":
			query = query.Where(fmt.Sprintf("%s = ?", query_indentify), fieldVal.Interface())

		case "lt":
			query = query.Where(fmt.Sprintf("%s = ?", query_indentify), fieldVal.Interface())

		default:
			return nil, fmt.Errorf("unsupported operator %s in filter tag of columns \"%s\"", operator, query_indentify)
		}
	}

	field := val.FieldByName("Preloads")

	if !field.IsValid() || !field.CanInterface() {
		return db, nil // Field not found or not accessible
	}

	preloaders := []string{}

	switch field.Kind() {
	case reflect.Slice:
		elemKind := field.Type().Elem().Kind()
		switch elemKind {
		case reflect.String:
			for i := 0; i < field.Len(); i++ {
				preloaders = append(preloaders, field.Index(i).String())
			}
		case reflect.Ptr:
			if field.Type().Elem().Elem().Kind() == reflect.String {
				for i := 0; i < field.Len(); i++ {
					strPtr := field.Index(i).Interface().(*string)
					if strPtr != nil {
						preloaders = append(preloaders, *strPtr)
					}
				}
			}
		}
	}

	if len(preloaders) > 0 {
		var err error

		query, err = r.buildPreloader(query, preloaders)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

func (r *BaseRepository) buildPreloader(db *gorm.DB, preloaders []string) (*gorm.DB, error) {
	for _, preloader := range preloaders {
		db = db.Preload(preloader)
	}
	return db, nil
}

func (r *BaseRepository) GetReadDB(ctx context.Context) *gorm.DB {
	if ctx.Value("db") == nil {
		return r.db.Slave
	}
	return ctx.Value("db").(*gorm.DB)
}
func (r *BaseRepository) GetWriteDB(ctx context.Context) *gorm.DB {
	if ctx.Value("db") == nil {
		return r.db.Master
	}
	return ctx.Value("db").(*gorm.DB)
}
