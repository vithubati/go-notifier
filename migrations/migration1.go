package migrations

const (
	migration1 = `
	CREATE TABLE IF NOT EXISTS notification
	(
		id 				varchar(36) PRIMARY KEY,
		resource        VARCHAR(50) NOT NULL,
		action          VARCHAR(50) NOT NULL,
		subject         VARCHAR(50) NOT NULL,
		message         VARCHAR(256) NOT NULL,
		createdAt       DATETIME,
		data            TEXT
	);

	CREATE TABLE IF NOT EXISTS deliverer
	(
		id 					varchar(36) PRIMARY KEY,
		type        		VARCHAR(50) NOT NULL,
		url          		VARCHAR(50) default '',
		channelId          	VARCHAR(50) default '',
		headers         	VARCHAR(50),
		credentials     	VARCHAR(256) default "",
		createdAt       	DATETIME,
		retry       		INT(11),
		intervalInSeconds   INT(11)
	);

	CREATE TABLE IF NOT EXISTS deliverer_resource
	(
		deliverer_id 		varchar(36) NOT NULL,
		resource 			varchar(36) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS delivery
	(
		id 					varchar(36) PRIMARY KEY,
		notificationId 		varchar(36) NOT NULL,
		delivererId 		varchar(36) NOT NULL,
		status          	VARCHAR(50) NOT NULL,
		attempt       		INT(11) default 0 ,
		createdAt       	DATETIME,
		updatedAt      		DATETIME
	);

	ALTER TABLE delivery 
		ADD CONSTRAINT delivererIdFK
		  FOREIGN KEY (delivererId)
		  REFERENCES deliverer (id)
		  ON DELETE NO ACTION
		  ON UPDATE NO ACTION,
		ADD CONSTRAINT notificationIdFK
		  FOREIGN KEY (notificationId)
		  REFERENCES notification (id)
		  ON DELETE NO ACTION
		  ON UPDATE NO ACTION;
	`
)
