package es

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/wayfair/model"
	"github.com/tealeg/xlsx"
)

const (
	_wayfairFamilyIDRegexExp = `\-(\w+)\.htm(l)?`
)

var (
	_wayfairFamilyIDRegex = regexp.MustCompile(_wayfairFamilyIDRegexExp)
)

func TestImport(t *testing.T) {
	excel, err := xlsx.OpenFile("./data.xlsx")
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}

	sheet := excel.Sheets[0]

	wfSKUs := make(map[string]string, len(sheet.Rows))

	for i, row := range sheet.Rows {
		if i < 1 {
			continue
		}

		status := strings.TrimSpace(row.Cells[6].Value)
		wfFamilyID := strings.TrimSpace(row.Cells[8].Value)
		eecItemNo := strings.TrimSpace(row.Cells[0].Value)
		if (eecItemNo != "" && !strings.Contains(eecItemNo, "N/A")) && (wfFamilyID != "" && !strings.Contains(wfFamilyID, "N/A")) && (status == "Active" || status == "Live Product") {
			itemsStr := wfSKUs[wfFamilyID]
			wfSKUs[wfFamilyID] = itemsStr + "," + eecItemNo
		}
		//  else {
		// 	log.Warnf("[%d] has empty family id", i)
		// }
	}

	// assert.Equal(t, len(sheet.Rows)-1, len(wfFamilyIDs))

	// wfSKUs = removeDuplicatedValues(wfSKUs)
	// log.Info(len(wfSKUs))

	wfItems := make([]*model.ItemDTO, 0, len(wfSKUs))

	for k, v := range wfSKUs {
		wfItems = append(wfItems, &model.ItemDTO{
			SKU:     k,
			ItemNOs: v + ",",
		})
	}

	esDAL, err := NewESItemDAL(
		elastic.SetURL("http://sa:Famous901@localhost:9200"),
		elastic.SetSniff(false),
	)
	if u.LogError(err) {
		return
	}

	err = esDAL.SaveItems(wfItems...)
	u.LogError(err)
}

func removeDuplicatedValues(src []string) []string {
	result := []string{}         //存放返回的不重复切片
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range src {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}
