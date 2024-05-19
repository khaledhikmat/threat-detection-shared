CREATE TABLE IF NOT EXISTS clipalerttypes (
		id INTEGER NOT NULL PRIMARY KEY,
        clipId TEXT NOT NULL,
        alertType TEXT NOT NULL,
        createdTime DATETIME NOT NULL
    );	
