USE adds;

DROP TABLE IF EXISTS adds;
CREATE TABLE adds (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	uid INT NOT NULL,
	title VARCHAR(255) NOT NULL,
	description TEXT,
	price DECIMAL(10,2) NOT NULL DEFAULT 0,
	category INT NOT NULL DEFAULT 0,
	region INT NOT NULL DEFAULT 0,
	negotiable BOOLEAN NOT NULL DEFAULT 0,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	image VARCHAR(100) DEFAULT "",
	view_count INT  NOT NULL DEFAULT 0
 );

DROP TABLE IF EXISTS photos;
CREATE TABLE photos (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	uid INT NOT NULL,
	aid INT NOT NULL,
	filename VARCHAR(100) DEFAULT "",
	screen TINYINT(1) DEFAULT 0,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
 );


DROP TABLE IF EXISTS users;
CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(100) NOT NULL,
	password BINARY(60),
	email VARCHAR(100) NOT NULL,
	token VARCHAR(100) NOT NULL,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO adds (uid, title, description, price, category, region, negotiable) VALUES
(1, "20 jant Bianchi çocuk bisikleti", "Çok iyi durumda, temiz kullanılmış, herhangi bir sorunu yok, lastikleri geçen yıl değişti", 50, 1, 1, 1),
(1, "i5-2240 işlemci ve anakart", "Test edip alabilirsiniz, güncel bütün oyunları oynayabilirsiniz.", 150, 1, 1, 0),
(2, "IKEA yataklı kanepe", "Açılınca 120 cm yatak oluyor. Ankara içinde istediğiniz adrese teslim edilir.", 179.99, 2, 2, 1),
(3, "Macbook Air 13.3 inch laptop", "Mid 2013 versiyonu. 120gb ssd, 4gb ram, 1,3 GHz Intel Core i5 işlemci", 2500, 1, 3, 0);

