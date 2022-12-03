-- +Up
CREATE SCHEMA IF NOT EXISTS  crypto;
CREATE TABLE IF NOT EXISTS crypto.quotes (
                             fsyms VARCHAR(255) NOT NULL,
                             tsyms VARCHAR(255) NOT NULL,
                             constraint fsyms_tsyms  PRIMARY KEY (fsyms,tsyms),
                             raw jsonb,
                             display jsonb,
                             error_message VARCHAR(255) default NULL,
                             created_at timestamp default now(),
                             updated_at timestamp default now()
);


