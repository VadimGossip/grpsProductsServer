package service

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/VadimGossip/grpsProductsServer/gen/products"
	"github.com/VadimGossip/grpsProductsServer/pkg/csv"
	domain "github.com/VadimGossip/grpsProductsServer/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, item domain.Product) error
	GetByName(ctx context.Context, name string) (domain.Product, error)
	UpdateByName(ctx context.Context, prod domain.Product) error
	List(ctx context.Context, paging domain.PagingParams, sorting domain.SortingFildsParams) ([]domain.Product, error)
}

type ProductService struct {
	repo Repository
}

func NewProductService(repo Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (ps *ProductService) Fetch(ctx context.Context, req *products.FetchRequest) (*products.FetchResponse, error) {
	data, err := csv.ParseCsvFromUrl(req.Url)
	if err != nil {
		return nil, err
	}

	type changeItem struct {
		price   int
		counter int
	}

	changesMap := make(map[string]changeItem)
	for idx := range data {
		name := data[idx][0]
		price, err := strconv.Atoi(data[idx][1])
		if err != nil {
			return nil, err
		}
		if val, ok := changesMap[name]; ok {
			if val.price != price {
				changesMap[name] = changeItem{
					price:   price,
					counter: val.counter + 1,
				}
			}
		}
		changesMap[name] = changeItem{
			price:   price,
			counter: 1,
		}
	}

	for key, val := range changesMap {
		p := domain.Product{
			Name:         key,
			Price:        val.price,
			ChangesCount: val.counter,
			Timestamp:    time.Now(),
		}
		repoProduct, err := ps.repo.GetByName(ctx, p.Name)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				p.ChangesCount = 0
				if err = ps.repo.Insert(ctx, p); err != nil {
					return nil, err
				}
			}
		}
		p.ChangesCount += repoProduct.ChangesCount
		if err = ps.repo.UpdateByName(ctx, p); err != nil {
			return nil, err
		}

	}
	return &products.FetchResponse{
		Status: "ok",
	}, nil
}

func (ps *ProductService) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	//paging := products.PagingParams{
	//	Offset: int(req.GetPagingOffset()),
	//	Limit:  int(req.GetPagingLimit()),
	//}
	//sorting := products.SortingParams{
	//	Field: req.SortField,
	//	Asc:   req.SortAsc,
	//}
	//
	//items, err := s.repo.List(ctx, paging, sorting)
	//if err != nil {
	//	return nil, err
	//}
	//var sorted_products []*products.ProductItem
	//
	//for _, x := range items {
	//	var sorted_product products.ProductItem
	//	sorted_product.Name = x.Name
	//	sorted_product.Price = int32(x.Price)
	//	sorted_product.Count = int32(x.ChangesCount)
	//	sorted_product.Timestamp = timestamppb.New(x.Timestamp)
	//	sorted_products = append(sorted_products, &sorted_product)
	//}

	return &products.ListResponse{}, nil
}
