SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_GenerateNumber;

DELIMITER $$
CREATE PROCEDURE sp_GenerateNumber(issueDate date,headerId varchar(45),outletId int,customerId int) 
BEGIN
	DECLARE Rn int;
    Declare Urut int DEFAULT 1;
    declare SeqNo int;
    declare headerIdSet varchar(255);
    declare referenceNo varchar(25);
    declare outletId varchar(5) default (select `code` from plants where id = outletId);
	declare Toroman varchar(5) default (select fn_ToRoman(month(issueDate))) ;
 
    SET headerIdSet = headerId;
    DROP TEMPORARY TABLE  IF EXISTS temp_invoice_number;
    CREATE TEMPORARY TABLE temp_invoice_number(id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,data_id int,seq_no int);
    if headerIdSet = "SalesOrder" THEN
		set seqNo = ifnull((select max(seq_no)+1 from sales_order where year(issue_date) = year(issueDate) and outlet_id = outletId ),1);
		set referenceNo  = concat(ifnull(outletId,"00"),'-SO/',Toroman,'-',date_format(issueDate,'%y'),'/',LPAD(seqNo,5,0)) ;
	else if headerIdSet = "DeliveryOrder" THEN
    END IF;
	select seqNo seqno, referenceNo `format`;
END$$
DELIMITER ;
