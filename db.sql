CREATE TABLE operations (
  id serial primary key,
  merchant_id varchar(20),
  application_name varchar(100),
  operation_name varchar(100),
  operation_ident varchar(100),
  description text,
  amount integer,
  xml_data xml,
  operation_created_at timestamp without time zone,
  operation_date_created_at date,
  created_at timestamp without time zone default now(),
  updated_at timestamp without time zone
);

CREATE UNIQUE INDEX index_operations_on_application_name_and_operation_name_and_operation_ident ON operations USING btree
  (application_name COLLATE pg_catalog."default", operation_name COLLATE pg_catalog."default", operation_ident COLLATE pg_catalog."default");

CREATE INDEX index_operations_on_operation_name ON operations USING btree
  (operation_name COLLATE pg_catalog."default");

GRANT ALL PRIVILEGES ON DATABASE eticket_billing_server_development TO eticket_billing_server_user;
GRANT ALL PRIVILEGES ON TABLE operations TO eticket_billing_server_user;
GRANT ALL PRIVILEGES ON SEQUENCE operations_id_seq TO eticket_billing_server_user;

CREATE OR REPLACE FUNCTION operation_insert() RETURNS TRIGGER AS '
  BEGIN
    NEW.operation_date_created_at := date(NEW.operation_created_at);
    RETURN NEW;
  END
' LANGUAGE plpgsql;

CREATE TRIGGER operation_insert BEFORE INSERT OR UPDATE ON operations FOR EACH ROW EXECUTE PROCEDURE operation_insert();
