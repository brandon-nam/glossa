-- National Assembly Bill Schema 
CREATE TABLE IF NOT EXISTS assembly_bill (
    id SERIAL PRIMARY KEY,
    bill_id INT UNIQUE,
    name TEXT,
    proposers TEXT,
    department TEXT,
	parliamentary_status TEXT,
	resolution_status TEXT,
    main_text TEXT,
    summary TEXT,
    categories TEXT,
    detail_url TEXT
);

-- Government Bill Schema 
CREATE TABLE IF NOT EXISTS government_bill (
    id SERIAL PRIMARY KEY,
    bill_id INT UNIQUE,
    name TEXT,
    proposers TEXT,
    department TEXT,
	parliamentary_status TEXT,
	resolution_status TEXT,
    main_text TEXT,
    summary TEXT,
    categories TEXT,
    detail_url TEXT
);