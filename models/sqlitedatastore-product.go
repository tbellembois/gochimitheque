package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"

	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

func (db *SQLiteDataStore) GetProducts(p GetProductsParameters) ([]Product, int, error) {
	var (
		products                                []Product
		count                                   int
		req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                  *sqlx.NamedStmt
		snstmt                                  *sqlx.NamedStmt
	)
	log.WithFields(log.Fields{"search": p.Search, "order": p.Order, "offset": p.Offset, "limit": p.Limit}).Debug("GetProducts")

	precreq.WriteString(" SELECT count(DISTINCT product.product_id)")
	presreq.WriteString(` SELECT product.product_id, 
	product.product_specificity, 
	name.name_label AS "name.name_label",
	casnumber.casnumber_label AS "casnumber.casnumber_label"`)
	comreq.WriteString(" FROM product")
	// get name
	comreq.WriteString(" JOIN name ON product.name = name.name_id")
	// get casnumber
	comreq.WriteString(" JOIN casnumber ON product.casnumber = casnumber.casnumber_id")
	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm, entity as e ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	`)
	comreq.WriteString(" WHERE name.name_label LIKE :search")
	postsreq.WriteString(" GROUP BY product.product_id")
	postsreq.WriteString(" ORDER BY name.name_label " + p.Order)

	// limit
	if p.Limit != constants.MaxUint64 {
		postsreq.WriteString(" LIMIT :limit OFFSET :offset")
	}

	// building count and select statements
	if cnstmt, db.err = db.PrepareNamed(precreq.String() + comreq.String()); db.err != nil {
		return nil, 0, db.err
	}
	if snstmt, db.err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); db.err != nil {
		return nil, 0, db.err
	}

	// building argument map
	m := map[string]interface{}{
		"search":   fmt.Sprint("%", p.Search, "%"),
		"personid": p.LoggedPersonID,
		"entityid": p.EntityID,
		"order":    p.Order,
		"limit":    p.Limit,
		"offset":   p.Offset}

	// select
	if db.err = snstmt.Select(&products, m); db.err != nil {
		return nil, 0, db.err
	}
	// count
	if db.err = cnstmt.Get(&count, m); db.err != nil {
		return nil, 0, db.err
	}

	//
	// getting symbols
	//
	for i, p := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT symbol_id, symbol_label, symbol_image FROM symbol")
		req.WriteString(" JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id")
		req.WriteString(" JOIN product ON productsymbols.productsymbols_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if db.err = db.Select(&products[i].Symbols, req.String(), p.ProductID); db.err != nil {
			return nil, 0, db.err
		}
	}

	return products, count, nil
}

func (db *SQLiteDataStore) GetProduct(id int) (Product, error) {
	var (
		product Product
	)
	return product, nil
}
func (db *SQLiteDataStore) DeleteProduct(id int) error {
	return nil
}
func (db *SQLiteDataStore) CreateProduct(p Product) (error, int) {
	return nil, 1
}
func (db *SQLiteDataStore) UpdateProduct(p Product) error {
	return nil
}
