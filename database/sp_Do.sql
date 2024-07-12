SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_Do;
DELIMITER $$
CREATE PROCEDURE sp_Do(theDate date,updatedAt date,uId int,warehouseIds varchar(75),warehousePlantIds varchar(75),outletIds varchar(75),customerIds varchar(75),plantIds varchar(75),productIds varchar(75),statusIds varchar(75),reportGrupId int,reportTypeId int,userId int,keyword varchar(255),in searchDetail int,in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000), in limitVal int, in offsetVal int ) 
BEGIN
	declare keywordSet varchar(8000);
    declare theDateSet varchar(100);
    declare updatedAtSet varchar(100);
    declare userIdSet   varchar(255);
    declare uIdSet   varchar(255);
    declare warehouseIdsSet   varchar(255);
    declare warehousePlantIdsSet   varchar(255);
	declare outletIdsSet   varchar(255);
	declare customerIdsSet   varchar(255);
	declare plantIdsSet   varchar(255);
	declare productIdsSet   varchar(255);
    declare statusIdsSet   varchar(255);
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
    
    set theDateSet = case when theDate is null then '' else concat("  and issue_date >= '",TheDate,"' and issue_date <= '",TheDate,"' ")   end;
    SET updatedAtSet  = case when updatedAt is null  then '' else concat(" and (date(t0.created_at) = date('",updatedAt,"')  or date(ifnull(t0.updated_at,t0.created_at)) = date('",updatedAt,"'))")   end;
    SET userIdSet = case when ( UserId is null or UserId = 0 ) then '' else concat( " and (t0.warehouse_plant_id in (select plant_id from sys_user_plant where user_id = ", userId, ") or t0.warehouse_id in (select company_id from sys_user_plant where user_id = ", userId, "))" ) end;
	SET uIdSet  = case when (uId is null  or uId = 0) then '' else concat(" and t0.id in (",uId,")")   end;
    SET warehouseIdsSet  = case when (warehouseIds is null or warehouseIds = '') then '' else concat(" and t0.warehouse_id in (",replace(warehouseIds,'''',''),")") end ;
    SET warehousePlantIdsSet  = case when (warehousePlantIds is null or warehousePlantIds = '') then '' else concat(" and t0.warehouse_plant_id in (",replace(warehousePlantIds,'''',''),")") end ;
    SET statusIdSet  = case when (statusIds is null or statusIds = '') then '' else concat(" and t0.status_id in (",replace(statusIds,'''',''),")") end ;
    set outletIdsSet = case when (outletIds is null or outletIds = '') then '' else concat(" and (t0.outlet_id in (",replace(outletIds,'''',''),")") end ;
    set customerIdsSet = case when (customerIds is null or customerIds = '') then '' else concat(" and (t0.customer_id in (",replace(customerIds,'''',''),")") end ;
    set plantIdsSet = case when (plantIds is null or plantIds = '') then '' else concat(" and t0.plant_id in (",replace(plantIds,'''',''),")") end ;
    set productIdsSet = case when (productIds is null or productIds = '') then '' else concat(" and t0.id in (select product_id from do_detail where deleted_at is null and t0.product_id in (",replace(productIds,'''',''),"))") end ;
    
	if searchDetail = 1 then
		SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and t0.id in (select do_id from do_detail t0 left join (select id,product_code,product_name,art_no,barcode from products) t1 on t1.id = t0.product_id 
			where deleted_at is null and (t1.product_code like '%",keyword,"%' or t1.product_name like '%",keyword,"%' or t1.art_no like '%",keyword,"%' or t1.barcode like '%",keyword,"%'))")   end;
    else
		SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and (t0.sales_order_no like '%",keyword,"%'  or t0.reference_no like '%",keyword,"%' or t1.warehouse_code like '%",keyword,"%' or t1.warehouse_name like '%",keyword,"%' or t2.warehouse_plant_code like '%",keyword,"%' 
			 or t2.warehouse_plant_name like '%",keyword,"%'   or t3.customer_code like '%",keyword,"%' or t3.customer_name like '%",keyword,"%' or t4.plant_code like '%",keyword,"%' or t4.plant_name like '%",keyword,"%' or t0.delivery_address like '%",keyword,"%' or t0.status_description like '%",keyword,"%'
              or t5.transporter_code like '%",keyword,"%' or t5.transporter_name like '%",keyword,"%' or t0.plate_no like '%",keyword,"%'  or t0.notes like '%",keyword,"%')")   end;
	end if;
    SET limitSet = case when LimitVal is null then '' else concat(" limit ",offsetVal,",",limitVal)    end;
    
    if reportGrupId = 1 then
		SET @s =  (concat ("select sum(do_count) all_dos,sum(case when status_id = 1 then do_count else 0 end) do_open,sum(case when status_id = 2 then do_count else 0 end) do_confirm,sum(case when status_id = 3 then do_count else 0 end) do_shipping,sum(case when status_id = 4 then do_count else 0 end) do_complete,
								sum(case when status_id = 99 then do_count else 0 end) do_cancel,
								'all_dos,do_open,do_confirm,do_shipping,do_complete,do_cancel' field_key,
								'Total,Open,Proses,Kirim,Selesai,Cancel' field_label,
								'' field_export,
								'' field_export_label,
								'all_dos,do_open,do_confirm,do_shipping,do_complete,do_cancel' field_int,
								'' field_footer,
								'' field_level
							from (
								select 
									t0.id,sales_order_id,sales_order_no,reference_no,issue_date,warehouse_id,t1.warehouse_code,t1.warehouse_name,t0.warehouse_plant_id,t2.warehouse_plant_code,t2.warehouse_plant_name,customer_id,t3.customer_code,t3.customer_name,t0.plant_id,t4.plant_code,t4.plant_name,
									delivery_address,transporter_id,t5.transporter_code,t5.transporter_name,courier_id,courier_name,plate_no,notes,status_id,status_description,1 do_count
								from dos t0
									left join (select id,`code` warehouse_code,`name` warehouse_name from companies) t1 on t1.id = t0.warehouse_id
									left join (select id,`code` warehouse_plant_code,`name` warehouse_plant_name from plants) t2 on t2.id = t0.warehouse_plant_id
									left join (select id,`code` customer_code,`name` customer_name from companies) t3 on t3.id = t0.customer_id
									left join (select id,`code` plant_code,`name` plant_name from plants) t4 on t4.id = t0.plant_id
									left join (select id,`code` transporter_code,`name` transporter_name from companies) t5 on t5.id = t0.transporter_id
								where 
									deleted_at is null",userIdSet,warehouseIdsSet,warehousePlantIdsSet,statusIdSet,outletIdsSet,customerIdsSet,plantIdsSet,keywordSet,") x ;"));
		
    else
		
    end if;
	PREPARE stmt FROM @s;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
