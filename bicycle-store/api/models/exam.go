package models

type StoreChange struct {
	Give_StoreId string `json:"give_store_id"`
	Get_StoreId  string `json:"get_store_id"`
	ProdId       string `json:"prod_id"`
	Count        string `json:"count"`
}

type ReadFrom struct {
	Give_StoreId   int    `json:"store_id"`
	Give_StoreName string `json:" store_name"`
	ProdId         int    `json:"product_id"`
	ProdName       string `json:"product_name"`
	Count          int    `json:"quantity"`
}

type ReadTo struct {
	Get_StoreId   int    `json:"store_id"`
	Get_StoreName string `json:" store_name"`
	ProdId        int    `json:"product_id"`
	ProdName      string `json:"product_name"`
	Count         int    `json:"quantity"`
}

type Valid struct {
	Store int `json:"store_id"`
	Prod  int `json:"product_id"`
}
