BEGIN;

CREATE TABLE IF NOT EXISTS public.companies
(
    mse_code character varying(10) COLLATE pg_catalog."default" NOT NULL,
    name character varying(150) COLLATE pg_catalog."default",
    CONSTRAINT companies_pkey PRIMARY KEY (mse_code)
)

CREATE TABLE IF NOT EXISTS public.daily_company_rates
(
    id character varying COLLATE pg_catalog."default" NOT NULL,
    no character varying COLLATE pg_catalog."default",
    high character varying COLLATE pg_catalog."default",
    low character varying COLLATE pg_catalog."default",
    code character varying(10) COLLATE pg_catalog."default" NOT NULL,
    buy numeric,
    sell numeric,
    pcp numeric,
    tcp numeric,
    vol integer,
    div_net numeric,
    div_yield numeric,
    earn_yield numeric,
    pe_ratio numeric,
    pbv_ratio numeric,
    cap money,
    profit money,
    shares bigint,
    date date NOT NULL,
    CONSTRAINT daily_company_rates_pkey PRIMARY KEY (id)
)

END;