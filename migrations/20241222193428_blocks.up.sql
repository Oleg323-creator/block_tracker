CREATE TABLE blocks (
                        id SERIAL PRIMARY KEY,
                        block_number BIGINT
);

INSERT INTO blocks (id, block_number) VALUES (1, 0);