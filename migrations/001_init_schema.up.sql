-- Drop existing tables if they exist
DROP TABLE IF EXISTS "booking_units";
DROP TABLE IF EXISTS "bookings";
DROP TABLE IF EXISTS "availabilities";
DROP TABLE IF EXISTS "products";

-- Create new tables
CREATE TABLE "products" (
    "id" VARCHAR(255) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "capacity" INT NOT NULL,
    "price" REAL NOT NULL,
    "currency" VARCHAR(50) NOT NULL
);

CREATE TABLE "availabilities" (
    "id" VARCHAR(255) PRIMARY KEY,
    "local_date" DATE NOT NULL,
    "status" VARCHAR(50) NOT NULL,
    "product_id" VARCHAR(50) NOT NULL,
    "vacancies" INT NOT NULL,
    "available" BOOLEAN NOT NULL,
    "price" REAL NOT NULL,
    "currency" VARCHAR(50) NOT NULL,
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE TABLE "bookings" (
    "id" VARCHAR(255) PRIMARY KEY,
    "status" VARCHAR(50) NOT NULL,
    "availability_id" VARCHAR(255) NOT NULL,
    "units" INT NOT NULL,
    "price" REAL NOT NULL,
    "currency" VARCHAR(50) NOT NULL,
    FOREIGN KEY ("availability_id") REFERENCES "availabilities" ("id")
);

CREATE TABLE "booking_units" (
    "id" VARCHAR(255) PRIMARY KEY,
    "booking_id" VARCHAR(255) NOT NULL,
    "price" REAL NOT NULL,
    "currency" VARCHAR(50) NOT NULL,
    FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id")
);