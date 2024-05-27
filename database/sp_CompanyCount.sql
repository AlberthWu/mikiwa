SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_CompanyCount;
DELIMITER $$
CREATE PROCEDURE sp_CompanyCount(updatedAt date,uId int,companyTypeId int,statusIds varchar(5),userId int, keyword varchar(255), in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000),in limitVal int, in offsetVal int )
BEGIN
	declare keywordSet varchar(8000);
    declare updatedAtSet   varchar(255);
    declare userIdSet   varchar(255);
    declare uIdSet   varchar(255);
    declare companyTypeIdSet   varchar(255);
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
    SET updatedAtSet  = case when updatedAt is null  then '' else concat(" and (date(t0.created_at) = date('",updatedAt,"')  or date(ifnull(t0.updated_at,t0.created_at)) = date('",updatedAt,"'))")   end;
    SET statusIdSet  = case when (statusIds is null or statusIds = '') then '' else concat(" and t0.status_id in (",replace(statusIds,'''',''),")") end ;
    set companyTypeIdSet = case when (companyTypeId is null or companyTypeId = '') then '' else concat(" AND t0.id in (select company_id from company_type where type_id in (",companyTypeId,"))") end ;
    SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" and (t0.code like '%",keyword,"%'  or t0.name like '%",keyword,"%' or t0.email like '%",keyword,"%' or t0.phone like '%",keyword,"%' 
		 or t0.address like '%",keyword,"%'  or t0.npwp like '%",keyword,"%' or t0.bank_no like '%",keyword,"%')")   end;
   
    SET limitSet = case when LimitVal is null then '' else concat(" limit ",limitVal)    end;
    set offsetSet = case when OffsetVal is null then '' else concat(" offset ",offsetVal)    end ;
    
    
	SET @s =  (concat ("select count(id) rn from (select * from (
						select t0.id,t0.parent_id,code,name,email,phone,fax,npwp,npwp_name,npwp_address,terms,credit,address,concat(t1.city_name,', ',t2.district_name,', ', t3.state_name) teritory,zip,city_id,t1.city_name,bank_id,t4.bank_name,bank_no,bank_account_name,bank_branch,is_cash,is_po,is_tax,is_receipt,status
						from
							companies t0
						left join
							(select id,name city_name,parent_id from cities) t1 on t1.id=t0.city_id
						Left Join
							(select id,name district_name,parent_id from cities) t2 On t2.id=t1.parent_id	
						Left Join
							(select id,name state_name from cities) t3 On t3.id=t2.parent_id
						left join
							(select id,name bank_name from banks) t4 on t4.id=t0.bank_id
						where deleted_at is null ",uIdSet,updatedAtSet,keywordSet,companyTypeIdSet,statusIdSet,") x where id != 0 ",ColumnSet,") y ;"));
	
	PREPARE stmt FROM @s;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
