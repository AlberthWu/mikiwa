SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_ProductUom;
DELIMITER $$
CREATE PROCEDURE sp_ProductUom(uId int,uomId int,userId int) 
BEGIN
	declare uIdSet varchar(255);
	declare uomIdSet varchar(255);
    SET uIdSet  = case when (uId is null or uId = 0) then '' else concat(' where t0.product_id = ',uId,'')   end;
    SET uomIdSet  = case when (uomId is null or uomId = 0) then '' else concat(' where uom_id = ',uomId,'')   end;
	SET @s =  (concat ('select product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,next_uom_id,next_uom_code,is_default from (
							select product_id,product_code,product_name,item_no,uom_id,uom_code,ratio,lag(uom_id) over (partition by product_id order by item_no) next_uom_id,is_default
							from product_uom t0
								left join (select id,product_code,product_name from products) t1 on t1.id=t0.product_id
								left join (select id,uom_code from uoms) t2 on t2.id=t0.uom_id
							',uIdSet,') x
								left join (select id,uom_code next_uom_code from uoms) t1 on t1.id=x.next_uom_id 
                            ',uomIdSet,' order by product_id,item_no;'));
    
    PREPARE stmt FROM @s;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
