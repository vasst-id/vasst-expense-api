-- Drop indexes (if not already dropped by cascade)
DROP INDEX IF EXISTS "vasst_ca".idx_message_organization_id;
DROP INDEX IF EXISTS "vasst_ca".idx_message_conversation_id;
DROP INDEX IF EXISTS "vasst_ca".idx_conversation_organization_id;
DROP INDEX IF EXISTS "vasst_ca".idx_conversation_user_id;
DROP INDEX IF EXISTS "vasst_ca".idx_conversation_contact_id;
DROP INDEX IF EXISTS "vasst_ca".idx_conversation_medium_id;
DROP INDEX IF EXISTS "vasst_ca".idx_contact_organization_id;
DROP INDEX IF EXISTS "vasst_ca".idx_contact_phone_number;
DROP INDEX IF EXISTS "vasst_ca".idx_user_organization_id;
DROP INDEX IF EXISTS "vasst_ca".idx_user_phone_number; -- user
DROP INDEX IF EXISTS "vasst_ca".idx_user_phone_number; -- tags

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS "vasst_ca".messages;
DROP TABLE IF EXISTS "vasst_ca".conversation;
DROP TABLE IF EXISTS "vasst_ca".contact_tags;
DROP TABLE IF EXISTS "vasst_ca".contact;
DROP TABLE IF EXISTS "vasst_ca".tags;
DROP TABLE IF EXISTS "vasst_ca".user;
DROP TABLE IF EXISTS "vasst_ca".organization_setting;
DROP TABLE IF EXISTS "vasst_ca".organization_knowledge;
DROP TABLE IF EXISTS "vasst_ca".organization_medium;
DROP TABLE IF EXISTS "vasst_ca".medium;
DROP TABLE IF EXISTS "vasst_ca".message_type;
DROP TABLE IF EXISTS "vasst_ca".role;
DROP TABLE IF EXISTS "vasst_ca".plan;
DROP TABLE IF EXISTS "vasst_ca".organization_category;
DROP TABLE IF EXISTS "vasst_ca".organization;

-- Drop schema
DROP SCHEMA IF EXISTS "vasst_ca" CASCADE;

-- Drop extension (optional, only if you want to remove uuid-ossp for the whole DB)
DROP EXTENSION IF EXISTS "uuid-ossp";