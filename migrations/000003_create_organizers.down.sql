-- Drop trigger before dropping the table
DROP TRIGGER IF EXISTS set_updated_at ON organizers;

-- Drop table
DROP TABLE IF EXISTS organizers CASCADE;

DROP TYPE IF EXISTS organization_type_enum;
