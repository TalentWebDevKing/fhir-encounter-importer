CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS patients (
  id UUID PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXISTS providers (
  id UUID PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXISTS encounters (
  id UUID PRIMARY KEY,
  patient_id UUID REFERENCES patients(id),
  provider_id UUID REFERENCES providers(id),
  service_date TIMESTAMP
);
