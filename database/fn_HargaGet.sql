SET GLOBAL log_bin_trust_function_creators = 1;

DROP FUNCTION IF EXISTS fn_HargaGet;

DELIMITER $$
CREATE FUNCTION fn_HargaGet(theDate date,customerId int,productId int,uomId int, qtyId decimal(5,2), discOne decimal(5,2), discTwo decimal(5,2), discTpr decimal(12,2),methodId int) 
RETURNS decimal(18,2)

BEGIN

	declare prices decimal(18,2);
    declare priceSet int;
    declare priceId int default 0;
    declare defaultRatio decimal(18,2);
    declare NettPrice decimal(18,2);
    set priceId = (select id from (select id,DENSE_RANK() OVER (PARTITION BY company_id,product_id,price_type order by effective_date desc,ifnull(expired_date,'9999-12-31'), id desc) dr from prices 
					where deleted_at is null and price_type = 'sales' and theDate between effective_date and ifnull(expired_date,'9999-12-31') and company_id = customerId and product_id=productId order by id desc) x where dr = 1 );
    -- new calculate
    if ifnull(priceId,0) = 0 then
		set prices = (select price from (
							select productId product_id,product_code,product_name,art_no,barcode,normal_price,case when conversion_uom_id = uomId then ceiling(round(normal_price/t3.final_ratio,0)/500)*500 else normal_price end price,qtyId qty,item_no,t0.uom_id,
								case when final_uom_id = uomId then t3.is_default_ratio else ratio end ratio,
								case when conversion_uom_id = uomId then uom_id when uomId = final_uom_id then uomId else lag(t0.uom_id) over (partition by product_id order by item_no) end packaging_id,
								case when conversion_uom_id = uomId then qtyId else qtyId*t0.is_default_ratio end final_qty,
								case when conversion_uom_id = uomId then uom_id else final_uom_id end final_uom_id,
								case when conversion_uom_id = uomId then t4.uom_code else final_uom_code end final_uom_code, qtyId*t0.final_ratio conversion_qty,t2.conversion_uom_id,t2.conversion_uom_code
							from product_uom t0 
								left join (select id,product_code,product_name,art_no,barcode from products where id=productId) t1 on 1=1
								left join (select uom_id conversion_uom_id,uom_code conversion_uom_code from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where product_id=productId and item_no = 1) t2 on 1=1
								left join (select uom_id final_uom_id,uom_code final_uom_code,is_default_ratio,price normal_price,final_ratio from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where is_default=1 and product_id=productId) t3 on 1=1
								left join (select id,uom_code from uoms where id = uomId) t4 on 1=1
							where product_id = productId )x 
								left join (select id,uom_code  from uoms) t1 on t1.id=x.uom_id
								left join (select id,uom_code packaging_code from uoms) t2 on t2.id=x.packaging_id
							where uom_id = uomId);
	-- database
    else
		-- set priceSet = (select price from (select t0.id,effective_date,expired_date,t2.uom_id,t2.uom_code,t1.normal_price,disc_one,disc_two,disc_tpr,
-- 							cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price,
-- 							DENSE_RANK() OVER (PARTITION BY t0.company_id,t0.product_id,price_type order by effective_date desc,ifnull(expired_date,'9999-12-31'), id desc) dr 
-- 						from prices t0
-- 							left join (select product_id,uom_id,price normal_price from product_uom where product_id=productId ) t1 on 1=1
-- 							left join (select price_id,uom_id,uom_code from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and  product_id=productId  and uom_id=uomId) t2 on t2.price_id=t0.id
-- 						where deleted_at is null and price_type = 'sales' and theDate between effective_date and ifnull(expired_date,'9999-12-31') and id in (select price_id from price_company where company_id = customerId) and t0.product_id=productId ) x where dr=1);
		
		set NettPrice = (select cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price from price_product_uom t0 left join (select id,price normal_price from products where deleted_at is null and id = productId) t1 on 1=1 where deleted_at is null and  is_default = 1 and price_id = priceId);
        set prices = (select price from (
							select productId product_id,product_code,product_name,art_no,barcode,normal_price,case when conversion_uom_id = uomId then ceiling(round(NettPrice/t3.final_ratio,0)/500)*500 else NettPrice end price,qtyId qty,item_no,t0.uom_id,
								case when final_uom_id = uomId then t3.is_default_ratio else ratio end ratio,
								case when conversion_uom_id = uomId then uom_id when uomId = final_uom_id then uomId else lag(t0.uom_id) over (partition by product_id order by item_no) end packaging_id,
								case when conversion_uom_id = uomId then qtyId else qtyId*t0.is_default_ratio end final_qty,
								case when conversion_uom_id = uomId then uom_id else final_uom_id end final_uom_id,
								case when conversion_uom_id = uomId then t4.uom_code else final_uom_code end final_uom_code, qtyId*t0.final_ratio conversion_qty,t2.conversion_uom_id,t2.conversion_uom_code
							from price_product_uom t0 
								left join (select id,product_code,product_name,price normal_price,art_no,barcode from products where id=productId) t1 on 1=1
								left join (select t0.id,uom_id conversion_uom_id,uom_code conversion_uom_code from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and price_id = priceId and item_no = 1) t2 on 1=1
								left join (select t0.id,uom_id final_uom_id,uom_code final_uom_code,is_default_ratio,final_ratio from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and  is_default=1 and  price_id = priceId) t3 on 1=1
								left join (select id,uom_code from uoms where id = uomId) t4 on 1=1
							where  deleted_at is null and price_id = priceId)x 
								left join (select id,uom_code  from uoms) t1 on t1.id=x.uom_id
								left join (select id,uom_code packaging_code from uoms) t2 on t2.id=x.packaging_id
							where  uom_id = uomId);
		end if;
		-- if ifnull(priceSet,0) = 0 then
-- 			set prices = (select price from product_uom where product_id = productId and uom_id = uomId);
-- 		else
-- 			set prices = priceSet;
-- 		end if;
    
		return (prices + ((prices*discOne)/100) + (prices + (prices*discOne)/100)*discTwo/100) + discTpr ;


END$$

DELIMITER ;