package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	List(ctx context.Context, paging domain.PagingParams, sorting domain.SortingParams) ([]domain.Product, error)
}

type ProductService struct {
	repo Repository
}

func NewProductService(repo Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (ps *ProductService) Fetch(ctx context.Context, req *products.FetchRequest) (*emptypb.Empty, error) {
	data, err := csv.ParseCsvFromUrl(req.Url)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	for idx := range data {
		name := data[idx][0]
		price, err := strconv.ParseInt(data[idx][1], 10, 64)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		p := domain.Product{
			Name:         name,
			Price:        price,
			ChangesCount: 0,
			Timestamp:    time.Now(),
		}

		repoProduct, err := ps.repo.GetByName(ctx, p.Name)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				if err = ps.repo.Insert(ctx, p); err != nil {
					return &emptypb.Empty{}, err
				}
			}
		}
		if repoProduct.Price != p.Price {
			p.ChangesCount = repoProduct.ChangesCount + 1
			if err = ps.repo.UpdateByName(ctx, p); err != nil {
				return &emptypb.Empty{}, err
			}
		}
	}
	return &emptypb.Empty{}, nil
}

func (ps *ProductService) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	paging := domain.PagingParams{
		Offset: req.GetPagingOffset(),
		Limit:  req.GetPagingLimit(),
	}

	sorting := domain.SortingParams{
		Field:    req.GetSortField().String(),
		SortType: req.GetSortType().String(),
	}

	items, err := ps.repo.List(ctx, paging, sorting)
	if err != nil {
		return nil, err
	}

	var respProducts []*products.ProductItem
	for idx := range items {
		item := products.ProductItem{
			ProductName: items[idx].Name,
			Price:       items[idx].Price,
			Count:       items[idx].ChangesCount,
			Timestamp:   timestamppb.New(items[idx].Timestamp),
		}
		respProducts = append(respProducts, &item)
	}

	return &products.ListResponse{
		Product: respProducts,
	}, nil
}
