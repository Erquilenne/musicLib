package repository

const getList = `SELECT id, group_name, song, text, link FROM songs`

const getText = `SELECT text FROM songs WHERE id = $1`

const deleteSong = `DELETE FROM songs WHERE id = $1`

const updateSong = `
    UPDATE songs 
    SET group_name = $1, 
        song = $2, 
        text = $3, 
        link = $4,
        release_date = $5 
    WHERE id = $6
    RETURNING id`

const createSong = `INSERT INTO songs (group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`
