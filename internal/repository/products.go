package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VadimGossip/grpsProductsServer/pkg/domain"
)

type Products struct {
	db *mongo.Database
}

func NewProducts(db *mongo.Database) *Products {
	return &Products{
		db: db,
	}
}

func (r *Products) Insert(ctx context.Context, item products.Product) error {
	_, err := r.db.Collection("products").InsertOne(ctx, item)
	return err
}

func (r *Products) GetByName(ctx context.Context, name string) (products.Product, error) {
	var p products.Product
	filter := bson.D{{Key: "product_name", Value: name}}
	err := r.db.Collection("products").FindOne(ctx, filter).Decode(&p)
	return p, err
}

func (r *Products) UpdateByName(ctx context.Context, p products.Product) error {
	filter := bson.D{{Key: "product_name", Value: p.Name}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "price", Value: p.Price},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "changes_count", Value: 1},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "timestamp", Value: p.Timestamp},
		}},
	}

	_, err := r.db.Collection("products").UpdateOne(ctx, filter, update)
	return err
}

func (r *Products) List(ctx context.Context, paging products.PagingParams, sorting products.SortingParams) ([]products.Product, error) {
	opts := options.Find()
	sortingType := 1
	if sorting.SortType == "desc" {
		sortingType = -1
	}

	opts.SetSort(bson.D{
		{Key: sorting.Field, Value: sortingType},
	})
	opts.SetSkip(paging.Offset)
	opts.SetLimit(paging.Limit)

	cur, err := r.db.Collection("products").Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var productsList []products.Product
	for cur.Next(ctx) {
		var item products.Product
		if err := cur.Decode(&item); err != nil {
			return nil, err
		}
		productsList = append(productsList, item)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return productsList, nil
}
