SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_PettyCashV2;

DELIMITER $$
CREATE PROCEDURE sp_PettyCashV2(TheDate datetime,TheDate2 datetime,in updatedAt datetime,voucherId int,accountId int,transactionOnly int,receiptStatus varchar(75),userId int,reportTypeId int,keyword varchar(255),in searchDetail int,in TheString varchar(8000),in TheParameter varchar(8000),in TheKeyword varchar(8000),in limitVal int, in offsetVal int) 
BEGIN
	declare theDateSet varchar(255);
    declare updatedAtSet   varchar(255);
	declare voucherIdSet   varchar(255);
	declare uidSet   varchar(255) default '';
    declare accountIdSet   varchar(255);
    declare userIdSet   varchar(255) default '';
    declare autId int default 0;
	declare menuId int default 0;
    declare limitSet  varchar(255);
    declare offsetSet  varchar(255); 
    declare receiptStatusSet varchar(255); 
    declare keywordSet varchar(8000);             
    declare keywordSetDetail varchar(8000);
    declare transactionOnlySet  varchar(255); 
    
    declare periodSet1 int;
    declare periodSet2 int;
    
    DECLARE Rn int;
    Declare Urut int DEFAULT 1;
    Declare ColumnName varchar(255) ;
    Declare ParameterName varchar(255);
    Declare KeywordName varchar(255);
    DECLARE ColumnSet varchar(8000) default '';

    set Rn = (CHAR_LENGTH(TheString) - CHAR_LENGTH(REPLACE(TheString, ',','')) + 1) ;
    
     WHILE Urut <= Rn DO
		Set ColumnName = SUBSTRING_INDEX(SUBSTRING_INDEX(TheString, ',', Urut),',',-1);
        Set ParameterName = SUBSTRING_INDEX(SUBSTRING_INDEX(TheParameter, ',', Urut),',',-1);
        Set KeywordName = SUBSTRING_INDEX(SUBSTRING_INDEX(TheKeyword, ',', Urut),',',-1);
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
    
    IF receiptStatus=  'Open' THEN
		SET receiptStatusSet  = concat(" AND t4.status_id = 0") ;
    ELSEIF  replace(receiptStatus," ","") = 'Approved' THEN
		SET receiptStatusSet  = concat(" AND t4.status_id = 1") ;
    ELSEIF  replace(receiptStatus," ","") = 'Posted' THEN
		SET receiptStatusSet  = concat(" AND t4.status_gl_id = 1") ;
    ELSE
		SET receiptStatusSet  = '';
    end if;
    set periodSet1 = concat(year(TheDate),LPAD(convert(month(TheDate),char(2)),2,0),LPAD(convert(day(TheDate),char(4)),2,0));
	set periodSet2 = concat(year(TheDate2),LPAD(convert(month(TheDate2),char(2)),2,0),LPAD(convert(day(TheDate2),char(4)),2,0));
    SET theDateSet  = case when (TheDate is null or TheDate = 0) then '' else concat(" AND t0.period between ",periodSet1," and ",periodSet2)   end;
    SET updatedAtSet  = case when updatedAt is null  then '' else concat(" and (date(t4.created_at) = date('",updatedAt,"')  or date(ifnull(t4.updated_at,t0.created_at)) = date('",updatedAt,"'))")   end;
	SET voucherIdSet  = case when voucherId is null  then '' else concat(" AND t0.voucher_id = ",voucherId)   end;
    SET uidSet  = case when voucherId is null  then '' else concat(" AND t0.id = ",voucherId)   end;
    SET accountIdSet  = case when (accountId is null  or accountId = 0) then '' else concat(" AND t0.id = ",accountId)   end;
    set transactionOnlySet = case when transactionOnly = 0 then '' else concat(" and (debet is not null or credit is not null)")  end;
    set menuId = (select id from sys_menus where form_name = 'petty_cash');
    set autId = ifnull((select permission_id from sys_role_menu_permission where  menu_id = menuId and permission_id = 6 and role_id in (select role_id from sys_user_role where user_id = userId) limit 1),0);
    if autId = 0 then
		SET userIdSet  = concat(" and t0.id in (select account_id from sys_user_account where user_id = ",userId,")")  ;
    end if;
    
	SET keywordSet  = case when (keyword is null   or keyword = '') then '' else concat(" AND (t0.memo like '%",keyword,"%' or t0.voucher_no like '%",keyword,"%' or t0.voucher_code like '%",keyword,"%' or t1.account_code_header like '%",keyword,"%' or t1.account_name_header like '%",keyword,"%' or t2.account_code like '%",keyword,"%' or t2.account_name like '%",keyword,"%')")   end;
    
    SET limitSet = case when LimitVal is null then '' else concat(" limit ",offsetVal,",",limitVal)    end;

    if reportTypeId = 0 then
		SET accountIdSet  = case when (accountId is null  or accountId = 0) then '' else concat(" AND t0.account_id_header = ",accountId)   end;
        if autId = 0 then
			SET userIdSet  = concat(" and t0.account_id_header in (select account_id from sys_user_account where user_id = ",userId,")")  ;
		end if;
		if voucherId= 0 then
			SET @s =  (concat ("select null issue_date,null voucher_no,null account_code,null account_name,null memo,null debet,null credit,null balance,null pic,null transaction_type,null status_id,null status_gl_id,null updated_by,null status_data,
									'issue_date,voucher_no,account_code,account_name,memo,debet,credit,balance,pic,transaction_type,status_id,status_gl_id,updated_by,status_data' field_key,
                                    'Tgl,No voucher,Kode,Nama,Memo,Debet,Kredit,Akhir,Pic,Transaksi,Status,Status gl,User,Status data' field_label,
                                    'issue_date,voucher_code,voucher_no,account_code_header,account_name_header,account_code,account_name,memo,debet,credit,balance,pic,transaction_type,status_id,status_gl_id,updated_by,status_data' field_export,
                                    'Tgl,Kode voucher,No voucher,Kode Akun, Nama Akum,Kode,Nama,Memo,Debet,Kredit,Akhir,Pic,Transaksi,Status,Status gl,User,Status data' field_export_label,
                                    'debet,credit,balance' field_int,
                                    'debet,credit,balance' field_footer,
                                    '' field_level;"));
        else
			select uidSet,theDateSet,receiptStatusSet,accountIdSet,userIdSet,keywordSet,voucherIdSet;
			SET @s =  (concat ("select id,issue_date,voucher_id,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,item_no,account_id,account_code,account_name,memo,debet,credit,balance,pic,transaction_type,status_id,status_gl_id,updated_by,status_data,
									'issue_date,voucher_no,account_code,account_name,memo,debet,credit,balance,pic,transaction_type,status_id,status_gl_id,updated_by,status_data' field_key,
                                    'Tgl,No voucher,Kode,Nama,Memo,Debet,Kredit,Akhir,Pic,Transaksi,Status,Status gl,User,Status data' field_label,
                                    'issue_date,voucher_code,voucher_no,account_code_header,account_name_header,account_code,account_name,memo,debet,credit,balance,pic,transaction_type,status_id,status_gl_id,updated_by,status_data' field_export,
                                    'Tgl,Kode voucher,No voucher,Kode Akun, Nama Akum,Kode,Nama,Memo,Debet,Kredit,Akhir,Pic,Transaksi,Status,Status gl,User,Status data' field_export_label,
                                    'debet,credit,balance' field_int,
                                    'debet,credit,balance' field_footer,
                                    '' field_level
                                from (
									select urut,id,voucher_id,issue_date,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,memo,debet,credit,item_no,account_id,account_code,account_name,
										sum((balance+debet)-credit) over (partition by account_id_header order by account_id_header,period,voucher_seq_no,id) balance,pic,transaction_type,status_id,status_desc,status_gl_id,status_gl_desc,updated_by,status_data
									from 
										(	select 0 urut,0 id,0 voucher_id,t0.account_id_header,t1.account_code_header,t1.account_name_header,date('",TheDate,"') issue_date,0 voucher_seq_no,'Opening' voucher_no,0 item_no,0 account_id,'' account_code,'' account_name,0 debet,0 credit, sum(debet)-sum(credit) balance,
												",periodSet1,"   period,'In' transaction_type,1 status_id, '' status_desc,1 status_gl_id,'' status_gl_desc,'' voucher_code,'' memo,'' pic,'' status_data,'' updated_by
											from petty_cash t0
												left join (select id,name_coa account_name_header,code_coa account_code_header from chart_of_accounts) t1 on t1.id = t0.account_id_header
												left join (select id,name_coa account_name,code_coa account_code from chart_of_accounts) t2 on t2.id = t0.account_id
                                                left join (select max(id) id,account_id from petty_cash_header t0 where deleted_at is null ",theDateSet,uidSet," group by account_id) t3 on t3.account_id = t0.account_id_header
											where t0.deleted_at is null  and t3.id is not null  and period < ",periodSet1,accountIdSet,userIdSet,"
												group  by t0.account_id_header,t1.account_code_header,t1.account_name_header
											union all
											select 1 urut,t0.id,t0.voucher_id,t0.account_id_header,t1.account_code_header,t1.account_name_header,t0.issue_date,t0.voucher_seq_no,t0.voucher_no,t0.item_no,t0.account_id,t2.account_code,t2.account_name,debet,credit, 0 balance,period,transaction_type,
												t4.status_id,t4.status_desc,t4.status_gl_id,t4.status_gl_desc,voucher_code,t0.memo,t0.pic,
												case when date(t0.created_at) = date(t0.updated_at) then 'NEW' else 'EDIT!!' end status_data,ifnull(t0.updated_by,t0.created_by) updated_by
											from petty_cash t0
												left join (select id,name_coa account_name_header,code_coa account_code_header from chart_of_accounts) t1 on t1.id = t0.account_id_header
												left join (select id,name_coa account_name,code_coa account_code from chart_of_accounts) t2 on t2.id = t0.account_id
												left join (select id,issue_date,status_id,case when status_id = 1 then 'Approved' else 'Open' end status_desc,status_gl_id,case when status_gl_id = 1 then 'Posted' else 'Open' end status_gl_desc,created_at,updated_at,created_by,updated_by from petty_cash_header where deleted_at is null) t4 on t4.id=t0.voucher_id
											where t0.deleted_at is null  ",theDateSet,receiptStatusSet,accountIdSet,userIdSet,keywordSet,voucherIdSet,"
										) x where account_id_header != 0 
										order by account_id_header,urut,period,voucher_seq_no,voucher_id,item_no ) y where account_id_header !=0 ",ColumnSet,limitSet,";"));
		end if;
	elseif reportTypeId = 1 then
		-- select periodSet1,theDateSet,receiptStatusSet,keywordSet,accountIdSet,userIdSet,transactionOnlySet,ColumnSet,limitSet;
        
		SET @s =  (concat ("select id,code_coa,name_coa,opening,debet,credit,balance,
								'code_coa,name_coa,opening,debet,credit,balance' field_key,
                                'Kode,Nama,Awal, Debet, Kredit, Akhir' field_label,
                                'code_coa,name_coa,opening,debet,credit,balance' field_export,
                                'Kode,Nama,Awal, Debet, Kredit, Akhir' field_export_label,
                                'opening,debet,credit,balance' field_int,
                                'opening,debet,credit,balance' field_footer,
                                '' field_level
							from 
								(select id,code_coa,name_coa,t1.opening,t1.debet,t1.credit,t1.balance
									from 
										chart_of_accounts t0
									left join
										(select t0.account_id_header account_id,t3.opening,sum(debet) debet,sum(credit) credit, ifnull(t3.opening,0) + sum(debet)-sum(credit) balance 
											from petty_cash t0 
												left join (select id,name_coa account_name_header,code_coa account_code_header from chart_of_accounts) t1 on t1.id = t0.account_id_header
												left join (select id,name_coa account_name,code_coa account_code from chart_of_accounts) t2 on t2.id = t0.account_id
												left join (select account_id_header account_id,sum(debet)-sum(credit) opening from petty_cash where deleted_at is null  and period < ",periodSet1,"  group by account_id_header) t3 on t3.account_id = t0.account_id_header
												left join (select id,status_id,case when status_id = 1 then 'Approved' else 'Open' end status_desc,case when status_gl_id = 1 then 'Posted' else 'Open' end status_gl_desc from petty_cash_header where deleted_at is null) t4 on t4.id = t0.voucher_id
											where 
												deleted_at is null ",theDateSet,receiptStatusSet,keywordSet,voucherIdSet,"
											group by t0.account_id_header,t3.opening) t1 on t1.account_id = t0.id
									where deleted_at is null and status_id = 1 and id not in (select parent_id from chart_of_accounts)  and is_header = 1 ",accountIdSet,userIdSet,transactionOnlySet,") x where id != 0 ",ColumnSet,limitSet,";"));
	elseif reportTypeId = 2 then
		SET accountIdSet  = case when (accountId is null  or accountId = 0) then '' else concat(" AND t0.account_id_header = ",accountId)   end;
        if autId = 0 then
			SET userIdSet  = concat(" and t0.account_id_header in (select account_id from sys_user_account where user_id = ",userId,")")  ;
		end if;
		SET @s =  (concat ("select urut,id,issue_date,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,debet,credit,balance,transaction_type,status_id,status_desc,status_gl_id,status_gl_desc,updated_by,status_data ,
								'issue_date,voucher_code,voucher_no,account_code_header,account_name_header,debet,credit,balance,transaction_type,status_id,status_gl_id,updated_by,status_data' field_key,
                                'Tgl,Kode,Nama,Kode voucher,No voucher,Debet,Kredit,Pic,Memo,Transaksi,User,Status, Status gl,Status data' field_label,
                                'issue_date,voucher_code,voucher_no,account_code_header,account_name_header,debet,credit,balance,transaction_type,status_id,status_gl_id,updated_by,status_data' field_export,
                                'Tgl,Kode,Nama,Kode voucher,No voucher,Debet,Kredit,Pic,Memo,Transaksi,User,Status, Status gl,Status data' field_export_label,
                                'debet,credit,balance' field_int,
                                'debet,credit,balance' field_footer,
                                '' field_level 
							from (
								select urut,id,issue_date,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,debet,credit,
									sum((balance+debet)-credit) over (partition by account_id_header order by account_id_header,period,voucher_seq_no,id) balance,transaction_type,status_id,status_desc,status_gl_id,status_gl_desc,updated_by,status_data
								from 
									(	select 0 urut,0 id,t0.account_id_header,t1.account_code_header,t1.account_name_header,date('",TheDate,"') issue_date,0 voucher_seq_no,'Opening' voucher_no,0 debet,0 credit, sum(debet)-sum(credit) balance,
											",periodSet1,"  period,'In' transaction_type,1 status_id, '' status_desc,1 status_gl_id,'' status_gl_desc,'' voucher_code,'' status_data,'' updated_by
										from petty_cash t0
											left join (select id,name_coa account_name_header,code_coa account_code_header from chart_of_accounts) t1 on t1.id = t0.account_id_header
											left join (select id,name_coa account_name,code_coa account_code from chart_of_accounts) t2 on t2.id = t0.account_id
                                            left join (select max(id) id,account_id from petty_cash_header t0 where deleted_at is null ",theDateSet,uidSet," group by account_id) t3 on t3.account_id = t0.account_id_header
										where t0.deleted_at is null  and t3.id is not null and period < ",periodSet1,accountIdSet,userIdSet,"
											group  by t0.account_id_header,t1.account_code_header,t1.account_name_header
										union all
										select 1 urut,t0.voucher_id id,t0.account_id_header,t1.account_code_header,t1.account_name_header,t0.issue_date,t0.voucher_seq_no,t0.voucher_no,sum(debet) debet,sum(credit) credit, 0 balance,period,transaction_type,
											t4.status_id,t4.status_desc,t4.status_gl_id,t4.status_gl_desc,voucher_code,
											case when date(t4.created_at) = date(t4.updated_at) then 'NEW' else 'EDIT!!' end status_data,ifnull(t4.updated_by,t4.created_by) updated_by
										from petty_cash t0
											left join (select id,name_coa account_name_header,code_coa account_code_header from chart_of_accounts) t1 on t1.id = t0.account_id_header
											left join (select id,name_coa account_name,code_coa account_code from chart_of_accounts) t2 on t2.id = t0.account_id
											left join (select id,issue_date,voucher_seq_no,status_id,case when status_id = 1 then 'Approved' else 'Open' end status_desc,status_gl_id,case when status_gl_id = 1 then 'Posted' else 'Open' end status_gl_desc,created_at,updated_at,created_by,updated_by from petty_cash_header where deleted_at is null) t4 on t4.id=t0.voucher_id
										where t0.deleted_at is null  ",theDateSet,receiptStatusSet,accountIdSet,userIdSet,keywordSet,voucherIdSet,"
											group by  t0.voucher_id,t0.account_id_header,t1.account_code_header,t1.account_name_header,t0.issue_date,t0.voucher_seq_no,t0.voucher_no,period,transaction_type,t4.status_id,t4.status_desc,t4.status_gl_id,t4.status_gl_desc,date(t4.created_at),date(t4.updated_at),ifnull(t4.updated_by,t4.created_by),voucher_code
									) x where account_id_header != 0 
									order by account_id_header,urut,period,voucher_seq_no,id) y where account_id_header != 0  ",ColumnSet,limitSet,";"));
    -- else
-- 		SET @s =  (concat ("select voucher_id,issue_date,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,debet,credit,memo_header,pic_header,updated_by,transaction_type,status_data,
-- 								'issue_date,account_code_header,account_name_header,voucher_code,voucher_no,debet,credit,pic_header,memo_header,transaction_type,updated_by,status_data' field_key,
--                                 'Tgl,Kode,Nama,Kode voucher,No voucher,Debet,Kredit,Pic,Memo,Transaksi,User,Status data' field_label,
--                                 'issue_date,account_code_header,account_name_header,voucher_code,voucher_no,debet,credit,pic_header,memo_header,transaction_type,updated_by,status_data' field_export,
--                                 'Tgl,Kode,Nama,Kode voucher,No voucher,Debet,Kredit,Pic,Memo,Transaksi,User,Status data' field_export_label,
--                                 'debet,credit' field_int,
--                                 'debet,credit' field_footer,
--                                 '' field_level
-- 							from (
-- 								select
-- 									t0.id,t1.id voucher_id,t1.issue_date,t1.voucher_seq_no,t1.voucher_code,t1.voucher_no,t0.account_id_header,t2.account_code_header,t2.account_name_header,t1.debet,t1.credit,t1.memo memo_header,t1.pic pic_header,item_no,account_id,t3.account_code,t3.account_name,t0.debet debet_detail,t0.credit credit_detail,t0.memo,
-- 									t0.created_at created_at_detail,t0.updated_at updated_at_detail,ifnull(t0.updated_by,t0.created_by) updated_by_detail,t1.transaction_type,
-- 									t1.created_at,t1.updated_at,ifnull(t1.updated_by,t1.created_by) updated_by,case when date(t1.created_at) = date(t1.updated_at) then 'NEW' else 'EDIT!!' end status_data
-- 								from 
-- 									petty_cash t0 
-- 								left join
-- 									(select id,issue_date,voucher_seq_no,voucher_code,voucher_no,debet,credit,period,status_id,case when status_id = 1 then 'Approved' else 'Open' end status_desc,case when status_gl_id = 1 then 'Posted' else 'Open' end status_gl_desc,pic,memo,created_at,updated_at,created_by,updated_by,transaction_type from petty_cash_header where deleted_at is null) t1 on t1.id=t0.voucher_id
-- 								left join
-- 									(select id,code_coa account_code_header,name_coa account_name_header from chart_of_accounts) t2 on t2.id=t0.account_id_header
-- 								left join
-- 									(select id,code_coa account_code,name_coa account_name from chart_of_accounts) t3 on t3.id=t0.account_id
-- 								where deleted_at is null ",theDateSet,updatedAtSet,receiptStatusSet,keywordSet,accountIdSet,userIdSet,") x where voucher_id != 0 ",ColumnSet,"
--                                 group by voucher_id,issue_date,voucher_seq_no,voucher_code,voucher_no,account_id_header,account_code_header,account_name_header,debet,credit,memo_header,pic_header,updated_by,status_data,transaction_type
-- 							order by issue_date desc ",limitSet,";"));
	end if;
	PREPARE stmt FROM @s;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;