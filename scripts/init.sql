CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS encounters;
DROP TABLE IF EXISTS providers;
DROP TABLE IF EXISTS patients;

CREATE TABLE IF NOT EXISTS patients (
  id TEXT PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXISTS providers (
  id TEXT PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXISTS encounters (
  id TEXT PRIMARY KEY,
  patient_id TEXT REFERENCES patients(id),
  provider_id TEXT REFERENCES providers(id),
  service_date TIMESTAMP
);
