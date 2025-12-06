-- Utworzenie tabeli customers
CREATE TABLE customers (
    email TEXT PRIMARY KEY,
    stripe_customer_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Usunięcie stripe_customer_id z forms
ALTER TABLE forms DROP COLUMN stripe_customer_id;

-- Dodanie foreign key do customers
ALTER TABLE forms ADD CONSTRAINT fk_forms_email
    FOREIGN KEY (email) REFERENCES customers(email) ON DELETE CASCADE;
