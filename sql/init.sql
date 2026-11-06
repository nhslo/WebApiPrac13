CREATE TABLE IF NOT EXISTS products (
 id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, category VARCHAR(50) NOT NULL,
 price NUMERIC(12,2) NOT NULL CHECK(price >= 0), quantity INTEGER NOT NULL CHECK(quantity >= 0)
);
INSERT INTO products (name,category,price,quantity)
SELECT s.name,s.category,s.price,s.quantity FROM (VALUES
 ('Laptop','Electronics',350000,4), ('Mouse','Accessories',5000,15), ('Monitor','Electronics',85000,3),
 ('Keyboard','Accessories',12000,8), ('Headphones','Accessories',25000,2)
) AS s(name,category,price,quantity) WHERE NOT EXISTS (SELECT 1 FROM products);
