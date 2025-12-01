CREATE TABLE "short_links" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "short_code" varchar(20) UNIQUE NOT NULL,
    "original_url" text NOT NULL,
    "title" varchar(255),
    "is_active" bool DEFAULT true,
    "expired_at" timestamp,
    "click_count" int DEFAULT 0,
    "last_clicked_at" timestamp,
    "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "created_by" int,
    "updated_by" int
);

-- Foreign keys
ALTER TABLE "short_links"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "short_links"
ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "short_links"
ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

-- Indexes
CREATE UNIQUE INDEX idx_short_links_short_code ON "short_links" ("short_code");

CREATE INDEX idx_short_links_user_id_created_at ON "short_links" ("user_id", "created_at");

CREATE INDEX idx_short_links_is_active ON "short_links" ("is_active")
WHERE
    "is_active" = true;