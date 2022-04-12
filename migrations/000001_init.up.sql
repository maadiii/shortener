BEGIN;
	CREATE SCHEMA IF NOT EXISTS "shortener";
	CREATE TABLE IF NOT EXISTS "shortener"."link"(
		id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		main VARCHAR(128) NOT NULL,
		shorted VARCHAR(9) NOT NULL,

		CONSTRAINT non_empty_main CHECK(main <> ''),
		CONSTRAINT non_empty_shorted CHECK(shorted <> '')
	);
COMMIT;
