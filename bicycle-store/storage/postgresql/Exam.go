package postgresql

import (
	"app/api/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

type R_Repo struct {
	db *pgxpool.Pool
}

func NewCodeRepo(db *pgxpool.Pool) *R_Repo {
	return &R_Repo{db: db}
}

func (r *R_Repo) Exam(req *models.StoreChange) (string, error) {
	apiCount, _ := strconv.Atoi(req.Count)
	var (
		give models.ReadFrom
		get  models.ReadTo
	)

	a, _ := r.ReadStocksF(context.Background(), req.Give_StoreId, req.ProdId)
	b, _ := r.ReadStocksG(context.Background(), req.Get_StoreId, req.ProdId)

	for _, v := range a {
		give = v
	}
	for _, v := range b {
		get = v
	}

	fmt.Println()
	fmt.Println("qaysi store dan olinadi = ", give)
	fmt.Println("qaysi store ga qoshiladi = ", get)
	fmt.Println()

	if give.Count-apiCount < 0 {
		fmt.Println("-- cant minus because count few null --")
		return "-- cant minus because count few null --", nil
	}

	// validation , bormi etgan store yomi
	vld, _ := r.Validator(context.Background())
	mapa := make(map[int]int)

	for _, v := range vld {
		mapa[v.Store] = v.Prod

	}

	nextStep := false

	for i, _ := range mapa {
		if i == get.Get_StoreId {
			nextStep = true
			break

		}

	}
	if nextStep == false {
		fmt.Println("stock unaqa store yoq")
		return "stock unaqa store yoq", nil

	}

	if nextStep == true {
		if give.Give_StoreId == get.Get_StoreId {
			fmt.Println("store oziga narsa jonatomidi")
			return "store oziga narsa jonatomidi", nil
		}
		get.Count += apiCount
		give.Count -= apiCount

		r.WriteChanged(context.Background(), give, get)

		fmt.Println("Changed")
		fmt.Println(r.ReadStocksF(context.Background(), req.Give_StoreId, req.ProdId))
		fmt.Println(r.ReadStocksG(context.Background(), req.Get_StoreId, req.ProdId))

	}

	return "all good", nil

}

func (r *R_Repo) ReadStocksF(ctx context.Context, idGive string, prod string) (from []models.ReadFrom, err error) {
	// Give
	rows, err := r.db.Query(ctx, ` SELECT stores.store_name,stocks.store_id,products.product_name,stocks.product_id, stocks.quantity
 FROM stocks
 JOIN products ON stocks.product_id = products.product_id
 JOIN stores ON stocks.store_id = stores.store_id
 WHERE stocks.product_id = $1  AND  stocks.store_id = $2 `, prod, idGive)

	if err != nil {
		fmt.Println("errore")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.ReadFrom
		err = rows.Scan(
			&a.Give_StoreName,
			&a.Give_StoreId,
			&a.ProdName,
			&a.ProdId,
			&a.Count,
		)

		from = append(from, a)
		if err != nil {
			return nil, err
		}

	}

	return from, err

}

func (r *R_Repo) ReadStocksG(ctx context.Context, idTo string, prod string) (to []models.ReadTo, err error) {
	// Get
	rows, err := r.db.Query(ctx, ` SELECT stores.store_name,stocks.store_id,products.product_name,stocks.product_id, stocks.quantity
 FROM stocks
 JOIN products ON stocks.product_id = products.product_id
 JOIN stores ON stocks.store_id = stores.store_id
 WHERE stocks.product_id = $1  AND  stocks.store_id = $2 `, prod, idTo)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.ReadTo
		err = rows.Scan(
			&a.Get_StoreName,
			&a.Get_StoreId,
			&a.ProdName,
			&a.ProdId,
			&a.Count,
		)

		to = append(to, a)
		if err != nil {
			return nil, err
		}

	}

	return to, err

}

func (r *R_Repo) WriteChanged(ctx context.Context, give models.ReadFrom, get models.ReadTo) (err error) {
	// Give logic
	queryGive := `UPDATE stocks SET quantity = $1 WHERE store_id = $2 AND product_id = $3`

	_, err = r.db.Exec(ctx, queryGive,
		give.Count,
		give.Give_StoreId,
		give.ProdId)

	if err != nil {
		return nil
	}

	// Get logic
	queryGet := `UPDATE stocks SET quantity = $1 WHERE store_id = $2 AND product_id = $3`

	_, err = r.db.Exec(ctx, queryGet,
		get.Count,
		get.Get_StoreId,
		get.ProdId)

	if err != nil {
		return nil
	}

	return err
}

func (r *R_Repo) Validator(ctx context.Context) (from []models.Valid, err error) {
	// Give
	rows, err := r.db.Query(ctx, ` SELECT store_id, product_id FROM stocks`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Valid
		err = rows.Scan(
			&a.Store,
			&a.Prod,
		)

		from = append(from, a)
		if err != nil {
			return nil, err
		}

	}

	return from, err

}
