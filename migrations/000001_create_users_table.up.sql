CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "fullName" varchar(255) NOT NULL,
    "email" varchar(255) UNIQUE NOT NULL,
    "role" varchar(20) NOT NULL DEFAULT 'customer',
    "password" text NOT NULL,
    "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
    "created_by" int,
    "updated_by" int
);

ALTER TABLE "users"
ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "users"
ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

CREATE INDEX idx_users_email ON "users" ("email");