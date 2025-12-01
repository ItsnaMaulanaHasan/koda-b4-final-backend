CREATE TABLE "sessions" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "refresh_token" text NOT NULL,
    "login_time" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "logout_time" timestamp,
    "expired_at" timestamp NOT NULL,
    "ip_address" varchar(45),
    "user_agent" text,
    "is_active" bool DEFAULT true,
    "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "created_by" int,
    "updated_by" int
);

ALTER TABLE "sessions"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions"
ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "sessions"
ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

CREATE INDEX idx_sessions_user_id ON "sessions" ("user_id");

CREATE INDEX idx_sessions_refresh_token ON "sessions" ("refresh_token");

CREATE INDEX idx_sessions_is_active ON "sessions" ("is_active")
WHERE
    "is_active" = true;