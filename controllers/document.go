package controllers

// import (
// 	"fmt"
// 	"mikiwa/models"
// 	"mikiwa/utils"
// 	"mime/multipart"
// 	"strings"
// 	"time"

// 	"github.com/beego/beego/v2/client/orm"
// 	"github.com/beego/beego/v2/core/validation"
// )

// type DcoumentController struct {
// 	BaseController
// }

// func (c *DcoumentController) Prepare() {
// 	c.Ctx.Request.Header.Set("token", "No Aut")
// 	c.BaseController.Prepare()
// }

// var form_document = "document"
// var t_document models.Document

// func (c *DcoumentController) Post() {
// 	var user_id, form_id int
// 	var user_name string
// 	var err error
// 	var deletedat string

// 	fmt.Print("Check :", user_id, form_id, user_name, "..")
// 	sess := c.GetSession("profile")
// 	if sess != nil {
// 		user_name = sess.(map[string]interface{})["username"].(string)
// 		user_id = sess.(map[string]interface{})["id"].(int)
// 	}
// 	form_id = FormName(form_document)
// 	write_aut := models.CheckPrivileges(user_id, form_id, Write)
// 	if !write_aut {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
// 		c.ServeJSON()
// 		return
// 	}

// 	fleet_id, _ := c.GetInt("fleet_id")
// 	doc_type_id, _ := c.GetInt("doc_type_id")
// 	folder_name := strings.TrimSpace(c.GetString("folder_name"))
// 	issue_date := strings.TrimSpace(c.GetString("issue_date"))
// 	expired_date := strings.TrimSpace(c.GetString("expired_date"))
// 	amount, _ := c.GetFloat("amount")
// 	files, errf := c.GetFiles("file")
// 	if errf != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("invalid_file :", errf.Error()))
// 		c.ServeJSON()
// 		return
// 	}

// 	valid := validation.Validation{}
// 	valid.Required(fleet_id, "fleet_id").Message("Is required")
// 	valid.Required(doc_type_id, "doc_type_id").Message("Is required")
// 	valid.Required(doc_no, "doc_no").Message("Is required")
// 	valid.Required(issue_date, "issue_date").Message("Is required")
// 	valid.Required(expired_date, "expired_date").Message("Is required")
// 	valid.Required(files, "files").Message("Is required")

// 	if valid.HasErrors() {
// 		out := make([]errors.ApiError, len(valid.Errors))
// 		for i, err := range valid.Errors {
// 			out[i] = errors.ApiError{Param: err.Key, Message: err.Message}
// 		}
// 		c.Ctx.ResponseWriter.WriteHeader(400)
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 400, "Invalid input field", out)
// 		c.ServeJSON()
// 		return
// 	}

// 	var fleets models.Fleet
// 	if err = models.Fleets().Filter("id", fleet_id).One(&fleets); err == orm.ErrNoRows {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, "Fleet unregistered/Illegal data")
// 		c.ServeJSON()
// 		return
// 	}

// 	if err != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprint("fleet_id :", err.Error()))
// 		c.ServeJSON()
// 		return
// 	}

// 	deletedat = fleets.DeletedAt.Format("2006-01-02")
// 	if deletedat != "0001-01-01" {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("fleet_id :'%v' has been deleted", fleets.PlateNo))
// 		c.ServeJSON()
// 		return
// 	}

// 	if fleets.Status == 0 {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid fleet", map[string]interface{}{"fleet_id": "'" + fleets.PlateNo + "' has been set as inactive"})
// 		c.ServeJSON()
// 		return
// 	}

// 	var doctype models.FleetDocType
// 	if err = models.FleetDocTypes().Filter("id", doc_type_id).Filter("deleted_at__isnull", true).One(&doctype); err == orm.ErrNoRows {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, "Doc type unregistered/Illegal data")
// 		c.ServeJSON()
// 		return
// 	}

// 	if err != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
// 		c.ServeJSON()
// 		return
// 	}

// 	deletedat = doctype.DeletedAt.Format("2006-01-02")
// 	if deletedat != "0001-01-01" {
// 		c.Ctx.ResponseWriter.WriteHeader(402)
// 		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("doc_type_id :'%v' has been deleted", doctype.Names))
// 		c.ServeJSON()
// 		return
// 	}

// 	issuedate, errDate := time.Parse("2006-01-02", issue_date)
// 	if errDate != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date :", errDate.Error()))
// 		c.ServeJSON()
// 		return
// 	}

// 	expdate, errExpDate := time.Parse("2006-01-02", expired_date)
// 	if errExpDate != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(401)
// 		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date :", errExpDate.Error()))
// 		c.ServeJSON()
// 		return
// 	}

// 	t_fleet_document = models.FleetDocument{
// 		FleetId:     fleet_id,
// 		IssueDate:   issuedate,
// 		ExpiredDate: expdate,
// 		DocTypeId:   doc_type_id,
// 		DocType:     doctype.Names,
// 		DocNo:       doc_no,
// 		Amount:      amount,
// 		CreatedBy:   user_name,
// 		UpdatedBy:   user_name,
// 	}

// 	d, err_ := t_fleet_document.Insert(t_fleet_document)
// 	errcode, errmessage := base.DecodeErr(err_)
// 	if err_ != nil {
// 		c.Ctx.ResponseWriter.WriteHeader(errcode)
// 		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
// 	} else {
// 		// insert document
// 		//  base.PutFilesFirebase(files, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id))
// 		for _, fileHeader := range files {
// 			if err := base.PostFilesToFirebase([]*multipart.FileHeader{fileHeader}, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id)); err != nil {
// 				c.Ctx.ResponseWriter.WriteHeader(401)
// 				utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error while posting files to Firebase: %s", err.Error()))
// 				c.ServeJSON()
// 				return
// 			}
// 		}

// 		v, _ := t_fleet_document.GetById(d.Id)
// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
// 	}

// 	c.ServeJSON()
// }

// func (c *DcoumentController) Put() {}

// func (c *DcoumentController) Delete() {}

// func (c *DcoumentController) GetAll() {}
// func (c *DcoumentController) GetOne() {}
