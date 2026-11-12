-- Demo catalog data for a new or development database.
-- Re-running this file updates these sample pizzas and replaces their
-- categories, images, and variants. It does not create fake users or reviews.
BEGIN;

INSERT INTO categories (name, slug, description) VALUES
    ('คลาสสิก', 'classic', 'พิซซ่าสูตรคลาสสิกที่เป็นที่รู้จัก'),
    ('มังสวิรัติ', 'vegetarian', 'หน้าไม่มีเนื้อสัตว์'),
    ('ชีส', 'cheese', 'พิซซ่าสำหรับคนรักชีส'),
    ('เนื้อสัตว์', 'meat', 'พิซซ่าหน้าเนื้อสัตว์หลากชนิด'),
    ('เผ็ด', 'spicy', 'พิซซ่ารสเผ็ด'),
    ('ซีฟู้ด', 'seafood', 'พิซซ่าหน้าอาหารทะเล'),
    ('รสหวานเค็ม', 'sweet-savory', 'รสชาติหวานและเค็มเข้ากัน')
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = NOW();

INSERT INTO pizzas (name, slug, description, status) VALUES
    ('Margherita', 'margherita', 'พิซซ่าสไตล์อิตาเลียนหน้าเรียบง่าย ใช้มะเขือเทศ มอซซาเรลลา และใบโหระพา', 'published'),
    ('Marinara', 'marinara', 'หน้ามะเขือเทศ กระเทียม ออริกาโน และน้ำมันมะกอก โดยทั่วไปไม่มีชีส', 'published'),
    ('Pepperoni', 'pepperoni', 'มอซซาเรลลากับเปปเปอโรนีรสเข้มข้น อบจนขอบไส้กรอกเกรียมนิด ๆ', 'published'),
    ('Hawaiian', 'hawaiian', 'แฮมและสับปะรดให้รสเค็มหวาน ตัดกับชีสและซอสมะเขือเทศ', 'published'),
    ('Quattro Formaggi', 'quattro-formaggi', 'พิซซ่าสี่ชีส ผสมชีสเนื้อนุ่มและชีสรสเข้มเพื่อให้ได้รสหลายมิติ', 'published'),
    ('Diavola', 'diavola', 'ซาลามีรสจัดกับพริก เพิ่มความเผ็ดบนฐานมะเขือเทศและมอซซาเรลลา', 'published'),
    ('BBQ Chicken', 'bbq-chicken', 'ไก่ซอสบาร์บีคิว หอมแดง และมอซซาเรลลา รสหวานรมควัน', 'published'),
    ('Meat Lovers', 'meat-lovers', 'รวมแฮม เปปเปอโรนี เบคอน และไส้กรอกบนหน้าชีส', 'published'),
    ('Seafood', 'seafood', 'กุ้งและปลาหมึกกับกระเทียม มอซซาเรลลา และซอสพิซซ่ารสอ่อน', 'published'),
    ('Pesto Mushroom', 'pesto-mushroom', 'เห็ดผัดหอม ๆ เพสโต และมอซซาเรลลา รสกลมกล่อมแบบไม่มีเนื้อสัตว์', 'published'),
    ('New York Cheese', 'new-york-cheese', 'พิซซ่าชีสหน้าเรียบง่าย แป้งบางพับได้ ได้แรงบันดาลใจจากสไตล์นิวยอร์ก', 'published'),
    ('Veggie Supreme', 'veggie-supreme', 'พริกหวาน เห็ด หอมใหญ่ มะกอก และมอซซาเรลลาสำหรับคนชอบผัก', 'published')
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    updated_at = NOW();

DELETE FROM pizza_categories pc
USING pizzas p
WHERE pc.pizza_id = p.id
  AND p.slug IN (
    'margherita', 'marinara', 'pepperoni', 'hawaiian',
    'quattro-formaggi', 'diavola', 'bbq-chicken', 'meat-lovers',
    'seafood', 'pesto-mushroom', 'new-york-cheese', 'veggie-supreme'
  );

INSERT INTO pizza_categories (pizza_id, category_id)
SELECT p.id, c.id
FROM (VALUES
    ('margherita', 'classic'),
    ('margherita', 'vegetarian'),
    ('marinara', 'classic'),
    ('marinara', 'vegetarian'),
    ('pepperoni', 'classic'),
    ('pepperoni', 'meat'),
    ('hawaiian', 'sweet-savory'),
    ('hawaiian', 'meat'),
    ('quattro-formaggi', 'cheese'),
    ('quattro-formaggi', 'vegetarian'),
    ('diavola', 'spicy'),
    ('diavola', 'meat'),
    ('bbq-chicken', 'meat'),
    ('meat-lovers', 'meat'),
    ('seafood', 'seafood'),
    ('pesto-mushroom', 'vegetarian'),
    ('pesto-mushroom', 'cheese'),
    ('new-york-cheese', 'classic'),
    ('new-york-cheese', 'cheese'),
    ('veggie-supreme', 'vegetarian')
) AS mapping(pizza_slug, category_slug)
JOIN pizzas p ON p.slug = mapping.pizza_slug
JOIN categories c ON c.slug = mapping.category_slug
ON CONFLICT DO NOTHING;

DELETE FROM pizza_images
WHERE pizza_id IN (
    SELECT id FROM pizzas
    WHERE slug IN (
        'margherita', 'marinara', 'pepperoni', 'hawaiian',
        'quattro-formaggi', 'diavola', 'bbq-chicken', 'meat-lovers',
        'seafood', 'pesto-mushroom', 'new-york-cheese', 'veggie-supreme'
    )
);

INSERT INTO pizza_images (pizza_id, image_url, alt_text, sort_order)
SELECT p.id, image.image_url, image.alt_text, 0
FROM (VALUES
    ('margherita', 'https://images.unsplash.com/photo-1579751626657-72bc17010498?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Margherita'),
    ('marinara', 'https://images.unsplash.com/photo-1571407970349-bc81e7e96d47?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Marinara'),
    ('pepperoni', 'https://images.unsplash.com/photo-1628840042765-356cda07504e?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Pepperoni'),
    ('hawaiian', 'https://images.unsplash.com/photo-1528137871618-79d2761e3fd5?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Hawaiian'),
    ('quattro-formaggi', 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Quattro Formaggi'),
    ('diavola', 'https://images.unsplash.com/photo-1628840042765-356cda07504e?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Diavola'),
    ('bbq-chicken', 'https://images.unsplash.com/photo-1528137871618-79d2761e3fd5?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for BBQ Chicken'),
    ('meat-lovers', 'https://images.unsplash.com/photo-1628840042765-356cda07504e?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Meat Lovers'),
    ('seafood', 'https://images.unsplash.com/photo-1571407970349-bc81e7e96d47?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Seafood'),
    ('pesto-mushroom', 'https://images.unsplash.com/photo-1593560708920-61dd98c46a4e?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Pesto Mushroom'),
    ('new-york-cheese', 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for New York Cheese'),
    ('veggie-supreme', 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&w=1200&q=80', 'Illustrative pizza photo for Veggie Supreme')
) AS image(pizza_slug, image_url, alt_text)
JOIN pizzas p ON p.slug = image.pizza_slug;

DELETE FROM pizza_variants
WHERE pizza_id IN (
    SELECT id FROM pizzas
    WHERE slug IN (
        'margherita', 'marinara', 'pepperoni', 'hawaiian',
        'quattro-formaggi', 'diavola', 'bbq-chicken', 'meat-lovers',
        'seafood', 'pesto-mushroom', 'new-york-cheese', 'veggie-supreme'
    )
);

INSERT INTO pizza_variants (pizza_id, name, size_label, price, currency, is_available, sort_order)
SELECT p.id, variant.name, variant.size_label, variant.price, 'THB', TRUE, variant.sort_order
FROM (VALUES
    ('margherita', 'Regular', 'กลาง 9 นิ้ว', 199.00::NUMERIC, 0),
    ('margherita', 'Large', 'ใหญ่ 12 นิ้ว', 299.00::NUMERIC, 1),
    ('marinara', 'Regular', 'กลาง 9 นิ้ว', 179.00::NUMERIC, 0),
    ('marinara', 'Large', 'ใหญ่ 12 นิ้ว', 269.00::NUMERIC, 1),
    ('pepperoni', 'Regular', 'กลาง 9 นิ้ว', 229.00::NUMERIC, 0),
    ('pepperoni', 'Large', 'ใหญ่ 12 นิ้ว', 349.00::NUMERIC, 1),
    ('hawaiian', 'Regular', 'กลาง 9 นิ้ว', 219.00::NUMERIC, 0),
    ('hawaiian', 'Large', 'ใหญ่ 12 นิ้ว', 329.00::NUMERIC, 1),
    ('quattro-formaggi', 'Regular', 'กลาง 9 นิ้ว', 249.00::NUMERIC, 0),
    ('quattro-formaggi', 'Large', 'ใหญ่ 12 นิ้ว', 369.00::NUMERIC, 1),
    ('diavola', 'Regular', 'กลาง 9 นิ้ว', 239.00::NUMERIC, 0),
    ('diavola', 'Large', 'ใหญ่ 12 นิ้ว', 359.00::NUMERIC, 1),
    ('bbq-chicken', 'Regular', 'กลาง 9 นิ้ว', 239.00::NUMERIC, 0),
    ('bbq-chicken', 'Large', 'ใหญ่ 12 นิ้ว', 359.00::NUMERIC, 1),
    ('meat-lovers', 'Regular', 'กลาง 9 นิ้ว', 259.00::NUMERIC, 0),
    ('meat-lovers', 'Large', 'ใหญ่ 12 นิ้ว', 389.00::NUMERIC, 1),
    ('seafood', 'Regular', 'กลาง 9 นิ้ว', 269.00::NUMERIC, 0),
    ('seafood', 'Large', 'ใหญ่ 12 นิ้ว', 399.00::NUMERIC, 1),
    ('pesto-mushroom', 'Regular', 'กลาง 9 นิ้ว', 229.00::NUMERIC, 0),
    ('pesto-mushroom', 'Large', 'ใหญ่ 12 นิ้ว', 339.00::NUMERIC, 1),
    ('new-york-cheese', 'Regular', 'กลาง 9 นิ้ว', 199.00::NUMERIC, 0),
    ('new-york-cheese', 'Large', 'ใหญ่ 12 นิ้ว', 299.00::NUMERIC, 1),
    ('veggie-supreme', 'Regular', 'กลาง 9 นิ้ว', 219.00::NUMERIC, 0),
    ('veggie-supreme', 'Large', 'ใหญ่ 12 นิ้ว', 329.00::NUMERIC, 1)
) AS variant(pizza_slug, name, size_label, price, sort_order)
JOIN pizzas p ON p.slug = variant.pizza_slug;

COMMIT;
