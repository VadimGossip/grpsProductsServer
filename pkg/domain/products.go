package products

import (
	"errors"
	"time"

	"github.com/VadimGossip/grpsProductsServer/gen/products"
)

const (
	SORTINGFIELD_PRODUCT_NAME  = "product_name"
	SORTINGFIELD_PRICE         = "price"
	SORTINGFIELD_CHANGES_COUNT = "changes_count"
	SORTINGFIELD_TIMESTAMP     = "timestamp"

	SORTINGTYPE_ASC  = "asc"
	SORTINGTYPE_DESC = "desc"
)

var (
	sortingFields = map[string]products.SortingField{
		SORTINGFIELD_PRODUCT_NAME:  products.SortingField_product_name,
		SORTINGFIELD_PRICE:         products.SortingField_price,
		SORTINGFIELD_CHANGES_COUNT: products.SortingField_changes_count,
		SORTINGFIELD_TIMESTAMP:     products.SortingField_timestamp,
	}

	sortingTypes = map[string]products.SortingType{
		SORTINGTYPE_ASC:  products.SortingType_asc,
		SORTINGTYPE_DESC: products.SortingType_desc,
	}
)

type Product struct {
	Name         string    `bson:"product_name"`
	Price        int       `bson:"price"`
	ChangesCount int       `bson:"changes_count"`
	Timestamp    time.Time `bson:"timestamp"`
}

type PagingParams struct {
	Offset int
	Limit  int
}

type SortingFildsParams struct {
	Name         string
	Price        string
	ChangesCount string
	Timestamp    string
}

func NewSortingFiledParams() SortingFildsParams {
	return SortingFildsParams{
		Name:         "name",
		Price:        "price",
		ChangesCount: "changes_count",
		Timestamp:    "timestamp",
	}
}

func ToPbSortingFields(field string) (products.SortingField, error) {
	val, ex := sortingFields[field]
	if !ex {
		return 0, errors.New("invalid entity")
	}

	return val, nil
}

func ToPbSortingTypes(sortingType string) (products.SortingType, error) {
	val, ex := sortingTypes[sortingType]
	if !ex {
		return 0, errors.New("invalid entity")
	}

	return val, nil
}
