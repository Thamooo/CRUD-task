CREATE TABLE public.clients (
    id integer NOT NULL,
    first_name text,
    last_name text,
    birth_date date,
    gender text,
    email text,
    address text
);
ALTER TABLE public.clients ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.clients_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);