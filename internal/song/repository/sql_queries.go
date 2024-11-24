package repository

const getList = `SELECT id, artist, title, text, link FROM songs ORDER BY ? ?  LIMIT ? OFFSET ?`

const getText = `SELECT text FROM songs WHERE title = ? AND artist = ?`

const deleteSong = `DELETE FROM songs WHERE id = ?`

const updateSong = `UPDATE songs SET artist = ?, title = ?, text = ?, link = ? WHERE id = ?`

const createSong = `INSERT INTO songs (artist, title, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`
