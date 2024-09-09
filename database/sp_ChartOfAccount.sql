SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_ChartOfAccount;
DELIMITER $$

CREATE PROCEDURE sp_ChartOfAccount(TheDate datetime,updatedAt datetime,parentLevelNo int,levelNo int,accountTypeId int,companyId int,salesTypeId int,userId int,keyword varchar(255),in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000),in limitVal int, in offsetVal int ) 
BEGIN
	
	declare keywordSet varchar(8000);
    declare theDateSet varchar(255);
    declare updatedAtSet varchar(255);
    declare limitSet  varchar(255);
    declare offsetSet  varchar(255);
    declare levelNoSet  varchar(255);
    declare parentLevelNoSet  varchar(255);
    declare accountTypeIdSet varchar(255);
    declare companyIdSet varchar(255);
    declare salesTypeIdSet varchar(255);
    declare autId int default 0;
	declare menuId int default 0;
	declare userIdSet  varchar(255) default '';
    
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
    
    SET theDateSet  = case when (TheDate is null or TheDate = 0) then '' else concat(" and t0.effective_date <= date('",TheDate,"')  and ifnull(t0.EXPIRED_DATE,date('",TheDate,"') ) >= date('",TheDate,"')")   end;
    SET updatedAtSet  = case when (updatedAt is null or updatedAt = 0) then '' else concat(" and (date(t0.created_at) = date('",updatedAt,"')  or date(t0.updated_at) = date('",updatedAt,"'))")   end;
    SET keywordSet  = case when (keyword is null or keyword = '') then '' else concat(" AND (t0.code_coa like '%",keyword,"%' or t0.name_coa like '%",keyword,"%' or t0.code_out like '%",keyword,"%' or t0.code_in like '%",keyword,"%')")   end;
	SET limitSet = case when LimitVal is null then '' else concat(" limit ",limitVal)    end;
    set offsetSet = case when OffsetVal is null then '' else concat(" offset ",offsetVal)    end ;
    SET parentLevelNoSet  = case when parentLevelNo =0 then '' else concat(" AND t0.level_no = ",parentLevelNo)   end;
    SET levelNoSet  = case when levelNo =0 then '' else concat(" AND t0.level_no = ",levelNo)   end;
    SET accountTypeIdSet  = case when (accountTypeId is null or accountTypeId = 0) then '' else concat(" AND t0.account_type_id = ",accountTypeId)   end;
    SET companyIdSet  = case when (companyId is null or companyId = 0) then '' else concat(" AND t0.company_id = ",companyId)   end;
    SET salesTypeIdSet  = '';
    -- SET accountTypeIdSet  = case when (accountTypeId is null or accountTypeId = 0) then '' else  concat(" AND t0.account_type_id in (select id from gl_account_type where component_account = '",accountTypeName,"')")   end;
    -- set menuId = (select id from sys_menus where form_name = 'chart_of_accounts');
--     set autId = ifnull((select permission_id from sys_role_menu_permission where  menu_id = menuId and permission_id = 6 and role_id in (select role_id from sys_user_role where user_id = userId) limit 1),0);
--     if autId = 0 then
-- 		SET userIdSet  = concat(" and t0.id in (select account_id from sys_user_account where user_id = ",userId,")")  ;
--     end if;
    
    
	
	SET @s =  (concat ("select t0.id,effective_date,expired_date,account_type_id,account_type_name,company_id,company_code,company_name,sales_type_id,sales_type_name,level_no,parent_id,t1.parent_code,t1.parent_name,code_coa,name_coa,code_out,code_in,journal_position,status_id,
							case when date(t0.created_at) = date(t0.updated_at) then 'NEW' else 'EDIT!!' end status_data
						from chart_of_accounts t0 
							left join (select id,code_coa parent_code,name_coa parent_name from chart_of_accounts) t1 on t1.id = t0.parent_id
						where deleted_at is null ",theDateSet,keywordSet,ColumnSet,levelNoSet,accountTypeIdSet,parentLevelNoSet,companyIdSet,salesTypeIdSet,updatedAtSet,userIdSet,"
							order by t0.id desc ",limitSet,offsetSet,"; "));
	
                        
    PREPARE stmt FROM @s;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

END$$

DELIMITER ;