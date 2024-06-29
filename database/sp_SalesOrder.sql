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
   SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and (t0.reference_no like '%",keyword,"%'  or t1.pool_name like '%",keyword,"%' or t2.outlet_name like '%",keyword,"%' or t3.customer_code like '%",keyword,"%' 
		 or t3.customer_name like '%",keyword,"%'   or t4.plant_name like '%",keyword,"%' or t0.delivery_address like '%",keyword,"%' or t5.employee_name like '%",keyword,"%' or t0.status_description like '%",keyword,"%')")   end;

    SET limitSet = case when LimitVal is null then '' else concat(" limit ",offsetVal,",",limitVal)    end;
	if reportTypeId = 0 then
		set uIdSet = case when (uId is null  or uId = 0) then '' else concat(" and t0.price_id in (",uId,")")   end;
		SET @s =  (concat ("select id,sales_order_id,reference_no,item_no,product_id,product_code,product_name,qty,uom_id,uom_code,ratio,packaging_id,packaging_code,final_qty,final_uom_id,final_uom_code,normal_price,price_id,price,disc1,disc1_amount,disc2,disc2_amount,disc_tpr,total_disc
								,nett_price,total,conversion_qty,conversion_uom_id,conversion_uom_code,status_data 
                                ,'product_code,product_name,qty,uom_code,ratio,packaging_code,final_qty,final_uom_code,price,disc1,disc2,disc_tpr,nett_price,total,status_data' field_key
                                ,'Kode,Nama,Qty,Sat,@,Sat,Total qty,Sat,Harga,Disc 1(%),Disc 2(%),Disc (Rp),Harga nett,Jumlah,Status Data' field_label
                                ,'product_code,product_name,qty,uom_code,ratio,packaging_code,final_qty,final_uom_code,price,disc1,disc2,disc_tpr,nett_price,total,status_data' field_export
                                ,'Kode,Nama,Qty,Sat,@,Sat,Total qty,Sat,Harga,Disc 1(%),Disc 2(%),Disc (Rp),Harga nett,Jumlah,Status Data' field_export_label
                                ,'qty,ratio,final_qty,price,disc1,disc2,disc_tpr,nett_price,total' field_int
                                ,'qty,final_qty,total' field_footer
                                ,'' field_level
							from (
								select t0.id,sales_order_id,reference_no,item_no,product_id,t1.product_code,t1.product_name,qty,uom_id,uom_code,ratio,packaging_id,packaging_code,final_qty,final_uom_id,final_uom_code,normal_price,price_id,price,disc1,disc1_amount,disc2,disc2_amount,disc_tpr,total_disc
									,nett_price,total,conversion_qty,conversion_uom_id,conversion_uom_code,case when date(t0.created_at) = date(t0.updated_at) then 'NEW' else 'EDIT!!' end status_data 
								from sales_order_detail t0
									left join (select id,product_code,product_name from products) t1 on t1.id = t0.product_id
								where 
									deleted_at is null
								and sales_order_id in 
									(select t0.id from sales_order t0
										left join (select id,`name` pool_name from pools) t1 on t1.id = t0.pool_id
										left join (select id,`name` outlet_name from plants) t2 on t2.id = t0.outlet_id
										left join (select id,`code` customer_code,`name` customer_name from companies) t3 on t3.id = t0.customer_id
										left join (select id,`name` plant_name from plants) t4 on t4.id = t0.plant_id
										left join (select id, employee_name from employees) t5 on t5.id = t0.employee_id
									where deleted_at is null ",theDateSet,dueDateSet,updatedAtSet,uIdSet,userIdSet,statusIdSet,employeeIdsSet,outletIdsSet,customerIdsSet,plantIdSet,productIdsSet,keywordSet,")) x ;"));
	else
		SET @s =  (concat ("select id,id sales_order_id,reference_no,issue_date,due_date,lead_time,pool_id,pool_name,outlet_id,outlet_name,customer_id,customer_code,customer_name,plant_id,plant_name,full_name,terms,delivery_address,employee_id,employee_name
								,subtotal,total_disc,dpp,ppn,ppn_amount,total
                                ,status_id,status_description,status_data
								,'reference_no,issue_date,due_date,lead_time,outlet_name,full_name,terms,delivery_address,employee_name,subtotal,total_disc,dpp,ppn_amount,total,status_description,status_data' field_key
								,'No SO,Tgl,Tgl akhir,Berlaku,Outlet,Customer,Alamat pengiriman,Salesman,Subtotal,Diskon,Dpp,Ppn,Total,Status,Status Data' field_label
								,'reference_no,issue_date,due_date,lead_time,outlet_name,full_name,terms,delivery_address,employee_name,subtotal,total_disc,dpp,ppn_amount,total,status_description,status_data' field_export
								,'No SO,Tgl,Tgl akhir,Berlaku,Outlet,Customer,Alamat pengiriman,Salesman,Subtotal,Diskon,Dpp,Ppn,Total,Status,Status Data' field_export_label
								,'terms,subtotal,total_disc,dpp,ppn_amount,total' field_int
								,'subtotal,total_disc,dpp,ppn_amount,total' field_footer
								,'sales_order_id' field_level 
							from (
								select t0.id,reference_no,issue_date,due_date,concat(lead_time,' hari') lead_time,pool_id,t1.pool_name,outlet_id,t2.outlet_name,customer_id,t3.customer_code,t3.customer_name,plant_id,t4.plant_name,case when plant_id = 0 then t3.customer_code else concat(t3.customer_code,' -- ',t4.plant_name) end full_name,terms,delivery_address,
									employee_id,t5.employee_name,subtotal,total_disc,dpp,ppn,ppn_amount,total,
                                    status_id,status_description,
									case when date(t0.created_at) = date(t0.updated_at) then 'NEW' else 'EDIT!!' end status_data
								from  sales_order t0
									left join (select id,`name` pool_name from pools) t1 on t1.id = t0.pool_id
                                    left join (select id,`name` outlet_name from plants) t2 on t2.id = t0.outlet_id
                                    left join (select id,`code` customer_code,`name` customer_name from companies) t3 on t3.id = t0.customer_id
                                    left join (select id,`name` plant_name from plants) t4 on t4.id = t0.plant_id
                                    left join (select id,employee_name from employees) t5 on t5.id = t0.employee_id
								where deleted_at is null ",theDateSet,dueDateSet,updatedAtSet,uIdSet,userIdSet,statusIdSet,employeeIdsSet,outletIdsSet,customerIdsSet,plantIdSet,productIdsSet,keywordSet,"
								) x where id != 0",ColumnSet,limitSet," ;"));
	end if;
	PREPARE stmt FROM @s;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
