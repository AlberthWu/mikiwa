SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_ConvertUom;
DELIMITER $$
CREATE PROCEDURE sp_ConvertUom(uId int,uomId int,qtyId decimal(12,2),userId int) 
BEGIN
	declare isDefault int default 0;
	
    select product_id,product_code,product_name,qty,uom_id,uom_code,ratio,packaging_id,packaging_code,final_qty,final_uom_id,final_uom_code,conversion_qty,conversion_uom_id,conversion_uom_code from (
    select uId product_id,product_code,product_name,qtyId qty,t0.uom_id,case when uomId = final_uom_id then t3.is_default_ratio else ratio end ratio,case when uomId = final_uom_id then uomId else lag(t0.uom_id) over (partition by product_id order by item_no) end packaging_id,qtyId*t0.is_default_ratio final_qty,final_uom_id,final_uom_code, qtyId*final_ratio conversion_qty,t2.conversion_uom_id,t2.conversion_uom_code
    from product_uom t0 
		left join (select id,product_code,product_name from products where id=uId) t1 on 1=1
        left join (select id,uom_id conversion_uom_id,uom_code conversion_uom_code from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where product_id=uId and item_no = 1) t2 on 1=1
        left join (select id,uom_id final_uom_id,uom_code final_uom_code,is_default_ratio from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where is_default=1 and product_id=uId) t3 on 1=1
    where product_id = uId )x 
		left join (select id,uom_code  from uoms) t1 on t1.id=x.uom_id
        left join (select id,uom_code packaging_code from uoms) t2 on t2.id=x.packaging_id
    where uom_id = uomId;
    
END$$

DELIMITER ;
