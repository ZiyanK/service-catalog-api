CREATE TABLE "users" (
  "user_uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" VARCHAR(50) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "services" (
  "service_id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "user_uuid" UUID NOT NULL,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "service_versions" (
  "sv_id" SERIAL PRIMARY KEY,
  "version" VARCHAR(8),
  "changelog" TEXT,
  "service_id" INTEGER NOT NULL,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "services" ADD FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid");

ALTER TABLE "service_versions" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("service_id");

-- Insert test user for tests
INSERT INTO users(user_uuid, email, password) VALUES ('d90f9b49-dcd9-4feb-8250-d013098e45ee', 'test@gmail.com', '$2a$04$YnrU9OJGE8ywjQ9yxIDZguyyRJucS4a5doOFk/3cde4OXxos6AkO.');