CREATE USER 'usagemeter'@'mysqlcontainer' IDENTIFIED BY '7890';
GRANT ALL ON usagemeter.* TO 'usagemeter'@'mysqlcontainer';
CREATE USER 'usagemeter'@'%' IDENTIFIED BY '7890';
GRANT ALL PRIVILEGES ON *.* TO 'usagemeter'@'%' WITH GRANT OPTION;

CREATE TABLE IF NOT EXISTS managerMinQuota (
  id int NOT NULL AUTO_INCREMENT,
  plan varchar(50) DEFAULT NULL,
  minLimit int DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS pendingRequest (
  id int NOT NULL AUTO_INCREMENT,
  userChain varchar(50) DEFAULT NULL,
  isPending tinyint(1) DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS userDevices (
  id int NOT NULL AUTO_INCREMENT,
  userChain varchar(50) DEFAULT NULL,
  deviceIP varchar(16) DEFAULT NULL,
  isActive tinyint(1) DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS userManagers (
  userChain varchar(50) DEFAULT NULL,
  managerChain varchar(50) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS userMinQuota (
  id int NOT NULL AUTO_INCREMENT,
  plan varchar(50) DEFAULT NULL,
  minLimit int DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
  userChain varchar(50) NOT NULL,
  email varchar(50) DEFAULT NULL,
  isManager tinyint(1) DEFAULT NULL,
  defaultQuota int DEFAULT NULL,
  isAdmin tinyint(1) DEFAULT NULL,
  PRIMARY KEY (userChain),
  UNIQUE KEY userChain (userChain)
);