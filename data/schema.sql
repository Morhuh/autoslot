PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT,
    email TEXT,
    website TEXT,
    description TEXT NOT NULL,
    specialties TEXT,
    working_hours TEXT,
    lat REAL DEFAULT 0,
    lng REAL DEFAULT 0,
    rating REAL DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    featured INTEGER DEFAULT 0,
    image_url TEXT
);

CREATE TABLE IF NOT EXISTS service_offerings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    price_from TEXT,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_brands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    brand TEXT NOT NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_amenities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    amenity TEXT NOT NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_gallery (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    author_name TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT NOT NULL,
    rating INTEGER NOT NULL,
    speed_rating INTEGER NOT NULL,
    price_rating INTEGER NOT NULL,
    quality_rating INTEGER NOT NULL,
    kindness_rating INTEGER NOT NULL,
    approved INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    excerpt TEXT NOT NULL,
    category TEXT NOT NULL,
    read_time TEXT NOT NULL,
    image_url TEXT,
    body TEXT NOT NULL,
    featured INTEGER NOT NULL DEFAULT 0,
    published_at DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    brand TEXT NOT NULL,
    model_name TEXT NOT NULL,
    years TEXT NOT NULL,
    engine TEXT NOT NULL,
    service_interval TEXT NOT NULL,
    known_issues TEXT NOT NULL,
    typical_costs TEXT NOT NULL,
    summary TEXT NOT NULL,
    image_url TEXT,
    featured INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    customer_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    vehicle_brand TEXT,
    vehicle_model TEXT,
    vehicle_year TEXT,
    requested_date TEXT NOT NULL,
    service_request TEXT NOT NULL,
    issue_details TEXT,
    status TEXT NOT NULL DEFAULT 'new',
    selected_slot_id INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS availability_slots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    starts_at TEXT NOT NULL,
    note TEXT,
    active INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_user_id) REFERENCES service_users(id) ON DELETE CASCADE
);
