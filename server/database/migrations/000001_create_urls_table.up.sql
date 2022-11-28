CREATE TABLE urls (
  "domain" TEXT NOT NULL,
  "path" TEXT NOT NULL,
  "target" TEXT NOT NULL,
  "comment" TEXT,
  "id" INTEGER PRIMARY KEY,
  UNIQUE("domain", "path")
);
