CREATE TABLE IF NOT EXISTS definitions (
   id SERIAL PRIMARY KEY,
   word TEXT not null,
   definition TEXT not null
);


INSERT INTO definitions (word, definition) VALUES ('rawr', 'To roar like a lion');

