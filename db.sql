CREATE TABLE fundsusd (
    id bigserial not null primary key,
    name varchar ,
    ticker varchar ,
    amount int ,
    priceperitem decimal,
    purchaseprice decimal,
    pricecurrent decimal,
    percentchanges decimal,
    yearlyinvestment decimal,
    clearmoney decimal,
    datepurchase date,
    datelastupdate date,
    type varchar,
);


INSERT INTO fundsusd (name, ticker, amount, priceperitem, purchaseprice, pricecurrent, percentchanges, yearlyinvestment, clearmoney, datepurchase, datelastupdate, type) VALUES ( 'MOMO', 'MOMO', 35, 12.39, 512.4 , 455, -11.2021857923, -37.4841219573, -57.8837, '2020-09-10 00:00:30', '2021-08-02 15:24:52', 'share' );