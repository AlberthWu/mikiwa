SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_SalesOrder;
DELIMITER $$
CREATE PROCEDURE sp_SalesOrder(theDate date,dueDate date,updatedAt date,uId int,employeeIds varchar(15),outletIds varchar(15),customerIds varchar(15),plantId int,productIds varchar(15),statusIds varchar(15),reportTypeId int,userId int,keyword varchar(255),in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000), in limitVal int, in offsetVal int ) 
BEGIN
	declare keywordSet varchar(8000);
    declare theDateSet varchar(100);
    declare dueDateSet varchar(255);
    declare updatedAtSet   varchar(255);
    declare userIdSet   varchar(255);
    declare uIdSet   varchar(255);
    declare employeeIdsSet   varchar(255);
	declare outletIdsSet   varchar(255);
	declare customerIdsSet   varchar(255);
	declare plantIdSet   varchar(255);
	declare productIdsSet   varchar(255);
    declare statusIdSet   varchar(255);
    declare limitSet  varchar(255);
    
    DECLARE Rn int;
    Declare Urut int DEFAULT 1;
    Declare ColumnName varchar(255) ;
    Declare ParameterName varchar(255);
    Declare KeywordName varchar(255);
    DECLARE ColumnSet varchar(8000) default '';
    
    set Rn = (CHAR_LENGTH(TheField) - CHAR_LENGTH(REPLACE(TheField, ',','')) + 1) ;
    
	WHILE Urut <= Rn DO
		Set ColumnName = SUBSTRING_INDEX(SUBSTRING_INDEX(TheField, ',', Urut),',',-1);
		Set ParameterName = SUBSTRING_INDEX(SUBSTRING_INDEX(MatchMode, ',', Urut),',',-1);
		Set KeywordName = SUBSTRING_INDEX(SUBSTRING_INDEX(ValueName, ',', Urut),',',-1);
		IF UPPER(ParameterName) = "CONTAINS" THEN
			set ColumnSet =  concat(ColumnSet," AND ",ColumnName," like '%",KeywordName,"%'");
		ELSEIF   UPPER(ParameterName) = "STARTSWITH" THEN
			set ColumnSet = concat(ColumnSet," AND ",ColumnName," like '",KeywordName,"%'");
		ELSEIF   UPPER(ParameterName) = "NOTCONTAINS" THEN
			set ColumnSet = concat(ColumnSet," AND ",ColumnName," not like '%",KeywordName,"%'");
		ELSEIF   UPPER(ParameterName) = "ENDSWITH" THEN
			set ColumnSet = concat(ColumnSet," AND ",ColumnName," like '%",KeywordName,"'");
		ELSEIF   UPPER(ParameterName) = "EQUALS" THEN
			set ColumnSet = concat(ColumnSet," AND ",ColumnName," = '",KeywordName,"'");
		ELSEIF   UPPER(ParameterName) = "NOTEQUALS" THEN
			set ColumnSet = concat(ColumnSet," AND ",ColumnName," not in ('",KeywordName,"')");
		end if;
		set Urut = Urut + 1;
	END WHILE;    
	SET ColumnSet = case when (ColumnSet is null or ColumnSet = '') then '' else ColumnSet end;
    
    set theDateSet = case when theDate is null then '' else concat("  and issue_date >= '",TheDate,"' issue_date <= '",TheDate,"' ")   end;
    set dueDateSet = case when dueDate is null then '' else concat("  and due_date >= '",dueDate,"' due_date <= '",dueDate,"' ")   end;
	SET updatedAtSet  = case when updatedAt is null  then '' else concat(" and (date(t0.created_at) = date('",updatedAt,"')  or date(ifnull(t0.updated_at,t0.created_at)) = date('",updatedAt,"'))")   end;
    set userIdSet = '';
	SET uIdSet  = case when (uId is null  or uId = 0) then '' else concat(" and t0.id in (",uId,")")   end;
    SET statusIdSet  = case when (statusIds is null or statusIds = '') then '' else concat(" and t0.status_id in (",replace(statusIds,'''',''),")") end ;
    set employeeIdsSet = case when (employeeIds is null or employeeIds = '') then '' else concat(" AND t0.employee_id in (",replace(employeeIds,'''',''),")") end ;
    set outletIdsSet = case when (outletIds is null or outletIds = '') then '' else concat(" and (t0.outlet_id in (",replace(outletIds,'''',''),")") end ;
    set customerIdsSet = case when (customerIds is null or customerIds = '') then '' else concat(" and (t0.customer_id in (",replace(customerIds,'''',''),")") end ;
    set plantIdSet = case when (plantId is null or plantId = 0) then '' else concat(" and (t0.plant_id = ",plantId) end ;
    set productIdsSet = case when (productIds is null or productIds = '') then '' else concat(" and t0.id in (select product_id from sales_order_detail where deleted_at is null and t0.product_id in (",replace(productIds,'''',''),"))") end ;
   --  SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and (t1.product_code like '%",keyword,"%'  or t1.product_name like '%",keyword,"%' or t3.company_code like '%",keyword,"%' or t3.company_name like '%",keyword,"%' 
-- 		 or t1.division_code like '%",keyword,"%'   or t1.product_type_name like '%",keyword,"%' or t2.uom_code  like '%",keyword,"%')")   end;
	set keywordSet= '';
    SET limitSet = case when LimitVal is null then '' else concat(" limit ",offsetVal,",",limitVal)    end;
	if reportTypeId = 0 then
		set uIdSet = case when (uId is null  or uId = 0) then '' else concat(" and t0.price_id in (",uId,")")   end;
		SET @s =  (concat ("select price_id,product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,is_default,normal_price,disc_one,disc_two,disc_tpr,price
								,'product_code,product_name,uom_code,ratio,normal_price,disc_one_disc_two,disc_tpr,price,is_default' field_key
								,'Kode,Nama,Uom,Ratio,Harga,Disc 1(%),Disc 2(%),Disc tpr,Harga,Default' field_label
								,'product_code,product_name,uom_code,ratio,normal_price,disc_one_disc_two,disc_tpr,price,is_default' field_export
								,'Kode,Nama,Uom,Ratio,Harga,Disc 1(%),Disc 2(%),Disc tpr,Harga,Default' field_export_label
								,'normal_price,disc_one,disc_two,disc_tpr,price' field_int
								,'' field_footer
								,'' field_level 
							from 
								(select 
									price_id,product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,is_default,t1.normal_price,disc_one*-1 disc_one,disc_two*-1 disc_two,disc_tpr*-1 disc_tpr,
									cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price
								from 
									price_product_uom t0
								left join 
									(select id,product_code,product_name,price normal_price from products) t1 on t1.id = t0.product_id
								left join 
									(select id,uom_code from uoms)t2 on t2.id = t0.uom_id
								where
									deleted_at is null ",uIdSet,"
								) x order by product_id,item_no;"));
	else
		SET @s =  (concat ("select id,effective_date,expired_date,company_id,company_code,company_name,product_id,product_code,product_name,product_division_id,product_division_code,product_type_id,product_type_name,normal_price,disc_one,disc_one_desc,disc_two,disc_tpr,price,uom_id,uom_code,ratio,sure_name,status_id,status_data
								,'effective_date,expired_date,company_code,product_code,product_name,sure_name,product_division_code,product_type_name,normal_price,disc_one_desc,disc_two_desc,disc_tpr_desc,price,status_data' field_key
								,'Tgl efektif,Tgl exp,Pelanggan,Kode,Nama,Alias,Divisi,Tipe,Normal,Disc 1,Disc 2,Disc tpr,Harga,Status Data' field_label
								,'effective_date,expired_date,company_code,product_code,product_name,sure_name,product_division_code,product_type_name,normal_price,disc_one_desc,disc_two_desc,disc_tpr_desc,price,status_data' field_export
								,'Tgl efektif,Tgl exp,Pelanggan,Kode,Nama,Alias,Divisi,Tipe,Normal,Disc 1,Disc 2,Disc tpr,Harga,Status Data' field_export_label
								,'normal_price,price' field_int
								,'' field_footer
								,'' field_level 
							from (
								select t0.id,effective_date,expired_date,company_id,t3.company_code,t3.company_name,product_id,t1.product_code,t1.product_name,t1.product_division_id,t1.product_division_code,t1.product_type_id,t1.product_type_name,t1.normal_price,
									disc_one,concat(disc_one,' %')  disc_one_desc,disc_two,concat(disc_two,' %') disc_two_desc,disc_tpr,concat('Rp',disc_tpr) disc_tpr_desc,
									cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price,
									uom_id,t2.uom_code,ratio,sure_name,status_id,case when date(t0.created_at) = date(t0.updated_at) then 'NEW' else 'EDIT!!' end status_data,
									DENSE_RANK() OVER (PARTITION BY t0.company_id,product_id,price_type order by effective_date desc,ifnull(expired_date,'9999-12-31'), id desc) dr
								from  prices t0
									left join (select t0.id,product_code,product_name,product_division_id,t1.product_division_code,product_type_id,t2.product_type_name,price normal_price from products t0
													left join (select id,division_code product_division_code from product_divisions) t1 on t1.id=t0.product_division_id
													left join (select id,type_name product_type_name from product_types) t2 on t2.id=t0.product_type_id) t1 on t1.id=t0.product_id
									left join (select id,uom_code from uoms) t2 on t2.id=t0.uom_id
									left join (select id,`code` company_code,`name` company_name from companies) t3 on t3.id=t0.company_id
								where deleted_at is null ",theDateSet,uIdSet,userIdSet,statusIdSet,updatedAtSet,divisionIdSet,typeIdSet,keywordSet,"
								) x where dr=1 ",ColumnSet,limitSet," ;"));
	end if;
	PREPARE stmt FROM @s;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
