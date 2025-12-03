CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    -- email właściciela formularza, który opłaca usługę obsługi formularz
    stripe_customer_id TEXT,  -- cus_123 (ten sam dla wszystkich formularzy usera)
    stripe_subscription_id TEXT UNIQUE,  -- sub_111, sub_222 (unikalny dla każdego formularza)
    is_paid BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
