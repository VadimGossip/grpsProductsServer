package repository

import (
	"context"
	//	"fmt"
	//	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VadimGossip/grpsProductsServer/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	filter := bson.D{{Key: "name", Value: p.Name}}
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

func (r *Products) List(ctx context.Context, paging products.PagingParams, sorting products.SortingFildsParams) ([]products.Product, error) {
	//opts := options.Find()
	//sortOpts := bson.D{
	//	{Key: fmt.Sprintf("%v", sorting.Field), Value: sorting.Asc},
	//}
	//
	//opts.SetSort(sortOpts)
	//opts.SetSkip(int64(paging.Offset))
	//opts.SetLimit(int64(paging.Limit))
	//
	//cur, err := r.db.Collection("products").Find(ctx, bson.D{}, opts)
	//if err != nil {
	//	return nil, err
	//}
	//defer cur.Close(ctx)
	//
	//var productsList []product.Product
	//for cur.Next(ctx) {
	//	var elem product.Product
	//	if err := cur.Decode(&elem); err != nil {
	//		return nil, err
	//	}
	//	productsList = append(productsList, elem)
	//}
	//
	//if err := cur.Err(); err != nil {
	//	return nil, err
	//}

	return nil, nil
}
