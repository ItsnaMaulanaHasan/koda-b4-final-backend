CREATE TABLE "clicks" (
    "id" serial PRIMARY KEY,
    "short_link_id" int NOT NULL,
    "clicked_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "ip_address" varchar(45),
    "referer" text,
    "user_agent" text,
    "country" varchar(100),
    "city" varchar(100),
    "device_type" varchar(20),
    "browser" varchar(50),
    "os" varchar(50),
    "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
);

ALTER TABLE "clicks"
ADD FOREIGN KEY ("short_link_id") REFERENCES "short_links" ("id") ON DELETE CASCADE;

ALTER TABLE "clicks"
ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "clicks"
ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

CREATE INDEX idx_clicks_short_link_id_clicked_at ON "clicks" ("short_link_id", "clicked_at");

CREATE INDEX idx_clicks_clicked_at ON "clicks" ("clicked_at");

CREATE INDEX idx_clicks_country ON "clicks" ("country");

CREATE INDEX idx_clicks_device_type ON "clicks" ("device_type");