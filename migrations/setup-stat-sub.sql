ALTER DATABASE stat_db SET timezone TO 'Europe/Moscow';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/*AI SUB*/

CREATE SUBSCRIPTION ai_sub CONNECTION 'port=5432 user=ai dbname=ai_db host=db-ai password=aipass' PUBLICATION ai_pub;

/*USERS SUB*/

CREATE SUBSCRIPTION users_sub CONNECTION 'port=5432 user=user dbname=users_db host=db-users password=userspass' PUBLICATION users_pub;
