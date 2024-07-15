SET GLOBAL log_bin_trust_function_creators = 1;

DROP PROCEDURE IF EXISTS sp_ConvertUom;
DELIMITER $$
CREATE PROCEDURE sp_ConvertUom(issueDate date,uId int,companyId int,uomId int,qtyId decimal(12,2),userId int) 
BEGIN
	declare priceId int default 0;
    declare NettPrice decimal(18,2);
    -- check customer is special or regular
    set priceId = (select id from (select id,DENSE_RANK() OVER (PARTITION BY company_id,product_id,price_type order by effective_date desc,ifnull(expired_date,'9999-12-31'), id desc) dr from prices 
					where deleted_at is null and price_type = 'sales' and issueDate between effective_date and ifnull(expired_date,'9999-12-31') and company_id = companyId and product_id=uId order by id desc) x where dr = 1 );
	if qtyId = 0 then
		set qtyId = 1;
    end if;
    -- if regular customer
	if ifnull(priceId,0) = 0 then
		select ifnull(priceId,0) price_id,product_id,product_code,product_name,art_no,barcode,normal_price,price,qty,item_no,uom_id,uom_code,ratio,packaging_id,packaging_code,final_qty,final_uom_id,final_uom_code,conversion_qty,conversion_uom_id,conversion_uom_code from (
		select uId product_id,product_code,product_name,art_no,barcode,normal_price,case when conversion_uom_id = uomId then ceiling(round(normal_price/t3.final_ratio,0)/500)*500 else normal_price end price,qtyId qty,item_no,t0.uom_id,
			case when final_uom_id = uomId then t3.is_default_ratio else ratio end ratio,
			case when conversion_uom_id = uomId then uom_id when uomId = final_uom_id then uomId else lag(t0.uom_id) over (partition by product_id order by item_no) end packaging_id,
			case when conversion_uom_id = uomId then qtyId else qtyId*t0.is_default_ratio end final_qty,
			case when conversion_uom_id = uomId then uom_id else final_uom_id end final_uom_id,
			case when conversion_uom_id = uomId then t4.uom_code else final_uom_code end final_uom_code, qtyId*t0.final_ratio conversion_qty,t2.conversion_uom_id,t2.conversion_uom_code
		from product_uom t0 
			left join (select id,product_code,product_name,art_no,barcode from products where id=uId) t1 on 1=1
			left join (select uom_id conversion_uom_id,uom_code conversion_uom_code from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where product_id=uId and item_no = 1) t2 on 1=1
			left join (select uom_id final_uom_id,uom_code final_uom_code,is_default_ratio,price normal_price,final_ratio from product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where is_default=1 and product_id=uId) t3 on 1=1
			left join (select id,uom_code from uoms where id = uomId) t4 on 1=1
		where product_id = uId )x 
			left join (select id,uom_code  from uoms) t1 on t1.id=x.uom_id
			left join (select id,uom_code packaging_code from uoms) t2 on t2.id=x.packaging_id
		where uom_id = uomId;
	-- if special customer
	else
		set NettPrice = (select cast((t1.normal_price + ((t1.normal_price*disc_one)/100) + (t1.normal_price + (t1.normal_price*disc_one)/100)*disc_two/100) + disc_tpr as decimal(18,2)) price from price_product_uom t0 left join (select id,price normal_price from products where deleted_at is null and id = uId) t1 on 1=1 where deleted_at is null and  is_default = 1 and price_id = priceId);
        select ifnull(priceId,0) price_id,product_id,product_code,product_name,art_no,barcode,normal_price,price,qty,item_no,uom_id,uom_code,ratio,packaging_id,packaging_code,final_qty,final_uom_id,final_uom_code,conversion_qty,conversion_uom_id,conversion_uom_code from (
		select uId product_id,product_code,product_name,art_no,barcode,normal_price,case when conversion_uom_id = uomId then ceiling(round(NettPrice/t3.final_ratio,0)/500)*500 else NettPrice end price,qtyId qty,item_no,t0.uom_id,
			case when final_uom_id = uomId then t3.is_default_ratio else ratio end ratio,
			case when conversion_uom_id = uomId then uom_id when uomId = final_uom_id then uomId else lag(t0.uom_id) over (partition by product_id order by item_no) end packaging_id,
			case when conversion_uom_id = uomId then qtyId else qtyId*t0.is_default_ratio end final_qty,
			case when conversion_uom_id = uomId then uom_id else final_uom_id end final_uom_id,
			case when conversion_uom_id = uomId then t4.uom_code else final_uom_code end final_uom_code, qtyId*t0.final_ratio conversion_qty,t2.conversion_uom_id,t2.conversion_uom_code
		from price_product_uom t0 
			left join (select id,product_code,product_name,price normal_price,art_no,barcode from products where id=uId) t1 on 1=1
			left join (select t0.id,uom_id conversion_uom_id,uom_code conversion_uom_code from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and price_id = priceId and item_no = 1) t2 on 1=1
			left join (select t0.id,uom_id final_uom_id,uom_code final_uom_code,is_default_ratio,final_ratio from price_product_uom t0 left join (select id,uom_code from uoms) t1 on t1.id=t0.uom_id where deleted_at is null and  is_default=1 and  price_id = priceId) t3 on 1=1
			left join (select id,uom_code from uoms where id = uomId) t4 on 1=1
		where  deleted_at is null and price_id = priceId)x 
			left join (select id,uom_code  from uoms) t1 on t1.id=x.uom_id
			left join (select id,uom_code packaging_code from uoms) t2 on t2.id=x.packaging_id
		where  uom_id = uomId;
    end if;
    
END$$

DELIMITER ;
