docker run --name pgadmin-container -p 5050:80 -e PGADMIN_DEFAULT_EMAIL=teste@teste.com -e PGADMIN_DEFAULT_PASSWORD=teste -d dpage/pgadmin4

# queries
# SELECT * FROM orders;

# SELECT * FROM products;

# SELECT * FROM orders WHERE id = 1067;

# SELECT * FROM users WHERE name = 'Junita Jast';

# SELECT * FROM orders WHERE id = 1067;
# SELECT * FROM orders WHERE user_id = 99;

# SELECT * FROM order_products WHERE order_id= 1067;

# SELECT
#     u.id,
#     u.name,
#     o.id AS order_id,
#     o.date,
#     op.product_id,
#     op.value
# FROM
#     users u
# JOIN
#     orders o ON u.id = o.user_id
# JOIN
#     order_products op ON o.id = op.order_id
# WHERE
#     u.id = 99;