# TEST 01
## Database Instalation 
- Mongod
- Please change .env for database cridentials
- please run script on this bellow :
  - db.products.insert( { sku: "120P90", name:"Google Home", price:49.99, qty: 10 } )
  - db.products.insert( { sku: "43N23P", name:"Macbook Pro", price:5399.99, qty: 5 } )
  - db.products.insert( { sku: "A304SD", name:"Lexa Speaker", price:109.5, qty: 10 } )
  - db.products.insert( { sku: "234234", name:"Raspberry Pi B", price:30, qty: 2 } )
  - db.promotions.insert( { promo_code: "promo3", promo_type:"discountAlexa", minimum_qty:3, discount: 10, sku:"A304SD", details:[] } )
  - db.promotions.insert( { promo_code: "promo1", promo_type:"free_sku", minimum_qty:1, sku:"43N23P", details:[{sku:"234234",qty:1}] } )
  - db.promotions.insert( { promo_code: "promo2", promo_type:"free_one_sku", minimum_qty:3, discount: 10, sku:"120P90", details:[] } )
- after run go, please use postman and read Testing01.postman_collection.json file for requirenment data
