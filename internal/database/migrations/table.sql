CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(25) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    credit_token INT DEFAULT 15,
    last_first_llm_used TIMESTAMP DEFAULT NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,  -- Bisa pakai UUID
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE TYPE credit_status AS ENUM ('pending', 'confirmed', 'cancelled');
CREATE TYPE feature_type AS ENUM ('chat-ai');

CREATE TABLE user_credit_reserved (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credit INT NOT NULL,
    feature_type feature_type NOT NULL,
    status credit_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE room_credit_reserved_conn (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(255) NOT NULL REFERENCES room_chat_train(room_code) ON DELETE CASCADE,
    user_credit_reserved_id INT NOT NULL REFERENCES user_credit_reserved(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);
