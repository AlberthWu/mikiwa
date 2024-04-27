package utils

type (
	Page struct {
		Page       int         `json:"page"`
		PageSize   int         `json:"pagesize"`
		TotalPage  int         `json:"totalpage"`
		TotalCount int         `json:"total"`
		FirstPage  bool        `json:"first_page"`
		LastPage   bool        `json:"last_page"`
		List       interface{} `json:"list"`
	}

	PageDynamic struct {
		Page             int         `json:"page"`
		PageSize         int         `json:"pagesize"`
		TotalPage        int         `json:"totalpage"`
		TotalCount       int         `json:"total"`
		FirstPage        bool        `json:"first_page"`
		LastPage         bool        `json:"last_page"`
		FieldKey         string      `json:"field_key"`
		FieldLabel       string      `json:"field_label"`
		FieldInt         string      `json:"field_int"`
		FieldLevel       string      `json:"field_level"`
		FieldExport      string      `json:"field_export"`
		FieldExportLabel string      `json:"field_export_label"`
		FieldFooter      string      `json:"field_footer"`
		List             interface{} `json:"list"`
	}
)

func Pagination(count int, page int, pageSize int, list interface{}) Page {

	tp := count / pageSize

	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}

	return Page{Page: page, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: page == 1, LastPage: page == tp, List: list}
}

func PaginationDynamic(count int, page int, pageSize int, fieldKey, fieldLabel, fieldInt, fieldLevel, fieldExport, fieldExportLabel, fieldFooter string, list interface{}) PageDynamic {

	tp := count / pageSize

	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}

	return PageDynamic{Page: page, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: page == 1, LastPage: page == tp, FieldKey: fieldKey, FieldLabel: fieldLabel, FieldInt: fieldInt, FieldLevel: fieldLevel, FieldExport: fieldExport, FieldExportLabel: fieldExportLabel, FieldFooter: fieldFooter, List: list}
}
