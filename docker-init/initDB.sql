-- CREATE SEQUENCE user_id_seq;
CREATE TABLE "user" (
    user_id INTEGER NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    chat_id INTEGER DEFAULT 0,
    PRIMARY KEY (user_id)
);
-- ALTER SEQUENCE user_id_seq OWNED BY user.user_id;

CREATE TABLE chat (
    chat_id INTEGER NOT NULL,
    head_person INTEGER,
    poll INTEGER,
    PRIMARY KEY (chat_id)
);

CREATE TABLE user_chat (
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id)
        REFERENCES "user" (user_id),
    FOREIGN KEY (chat_id)
        REFERENCES chat (chat_id)
);