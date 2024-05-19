CREATE TABLE IF NOT EXISTS clips (
		id TEXT NOT NULL PRIMARY KEY,
        localReference TEXT NOT NULL,
        cloudReference TEXT NOT NULL,
        storageProvider
        capturer TEXT NOT NULL,
        camera TEXT NOT NULL,
        region TEXT NOT NULL,
        location TEXT NOT NULL,
        priority TEXT NOT NULL,
        frames INTEGER NOT NULL,
        beginTime DATETIME NOT NULL,
        endTime DATETIME NOT NULL,
        prevClip TEXT,
        analytics TEXT,
        alertTypes TEXT,
        mediaIndexerTypes TEXT,
        tags TEXT
    );	