# A Very Simple Bockchain For Tracking Rent Payment

This repo followed tutorial from: https://github.com/codehakase/blockchain-golang-tutorial

## Installation
Clone Repository
```
git clone https://github.com/harrynv11/rental-tracker-blockchain.git
```

## Run
```
go run app.go
```

## Create Blocks with Postman
Create a new renter
![Imgur](https://i.imgur.com/N23Jz31.png)

Store the renter payment to block
![Imgur](https://i.imgur.com/8zFPvkY.png)

## Or with cURL
Create a new renter
```
$ curl -X POST http://localhost:3000/new 
  -H "Content-Type: application/json" 
  -d '{"name":"Na Na", "join_date":"2019-17-04"}'
```

Create a new block with payment date
```
$ curl -X POST http://localhost:3000
  -H "Content-Type: application/json" 
  -d '{"renter":"Na Na", "pay_date":"2019-18-04"}'
```

All blocks are also recorded at http://localhost:1911
![Imgur](https://i.imgur.com/7TeRv3p.png)
