
CREATE TABLE IF NOT EXISTS "cities" (
                          "city_id"	INTEGER NOT NULL UNIQUE,
                          "title_ru"	TEXT,
                          PRIMARY KEY("city_id")
);

CREATE TABLE IF NOT EXISTS "users" (
                         "id"	INTEGER NOT NULL UNIQUE,
                         "discord_id"	TEXT NOT NULL UNIQUE,
                         "score"	INTEGER DEFAULT 0,
                         PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE INDEX "idx_cities_title_ru" ON "cities" ("title_ru");
