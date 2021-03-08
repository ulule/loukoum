CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    deleted_at timestamp with time zone,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    is_staff boolean DEFAULT false NOT NULL
);

INSERT INTO users (email, first_name, last_name, is_staff)
    VALUES ('tech@ulule.com', 'Ulule', 'Tech', true);

INSERT INTO users (email, first_name, last_name, is_staff)
    VALUES ('thomas.leroux@ulule.com', 'Thomas', 'LE ROUX', true);

CREATE TABLE news (
    id SERIAL NOT NULL PRIMARY KEY,
    published_at timestamp with time zone,
    deleted_at timestamp with time zone,
    status character varying(255) NOT NULL
);

INSERT INTO news (status) VALUES ('draft');

CREATE TABLE comments (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    user_id integer NOT NULL,
    message character varying(2048) NOT NULL
);

ALTER TABLE ONLY comments
    ADD CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY comments
    ADD CONSTRAINT comments_email_idx UNIQUE (email);

INSERT INTO comments (email, user_id, status, message, created_at)
    VALUES ('thomas.leroux@ulule.com', 2, 'waiting', 'Hello world', NOW());
