CREATE TABLE IF NOT EXISTS clipanalytics (
		id INTEGER NOT NULL PRIMARY KEY,
        clipId TEXT NOT NULL,
        analytic TEXT NOT NULL,
        createdTime DATETIME NOT NULL
    );	
