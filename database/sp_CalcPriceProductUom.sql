SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_CalcPriceProductUom;
DELIMITER $$
CREATE PROCEDURE sp_CalcPriceProductUom(uId int,priceId int,userId int) 
BEGIN
	declare itemNo int default 1;
    declare rn int;
    declare isDefault int default 0;
    declare qty decimal(12,4);
    declare qtyDefault decimal(12,4);
    set rn = (select max(item_no) from price_product_uom where deleted_at is null and price_id = priceId and product_id = uId);
    set isDefault = (select item_no from price_product_uom where deleted_at is null and price_id = priceId and product_id = uId and is_default = 1);
    
    while itemNo <= rn do
		-- final_ratio
        if itemNo = 1 then
			set qty = 1;
		else
			set qty = (select final_ratio from price_product_uom where deleted_at is null and price_id = priceId and product_id = uId and item_no = itemNo-1);
        end if;
		
        -- default_ratio
        if itemNo = isDefault then
			set qtyDefault = 1;
		else
			set qtyDefault = ifnull((select is_default_ratio from price_product_uom where deleted_at is null and price_id = priceId and product_id = uId and item_no = itemNo-1),0)*ifnull((select ratio from price_product_uom where deleted_at is null and price_id = priceId and product_id = uId and item_no = itemNo),0);
        end if;
        
        update price_product_uom set final_ratio = qty*ratio ,is_default_ratio = qtyDefault where deleted_at is null and price_id = priceId and product_id = uId and item_no = itemNo;
		set itemNo = itemNo+1;
    end while;
    
	select * from price_product_uom t0
		left join (select id,uom_code from uoms) t1 on t1.id= t0.uom_id
	where deleted_at is null and price_id = priceId and product_id = uId order by item_no;
END$$

DELIMITER ;
