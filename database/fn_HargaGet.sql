SET GLOBAL log_bin_trust_function_creators = 1;

DROP FUNCTION IF EXISTS fn_HargaGet;

DELIMITER $$
CREATE FUNCTION fn_HargaGet(theDate date,customerId int,productId int,uomId int, qtyId decimal(5,2), discOne decimal(5,2), discTwo decimal(5,2), discTpr decimal(12,2),methodId int) 
RETURNS decimal(18,2)

BEGIN

	declare price decimal(18,2);
	declare prices decimal(18,2);
    declare priceSet int;
    
    -- new calculate
    if methodId = 0 then
		set prices = (select price from product_uom where product_id = productId and uom_id = uomId);
	-- database
    else
		set priceSet = (select price from (select t0.id,effective_date,expired_date,t2.uom_id,t2.uom_code,t1.normal_price,disc_one,disc_two,disc_tpr,
							cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price,
							DENSE_RANK() OVER (PARTITION BY t0.company_id,t0.product_id,price_type order by effective_date desc,ifnull(expired_date,'9999-12-31'), id desc) dr 
						from prices t0
							left join (select product_id,uom_id,price normal_price from product_uom where product_id=productId and uom_id=uomId) t1 on 1=1
							left join (select price_id,uom_id,uom_code from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and  product_id=productId  and uom_id=uomId) t2 on t2.price_id=t0.id
						where deleted_at is null and price_type = 'sales' and theDate between effective_date and ifnull(expired_date,'9999-12-31') and id in (select price_id from price_company where company_id = customerId) and t0.product_id=productId) x where dr=1);
		if ifnull(priceSet,0) = 0 then
			set prices = (select price from product_uom where product_id = productId and uom_id = uomId);
		else
			set prices = priceSet;
		end if;
    end if;
    
    set price = (prices + ((prices*discOne)/100) + (prices + (prices*discOne)/100)*discTwo/100) + discTpr ;
    return price;

END$$

DELIMITER ;