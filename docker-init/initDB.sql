-- CREATE SEQUENCE user_id_seq;
CREATE TABLE "user" (
    user_id BIGINT NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    chat_id BIGINT DEFAULT 0,
    PRIMARY KEY (user_id)
);

-- ALTER SEQUENCE user_id_seq OWNED BY user.user_id;
CREATE TABLE chat (
    chat_id BIGINT NOT NULL,
    poll VARCHAR(255),
    PRIMARY KEY (chat_id)
);

CREATE TABLE user_chat (
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    is_head BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES "user" (user_id),
    FOREIGN KEY (chat_id) REFERENCES chat (chat_id)
);