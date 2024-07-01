package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type DoController struct {
	base.BaseController
}

func (c *DoController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderDo struct {
		SalesOrderId    int                   `json:"sales_order_id"`
		IssueDate       string                `json:"issue_date"`
		WarehouseId     int                   `json:"warehouse_id"`
		CustomerId      int                   `json:"customer_id"`
		PlantId         int                   `json:"plant_id"`
		LogisticsId     int                   `json:"logistics_id"`
		CourierId       int                   `json:"courier_id"`
		StatusId        int8                  `json:"status_id"`
		DeliveryAddress string                `json:"delivery_address"`
		PlateNo         string                `json:"plate_no"`
		Notes           string                `json:"notes"`
		UploadFile      models.DocumentList   `json:"upload_file"`
		Detail          []InputDetailDoDetail `json:"detail"`
	}

	InputDetailDoDetail struct {
		Id         int     `json:"id"`
		CategoryId int     `json:"category_id"`
		ProductId  int     `json:"product_id"`
		Qty        float64 `json:"qty"`
		UomId      int     `json:"uom_id"`
		Memo       string  `json:"memo"`
	}
)

func (c *DoController) Post() {
	// o := orm.NewOrm()
	// var user_id, form_id int
	// var user_name string
	// var err error
	// var folderName string = "delivery_order"
	// var status_id int8 = base.OpenDo
	// sess := c.GetSession("profile")
	// if sess != nil {
	// 	user_id = sess.(map[string]interface{})["id"].(int)
	// 	user_name = sess.(map[string]interface{})["username"].(string)
	// }

	// form_id = base.FormName(form_do)
	// write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	// write_aut = true
	// if !write_aut {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
	// 	c.ServeJSON()
	// 	return
	// }

	// var i int = 0
	// var ob InputHeaderDo
	// var inputDetail []models.DoDetail

	// body := c.Ctx.Input.RequestBody
	// err = json.Unmarshal(body, &ob)
	// if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 	c.ServeJSON()
	// 	return
	// }

	// ob.CourierId = 1
	// ob.StatusId = status_id
	// valid := validation.Validation{}
	// valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	// valid.Required(ob.WarehouseId, "warehouse_id").Message("Is required")
	// valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	// valid.Required(ob.LogisticsId, "logistics_id").Message("Is required")
	// valid.Required(ob.PlateNo, "plate_no").Message("Is required")
	// valid.Required(strings.TrimSpace(ob.DeliveryAddress), "delivery_address").Message("Is required")

	// if len(ob.Detail) == 0 {
	// 	valid.AddError("detail", "Detail list is required")
	// }

	// if valid.HasErrors() {
	// 	out := make([]utils.ApiError, len(valid.Errors))
	// 	for i, err := range valid.Errors {
	// 		out[i] = utils.ApiError{Param: err.Key, Message: err.Message}
	// 	}
	// 	c.Ctx.ResponseWriter.WriteHeader(400)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 400, "Invalid input field", out)
	// 	c.ServeJSON()
	// 	return
	// }

	// var customers models.Company
	// var plants models.Plant
	// err = models.Companies().Filter("id", ob.CustomerId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Customer).One(&customers)
	// if err == orm.ErrNoRows {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
	// 	c.ServeJSON()
	// 	return
	// }

	// if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 	c.ServeJSON()
	// 	return
	// }

	// if customers.Status == 0 {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", customers.Code))
	// 	c.ServeJSON()
	// 	return
	// }

	// if ob.PlantId != 0 {
	// 	err = models.Plants().Filter("deleted_at__isnull", true).Filter("id", ob.PlantId).Filter("company_id", ob.CustomerId).One(&plants)
	// 	if err == orm.ErrNoRows {
	// 		c.Ctx.ResponseWriter.WriteHeader(401)
	// 		utils.ReturnHTTPError(&c.Controller, 401, "Plant unregistered/Illegal data")
	// 		c.ServeJSON()
	// 		return
	// 	}

	// 	if err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(401)
	// 		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 		c.ServeJSON()
	// 		return
	// 	}

	// 	if plants.Status == 0 {
	// 		c.Ctx.ResponseWriter.WriteHeader(402)
	// 		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", plants.Name))
	// 		c.ServeJSON()
	// 		return
	// 	}
	// }

	// var warehouse models.Company
	// err = models.Companies().Filter("id", ob.WarehouseId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Warehouse).One(&warehouse)
	// if err == orm.ErrNoRows {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, "Warehouse unregistered/Illegal data")
	// 	c.ServeJSON()
	// 	return
	// }

	// if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 	c.ServeJSON()
	// 	return
	// }

	// if warehouse.Status == 0 {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", warehouse.Code))
	// 	c.ServeJSON()
	// 	return
	// }

	// var logistics models.Company
	// err = models.Companies().Filter("id", ob.LogisticsId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Transporter).One(&logistics)
	// if err == orm.ErrNoRows {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, "Logistic unregistered/Illegal data")
	// 	c.ServeJSON()
	// 	return
	// }

	// if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 	c.ServeJSON()
	// 	return
	// }

	// if logistics.Status == 0 {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", logistics.Code))
	// 	c.ServeJSON()
	// 	return
	// }

	// var products models.Product
	// var productUom models.ProductUom
	// wg = new(sync.WaitGroup)
	// var mutex sync.Mutex
	// resultChan := make(chan utils.ResultChan, len(ob.Detail))
	// var queryResults []utils.ResultChan
	// wg.Add(len(ob.Detail))
	// for _, v := range ob.Detail {
	// 	go func(v InputDetailDoDetail) {
	// 		defer wg.Done()
	// 		mutex.Lock()
	// 		defer mutex.Unlock()
	// 		if err = models.Products().Filter("id", v.ProductId).Filter("deleted_at__isnull", true).Filter("product_type_id", 3).One(&products); err == orm.ErrNoRows {
	// 			resultChan <- utils.ResultChan{Id: v.ProductId, Data: "Invalid product", Message: "product unregistered/Illegal data"}
	// 			return
	// 		}

	// 		if err != nil {
	// 			resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
	// 			return
	// 		}

	// 		if products.StatusId == 0 {
	// 			resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: fmt.Sprintf("'%v' has been set as inactive", products.ProductCode)}
	// 			return
	// 		}

	// 		err = o.Raw("select * from product_uom where product_id = " + utils.Int2String(v.ProductId) + " and uom_id = " + utils.Int2String(v.UomId)).QueryRow(&productUom)
	// 		if err == orm.ErrNoRows {
	// 			resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: fmt.Sprintf("product uom unregistered/Illegal data")}
	// 			return
	// 		}

	// 		if err != nil {
	// 			resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
	// 			return
	// 		}
	// 	}(v)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(resultChan)
	// }()

	// for result := range resultChan {
	// 	if result.Message != "" {
	// 		queryResults = append(queryResults, utils.ResultChan{
	// 			Id:      result.Id,
	// 			Data:    result.Data,
	// 			Message: result.Message,
	// 		})
	// 	}
	// }

	// if len(queryResults) != 0 {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Error", map[string]interface{}{"Invalid field": queryResults})
	// 	c.ServeJSON()
	// 	return
	// }

	// thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	// if errDate != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
	// 	c.ServeJSON()
	// 	return
	// }

	// seqno, referenceno := models.GenerateNumber(thedate, 1, ob.CustomerId)

	// t_sales_order = models.SalesOrder{
	// 	IssueDate:       thedate,
	// 	ReferenceNo:     referenceno,
	// 	SeqNo:           seqno,
	// 	DueDate:         dueDate,
	// 	PoolId:          1,
	// 	OutletId:        ob.OutletId,
	// 	OutletName:      outlet.Name,
	// 	CustomerId:      ob.CustomerId,
	// 	CustomerCode:    customers.Code,
	// 	PlantId:         ob.PlantId,
	// 	PlantName:       plants.Name,
	// 	Terms:           customers.Terms,
	// 	DeliveryAddress: ob.DeliveryAddress,
	// 	EmployeeId:      ob.EmployeeId,
	// 	EmployeeName:    "",
	// 	LeadTime:        ob.LeadTime,
	// 	Subtotal:        subtotal_,
	// 	TotalDisc:       totalDisc_,
	// 	Dpp:             dpp_amount,
	// 	Ppn:             ppn,
	// 	PpnAmount:       ppn_amount,
	// 	Total:           total,
	// 	StatusId:        ob.StatusId,
	// 	CreatedBy:       user_name,
	// 	UpdatedBy:       user_name,
	// }

	// for k, v := range ob.Detail {
	// 	i = 0
	// 	wg.Add(1)
	// 	go func(k int, v InputDetailSalesOrder) {
	// 		priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)
	// 		if priceRtn == nil {
	// 			disc1 = 0
	// 			disc2 = 0
	// 			disctpr = 0
	// 			nettprice = 0
	// 			subtotal = 0
	// 			price = 0
	// 			normal_price = 0
	// 		} else {
	// 			if priceRtn.Price == 0 {
	// 				disc1 = 0
	// 				disc2 = 0
	// 				disctpr = 0
	// 				nettprice = 0
	// 				subtotal = 0
	// 				price = 0
	// 				normal_price = 0
	// 			} else {
	// 				disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
	// 				disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
	// 				disctpr = v.DiscTpr * -1
	// 				price = priceRtn.Price
	// 				normal_price = priceRtn.NormalPrice
	// 				nettprice = price + disc1 + disc2 + disctpr
	// 				subtotal = priceRtn.FinalQty * nettprice
	// 			}

	// 		}
	// 		defer wg.Done()
	// 		mutex.Lock()
	// 		if v.Id == 0 {
	// 			inputDetail = append(inputDetail, models.SalesOrderDetail{
	// 				SalesOrderId:      t_sales_order.Id,
	// 				ReferenceNo:       referenceno,
	// 				IssueDate:         thedate,
	// 				DueDate:           dueDate,
	// 				ItemNo:            k + 1,
	// 				ProductId:         v.ProductId,
	// 				ProductCode:       priceRtn.ProductCode,
	// 				Qty:               v.Qty,
	// 				UomId:             v.UomId,
	// 				UomCode:           priceRtn.UomCode,
	// 				Ratio:             priceRtn.Ratio,
	// 				PackagingId:       priceRtn.PackagingId,
	// 				PackagingCode:     priceRtn.PackagingCode,
	// 				FinalQty:          priceRtn.FinalQty,
	// 				FinalUomId:        priceRtn.FinalUomId,
	// 				FinalUomCode:      priceRtn.FinalUomCode,
	// 				NormalPrice:       normal_price,
	// 				PriceId:           priceRtn.PriceId,
	// 				Price:             price,
	// 				Disc1:             v.Disc1,
	// 				Disc1Amount:       disc1,
	// 				Disc2:             v.Disc2,
	// 				Disc2Amount:       disc2,
	// 				DiscTpr:           disctpr,
	// 				TotalDisc:         disc1 + disc2 + disctpr,
	// 				NettPrice:         nettprice,
	// 				Total:             subtotal,
	// 				LeadTime:          v.LeadTime,
	// 				ConversionQty:     priceRtn.ConversionQty,
	// 				ConversionUomId:   priceRtn.ConversionUomId,
	// 				ConversionUomCode: priceRtn.ConversionUomCode,
	// 				CreatedBy:         user_name,
	// 				UpdatedBy:         user_name,
	// 			})
	// 			i += 1
	// 		}
	// 		mutex.Unlock()
	// 	}(k, v)
	// }
	// wg.Wait()

	// d, err_ := t_sales_order.InsertWithDetail(t_sales_order, inputDetail)
	// errcode, errmessage = base.DecodeErr(err_)
	// if err_ != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(errcode)
	// 	utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	// 	c.ServeJSON()
	// 	return
	// } else {
	// 	if err = base.PostFirebaseRaw(ob.UploadFile, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id)); err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(401)
	// 		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
	// 	} else {
	// 		v, err := t_sales_order.GetById(d.Id, user_id)
	// 		errcode, errmessage = base.DecodeErr(err)
	// 		if err != nil {
	// 			c.Ctx.ResponseWriter.WriteHeader(errcode)
	// 			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	// 		} else {
	// 			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	// 		}
	// 	}
	// }

	// c.ServeJSON()
}
