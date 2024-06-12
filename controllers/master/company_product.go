package controllers

// import (
// 	"encoding/json"
// 	"fmt"
// 	base "mikiwa/controllers"
// 	"mikiwa/models"
// 	"mikiwa/utils"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/beego/beego/v2/client/orm"
// )

// type CompanyProductController struct {
// 	base.BaseController
// }

// func (c *CompanyProductController) Prepare() {
// 	c.Ctx.Request.Header.Set("token", "No Aut")
// 	c.BaseController.Prepare()
// }

// type (
// 	InputHeaderCompanyProduct struct {
// 		Id            int     `json:"id"`
// 		EffectiveDate string  `json:"effective_date"`
// 		ExpiredDate   string  `json:"expired_date"`
// 		ProductId     int     `json:"product_id"`
// 		SureName      string  `json:"sure_name"`
// 		UomIdOne      int     `json:"uom_id_one"`
// 		Ratio         float64 `json:"ratio"`
// 		DiscOne       float64 `json:"disc_one"`
// 		DiscTwo       float64 `json:"disc_two"`
// 		DiscTpr       float64 `json:"disc_tpr"`
// 		StatusId      int8    `json:"status_id"`
// 	}
// )

// func (c *CompanyProductController) Post() {
// 	o := orm.NewOrm()
// 	var user_id, form_id int
// 	var user_name string
// 	var err error
// 	sess := c.GetSession("profile")
// 	if sess != nil {
// 		user_id = sess.(map[string]interface{})["id"].(int)
// 		user_name = sess.(map[string]interface{})["username"].(string)
// 	}

// 	form_id = base.FormName(form_company_product)
// 	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
// 	write_aut = true
// 	if !write_aut {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
// 		c.ServeJSON()
// 		return
// 	}

// 	var i int = 0
// 	var ob []InputHeaderCompanyProduct
// 	var input []models.CompanyProduct

// 	body := c.Ctx.Input.RequestBody
// 	err = json.Unmarshal(body, &ob)
// 	if err != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
// 		c.ServeJSON()
// 		return
// 	}

// 	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
// 	var querydata models.Company
// 	err = models.Companies().Filter("id", id).Filter("deleted_at__isnull", true).Filter("status", 1).Filter("CompanyTypes__TypeId__Id", base.Customer).One(&querydata)
// 	if err == orm.ErrNoRows {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPError(&c.Controller, 402, "Customer id unregistered/Illegal data")
// 		c.ServeJSON()
// 		return
// 	}

// 	if err != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
// 		c.ServeJSON()
// 		return
// 	}

// 	var deleteIds []string
// 	var joinId string
// 	var issuedate time.Time
// 	var expireddate *time.Time
// 	var errdate error
// 	for _, v := range ob {
// 		issuedate, errdate = time.Parse("2006-01-02", v.EffectiveDate)
// 		if errdate != nil {
// 			c.Ctx.ResponseWriter.WriteHeader(401)
// 			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date :", errdate.Error()))
// 			c.ServeJSON()
// 			return
// 		}

// 		if v.ExpiredDate == "" {
// 			expireddate = nil
// 		} else {
// 			expiredthedate, errdate := time.Parse("2006-01-02", v.ExpiredDate)
// 			if errdate != nil {
// 				c.Ctx.ResponseWriter.WriteHeader(401)
// 				utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date: ", errdate.Error()))
// 				c.ServeJSON()
// 				return
// 			}
// 			expireddate = &expiredthedate
// 		}

// 		if v.Id != 0 {
// 			deleteIds = append(deleteIds, utils.Int2String(v.Id))
// 		}

// 		deleteIds = append(deleteIds, utils.Int2String(v.Id))
// 	}

// 	if len(deleteIds) == 0 {
// 		joinId = "0"
// 	} else {
// 		joinId = strings.Join(deleteIds, ",")
// 	}
// 	o.Raw("update company_product set deleted_at = now(), deleted_by = '" + user_name + "' where deleted_at is null and company_id  = " + utils.Int2String(id) + " and id not in (" + joinId + ")").Exec()

// 	var products models.Product
// 	var wg = new(sync.WaitGroup)
// 	var mutex sync.Mutex
// 	for k, v := range ob {
// 		wg.Add(1)
// 		go func(k int, v InputHeaderCompanyProduct) {
// 			defer wg.Done()
// 			mutex.Lock()
// 			models.Products().Filter("id", v.ProductId).Filter("deleted_at__isnull", true).Filter("product_type_id", 3).One(&products)
// 			if v.Id == 0 {
// 				input = append(input, models.CompanyProduct{
// 					EffectiveDate: issuedate,
// 					ExpiredDate:   expireddate,
// 					CompanyId:     id,
// 					CompanyCode:   querydata.Code,
// 					ProductId:     v.ProductId,
// 					ProductCode:   products.ProductCode,
// 					NormalPrice:   products.Price,
// 				})
// 				i += 1
// 			} else {
// 				t_company_product.Id = v.Id
// 				t_company_product.IssueDate = issuedate
// 				t_company_product.ItemNo = k + 1
// 				t_company_product.Update()
// 			}
// 			mutex.Unlock()
// 		}(k, v)
// 	}

// 	wg.Wait()
// 	o.InsertMulti(i, input)

// 	v, err_ := t_fleet_attendance.JobList(issue_date, "", user_id, 0, 1)
// 	errcode, errmessage = base.DecodeErr(err_)
// 	if err_ != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(errcode)
// 		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
// 	} else {
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
// 	}
// 	c.ServeJSON()
// }

// func (c *CompanyProductController) Delete() {}
// func (c *CompanyProductController) GetOne() {}
// func (c *CompanyProductController) GetAll() {}
