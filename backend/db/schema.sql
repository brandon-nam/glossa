-- Bill Schema 
CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    bill_id INT UNIQUE,
    name TEXT,
    proposers TEXT,
    main_text TEXT,
    summary TEXT,
    categories TEXT,
    detail_url TEXT
);