CREATE TABLE songs
(
    id SERIAL PRIMARY KEY,
    song            VARCHAR(255),
    group_name      VARCHAR(255),
    release_date    VARCHAR(12),
    text            TEXT, 
    link            VARCHAR(255)
);
