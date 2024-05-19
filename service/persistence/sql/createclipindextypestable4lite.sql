CREATE TABLE IF NOT EXISTS clipindextypes (
		id INTEGER NOT NULL PRIMARY KEY,
        clipId TEXT NOT NULL,
        indexType TEXT NOT NULL,
        createdTime DATETIME NOT NULL
    );	
