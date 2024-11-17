-- migrations/202311050001_create_capsules_table.up.sql
CREATE TABLE capsules (
    id SERIAL PRIMARY KEY,
    sender_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    send_at TIMESTAMP NOT NULL,
    message VARCHAR(4096) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    recipient_tg_username VARCHAR(50),
    files_folder_UUID UUID
);
