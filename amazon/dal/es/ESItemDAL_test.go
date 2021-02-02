package es

import (
	"fmt"
	"strings"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/model"
	"github.com/tealeg/xlsx"
)

func TestESItemDAL_ImportItems(t *testing.T) {
	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://192.168.188.200:9200"),
	)
	if u.LogError(err) {
		return
	}

	excel, err := xlsx.OpenFile("./data.xlsx")
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}

	sheet := excel.Sheets[0]
	items := make([]*model.ItemDTO, 0, len(sheet.Rows))
	for i, row := range sheet.Rows {
		if i <= 1 {
			continue
		}
		var strs []string
		for _, cell := range row.Cells {
			text := cell.String()
			strs = append(strs, text)
		}
		items = append(items, &model.ItemDTO{
			ItemNo: strings.TrimSpace(strs[1]),
			ASIN:   strings.TrimSpace(strs[3]),
		})
	}

	err = esDAL.SaveItems(items...)
	u.LogError(err)
}

func TestESItemDAL_GetItems(t *testing.T) {
	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://localhost:9200"),
	)
	if u.LogError(err) {
		return
	}

	rs, err := esDAL.GetItems(&model.ItemQuery{
		ASIN:     "AAAAAAA",
		ItemNo:   "Item0001",
		PageSize: 10000,
		Status:   2,
	})
	u.LogError(err)
	assert.NotEmpty(t, rs)
}

func TestESItemDAL_GetAllItems(t *testing.T) {
	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://localhost:9200"),
	)
	if u.LogError(err) {
		return
	}

	rs, err := esDAL.GetAllItems(&model.ItemQuery{
		Status: -1,
	})
	u.LogError(err)
	assert.NotEmpty(t, rs)
}

func TestESItemDAL_SaveItems(t *testing.T) {
	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://localhost:9200"),
	)
	if u.LogError(err) {
		return
	}

	err = esDAL.SaveItems(&model.ItemDTO{
		ASIN:   "AAAAAAA",
		ItemNo: "Item0001",
		Status: 2,
	})
	u.LogError(err)
}

func TestESItemDAL_DeleteItems(t *testing.T) {
	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://localhost:9200"),
	)
	if u.LogError(err) {
		return
	}

	err = esDAL.DeleteItems(&model.ItemDTO{
		ASIN:   "AAAAAAA",
		Status: 2,
	})
	u.LogError(err)
}
