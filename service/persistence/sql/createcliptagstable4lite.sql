CREATE TABLE IF NOT EXISTS cliptags (
		id INTEGER NOT NULL PRIMARY KEY,
        clipId TEXT NOT NULL,
        tag TEXT NOT NULL,
        createdTime DATETIME NOT NULL
    );	
