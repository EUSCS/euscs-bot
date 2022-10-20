package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var migrations = []string{
	migration0,
	migration1,
	migration2,
}

func migrate() error {
	var start int
	_, err := db.Exec(migrations[0])
	if err != nil {
		return err
	}
	start, err = getLatestMigration()
	if err != nil {
		return err
	}
	for i := start + 1; i < len(migrations); i++ {
		log.Info(fmt.Sprintf("applying migration %d", i))
		_, err = db.Exec(migrations[i])
		if err != nil {
			return err
		}
		_, err = db.Exec(`INSERT INTO migrations (version) VALUES (?)`, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func getLatestMigration() (int, error) {
	ver := 0
	row := db.QueryRow(`SELECT version
	FROM migrations
	ORDER BY version DESC
	LIMIT 1`)
	err := row.Scan(&ver)
	return ver, err
}

var migration0 = `CREATE TABLE IF NOT EXISTS migrations (
	version int
);
INSERT INTO migrations (version) VALUES (0);
`

var migration1 = `CREATE TABLE players (
    discordID text UNIQUE,
		elo int DEFAULT 1500 NOT NULL,
		role text DEFAULT "" NOT NULL,
		queuing int DEFAULT 0 NOT NULL,
		osuser text DEFAULT "",
		PRIMARY KEY(discordID)
);`

var migration2 = `CREATE TABLE matches (
	matchID text UNIQUE,
	messageID text UNIQUE,
	threadID text UNIQUE,
	timestamp int,
	running int DEFAULT 1 NOT NULL,
	team1score int DEFAULT 0 NOT NULL,
	team2score int DEFAULT 0 NOT NULL,
	PRIMARY KEY(matchID)
);
CREATE TABLE matchesplayers (
	matchID text,
	team int,
	playerID text,
	FOREIGN KEY (playerID) REFERENCES players(discordID),
	FOREIGN KEY (matchID) REFERENCES matches(matchID),
	PRIMARY KEY (playerID,matchID)
);
`