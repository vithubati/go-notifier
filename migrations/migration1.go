package migrations

const (
	migration1 = `
	CREATE TABLE IF NOT EXISTS notification
	(
		id 				varchar(36) PRIMARY KEY,
		topic           VARCHAR(50) NOT NULL,
		action          VARCHAR(50) NOT NULL,
		subject         VARCHAR(50) NOT NULL,
		message         VARCHAR(256) NOT NULL,
		createdAt       DATETIME DEFAULT NULL,
		data            TEXT
	);
`
	migration2 = `
	CREATE TABLE IF NOT EXISTS deliverer
	(
		id 					varchar(36) PRIMARY KEY,
		type        		VARCHAR(50) NOT NULL,
		url          		VARCHAR(50) default '',
		channelId          	VARCHAR(50) default '',
		headers         	VARCHAR(50),
		credentials     	VARCHAR(256) default '',
		createdAt       	DATETIME,
		retry       		INT(11),
		intervalInSeconds   INT(11)
	);
`
	migration3 = `
	CREATE TABLE IF NOT EXISTS deliverer_topic
	(
		deliverer_id 		varchar(36) NOT NULL,
		topic 			varchar(36) NOT NULL
	);
`
	migration4 = `
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
`
	migration5 = `
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
	migration6 = `
	ALTER TABLE deliverer_topic
		ADD INDEX delevererFK_idx (deliverer_id ASC);
`
	migration7 = `
	ALTER TABLE deliverer_topic 
		ADD CONSTRAINT delevererFK
		  FOREIGN KEY (deliverer_id)
		  REFERENCES deliverer (id)
		  ON DELETE CASCADE
		  ON UPDATE NO ACTION;
`
)
