CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(25) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    credit_token INT DEFAULT 15,
    last_first_llm_used TIMESTAMP DEFAULT NULL
    multifa_enabled BOOLEAN DEFAULT FALSE,
    phone_number VARCHAR(20) DEFAULT NULL,
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

CREATE TYPE otp_purpose AS ENUM ('login_2fa', 'password_reset', 'phone_verification');
CREATE TYPE otp_channel AS ENUM ('whatsapp', 'sms', 'email');

CREATE TABLE user_otps (
    id            SERIAL        PRIMARY KEY,
    user_id       INT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose       otp_purpose NOT NULL,
    channel       otp_channel NOT NULL,
    code_hash     TEXT        NOT NULL,
    created_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at    TIMESTAMP   NOT NULL,
    used          BOOLEAN     NOT NULL DEFAULT FALSE,
    used_at       TIMESTAMP
);

-- Index to quickly find the latest unused code for a user+purpose
CREATE INDEX idx_user_otps_user_purpose 
ON user_otps(user_id, purpose, used, expires_at);

-- (Optional) background job can DELETE FROM user_otps WHERE expires_at < NOW() OR used = TRUE;
