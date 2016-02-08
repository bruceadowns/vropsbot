package common

const (
	// SlackItemTypeUser is the db constant for a user type
	SlackItemTypeUser = 1

	// SlackItemTypeGroup is the db constant for a group type
	SlackItemTypeGroup = 2

	// SlackItemTypeChannel is the db constant for a channel type
	SlackItemTypeChannel = 3
)

var createDbSQL = `

--
-- sqlite3 initialization database script
--

-- generic key / value store

CREATE TABLE IF NOT EXISTS 'KVStore' (
  'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
  , 'Key' TEXT
  , 'Value' TEXT
);

-- plugin heartbeat

CREATE TABLE IF NOT EXISTS 'HeartBeat' (
	'Name' TEXT NOT NULL PRIMARY KEY UNIQUE
	, 'When' TEXT
);

-- slack user/group/channel cache

CREATE TABLE IF NOT EXISTS 'SlackItemType' (
  'ID' INTEGER NOT NULL PRIMARY KEY UNIQUE
  , 'Name' TEXT
);

CREATE TABLE IF NOT EXISTS 'SlackItem' (
  'ID' TEXT NOT NULL PRIMARY KEY UNIQUE
  , 'Name' TEXT
	, 'Type' INTEGER
  , FOREIGN KEY('Type') REFERENCES SlackItemType(ID)
);

-- pull and push actions

CREATE TABLE IF NOT EXISTS 'PullActions' (
	'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
  , 'When' TEXT
	, 'Channel' TEXT
	, 'User' TEXT
	, 'Request' TEXT
  , 'Response' TEXT
  , 'IsError' INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS 'PushActions' (
	'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
  , 'When' TEXT
  , 'PluginName' TEXT
	, 'Channel' TEXT
	, 'Request' TEXT
  , 'Response' TEXT
  , 'Before' TEXT
  , 'After' TEXT
  , 'IsError' INTEGER DEFAULT 0
);

-- seed cross reference table for slack item types

INSERT INTO 'SlackItemType' ( 'ID', 'Name' )
VALUES
( 1, 'User' ),
( 2, 'Group' ),
( 3, 'Channel' );

`
