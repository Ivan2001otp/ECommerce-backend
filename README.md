# ECommerce Backend in Go
Backend for E-Commerce with JWT authentication.

## Models
- Product
- CategoryOfProduct
- User
- Order

## Endpoints
- localhost:8080/users/signup(POST)
- localhost:8080/users/login(POST)
- localhost:8080/users(GET)
- localhost:8080/admin/category/add(POST)
- localhost:8080/admin/categories/:category_id(GET/UPDATE)
- localhost:8080/admin/categories(GET)
- localhost:8080/admin/product/add(POST)
- localhost:8080/admin/products(GET)
- localhost:8080/admin/products/:product_id(GET)
- localhost:8080/admin/product/:product_id(UPDATE)
- localhost:8080/order/:user_id(POST)
- localhost:8080/admin/orders(GET)
- localhost:8080/order/:order_id(GET)
