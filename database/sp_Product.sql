SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_Product;
DELIMITER $$
CREATE PROCEDURE sp_Product(updatedAt date,uId int,divisionIds varchar(7),typeIds varchar(7),purchaseIds varchar(5),salesIds varchar(5),productionIds varchar(5),in statusId int,userId int, reportType int, keyword varchar(255), in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000),in limitVal int, in offsetVal int )
BEGIN
	declare keywordSet varchar(8000);
    declare updatedAtSet   varchar(255);
    declare userIdSet   varchar(255);
    declare uIdSet   varchar(255);
    declare divisionIdSet   varchar(255);
    declare purchaseIdSet   varchar(255);
    declare salesIdSet   varchar(255);
    declare productionIdSet   varchar(255);
    declare typeIdSet   varchar(255);
    declare statusIdSet   varchar(255);
    
    declare limitSet  varchar(255);
    declare offsetSet  varchar(255);
    
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
    
	SET uIdSet  = case when (uId is null  or uId = 0) then '' else concat(" and t0.id in (",uId,")")   end;
    SET statusIdSet  = case when statusId is null  then '' else concat(" and t0.status_id in (",statusId,")")   end;
    SET updatedAtSet  = case when updatedAt is null  then '' else concat(" and (date(t0.created_at) = date('",updatedAt,"')  or date(ifnull(t0.updated_at,t0.created_at)) = date('",updatedAt,"'))")   end;
    set divisionIdSet = case when (divisionIds is null or divisionIds = '') then '' else concat(" AND t0.product_division_id in (",replace(divisionIds,'''',''),")") end ;
    set typeIdSet = case when (typeIds is null or typeIds = '') then '' else concat(" and (t0.product_type_id in (",replace(typeIds,'''',''),")") end ;
    set purchaseIdSet = case when (purchaseIds is null or purchaseIds = '') then '' else concat(" AND t1.is_purchase in (",replace(purchaseIds,'''',''),")") end ;
    set salesIdSet = case when (salesIds is null or salesIds = '') then '' else concat(" AND t1.is_sales in (",replace(salesIds,'''',''),")") end ;
    set productionIdSet = case when (productionIds is null or productionIds = '') then '' else concat(" AND t1.is_production in (",replace(productionIds,'''',''),")") end ;
    SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and (t0.product_code like '%",keyword,"%'  or t0.product_name like '%",keyword,"%' or t0.product_division_code like '%",keyword,"%' or t0.product_division_name like '%",keyword,"%' 
		 or t0.product_type_name like '%",keyword,"%'   or t0.product_uom_code like '%",keyword,"%')")   end;
   
    SET limitSet = case when LimitVal is null then '' else concat(" limit ",limitVal)    end;
    set offsetSet = case when OffsetVal is null then '' else concat(" offset ",offsetVal)    end ;
    
    -- product_uom
    if reportType = 0 then
		SET @s =  (concat ("select product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,is_default
								,'product_code,product_name,uom_code,ratio,is_default' field_key
								,'Kode,Nama,Uom,Ratio,Default' field_label
								,'product_code,product_name,uom_code,ratio,is_default' field_export
								,'Kode,Nama,Uom,Ratio,Default' field_export_label
								,'ratio' field_int
								,'' field_footer
								,'' field_level 
							from 
								(select 
									product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,is_default 
								from 
									product_uom t0
								left join 
									(select id,product_code,product_name from products) t1 on t1.id = t0.product_id
								left join 
									(select id,uom_code from uoms)t2 on t2.id = t0.uom_id
                                where
									product_id in (select t0.id from products t0
													left join
														(select id,is_purchase,is_sales,is_production from product_types) t1 on t1.id = t0.product_type_id
													where deleted_at is null ",uIdSet,updatedAtSet,keywordSet,divisionIdSet,typeIdSet,statusIdSet,purchaseIdSet,salesIdSet,productionIdSet,")
								) x order by product_id,item_no;"));
    else
		SET @s =  (concat ("select id,product_code,product_name,serial_number,product_type_id,product_type_name,product_division_id,product_division_code,product_division_name,uom_id,uom_code,lead_time,status_id,status_data,is_purchase,is_sales,is_production
								,'product_code,product_name,serial_number,lead_time,uom_code,product_type_name,division_code,is_purchase,is_sales,is_production,status_id,status_data' field_key
                                ,'Kode,Nama,Serial,ETA,Uom,Tipe,Divisi,Pembelian,Sales,Produksi,Status,Status Data' field_label
                                ,'product_code,product_name,serial_number,lead_time,uom_code,product_type_name,division_code,is_purchase,is_sales,is_production,status_id,status_data' field_export
                                ,'Kode,Nama,Serial,ETA,Uom,Tipe,Divisi,Pembelian,Sales,Produksi,Status,Status Data' field_export_label
                                ,'lead_time' field_int
                                ,'' field_footer
                                ,'product_id' field_level
							from
								(select t0.id,product_code,product_name,serial_number,product_type_id,product_type_name,product_division_id,product_division_code,product_division_name,uom_id,uom_code,lead_time,status_id,case when date(created_at) = date(ifnull(updated_at,created_at)) then 'NEW' else 'EDIT!!' end status_data
									,is_purchase,is_sales,is_production
                                from 
									products t0
								left join
									(select id,is_purchase,is_sales,is_production from product_types) t1 on t1.id = t0.product_type_id
								where deleted_at is null ",uIdSet,updatedAtSet,keywordSet,divisionIdSet,typeIdSet,statusIdSet,purchaseIdSet,salesIdSet,productionIdSet,") x where id != 0 ",ColumnSet,limitSet,offsetSet," ;"));
	end  if;
	
	PREPARE stmt FROM @s;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
