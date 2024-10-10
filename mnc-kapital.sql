--
-- PostgreSQL database dump
--

-- Dumped from database version 13.16 (Debian 13.16-1.pgdg120+1)
-- Dumped by pg_dump version 13.16 (Debian 13.16-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: transaction_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_type AS ENUM (
    'TOPUP',
    'PAYMENT',
    'TRANSFER'
);


ALTER TYPE public.transaction_type OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL,
    applied_at timestamp with time zone NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    id uuid NOT NULL,
    user_id uuid,
    transaction_type public.transaction_type NOT NULL,
    amount bigint NOT NULL,
    remarks character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    phone_number character varying(255) NOT NULL,
    pin character varying(255) NOT NULL,
    address character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: view_transactions; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.view_transactions AS
 SELECT t.id,
    t.user_id,
    t.transaction_type,
        CASE
            WHEN (t.amount >= 0) THEN 'CREDIT'::text
            WHEN (t.amount < 0) THEN 'DEBIT'::text
            ELSE NULL::text
        END AS account_type,
    t.amount,
    t.remarks,
    t.created_at,
    ( SELECT COALESCE(sum(t_prev.amount), (0)::numeric) AS "coalesce"
           FROM public.transactions t_prev
          WHERE ((t_prev.user_id = t.user_id) AND (t_prev.created_at < t.created_at))) AS balance_before,
    ( SELECT COALESCE(sum(t_prev.amount), (0)::numeric) AS "coalesce"
           FROM public.transactions t_prev
          WHERE ((t_prev.user_id = t.user_id) AND (t_prev.created_at <= t.created_at))) AS balance_after
   FROM public.transactions t
  ORDER BY t.created_at;


ALTER TABLE public.view_transactions OWNER TO postgres;

--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, applied_at) FROM stdin;
0000	2024-10-10 16:22:43.616762+00
0001	2024-10-10 17:03:03.57117+00
0002	2024-10-10 19:48:49.611828+00
0003	2024-10-10 20:38:43.002269+00
0004	2024-10-10 21:17:15.879221+00
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (id, user_id, transaction_type, amount, remarks, created_at) FROM stdin;
11a34ee7-2cbc-47a0-a053-7bc5fb3e03d9	9ae4eede-7c8f-4807-a276-f27dd718d258	TOPUP	50000		2024-10-10 20:38:46.357696+00
53deeed2-8b89-4e52-b1d2-2573cf7de622	9ae4eede-7c8f-4807-a276-f27dd718d258	TOPUP	50000		2024-10-10 20:40:47.076634+00
d539a8c8-b486-490b-bbf3-60d6632f6971	9ae4eede-7c8f-4807-a276-f27dd718d258	TOPUP	50000		2024-10-10 20:43:49.403142+00
7a3dadea-38a9-4143-b471-5f6225d4bcf0	9ae4eede-7c8f-4807-a276-f27dd718d258	TOPUP	50000		2024-10-10 21:00:41.231632+00
d9fcec1f-f230-40e9-9e52-8e28a51d3e1c	9ae4eede-7c8f-4807-a276-f27dd718d258	PAYMENT	-50000	Pembayaran Voucher Pulsa	2024-10-10 21:01:55.537683+00
61d513c1-21aa-48a1-8d71-6fbeaa7919ec	9ae4eede-7c8f-4807-a276-f27dd718d258	TRANSFER	-50000	Hadiah Ulang Tahun	2024-10-10 21:09:01.158122+00
be94ef5c-d4ac-4471-a208-72dc2d03e932	9dfa9b39-5c5a-4bde-b541-c5ebec9431da	TRANSFER	50000	Hadiah Ulang Tahun	2024-10-10 21:09:01.162949+00
d0c8a025-d562-4850-b0ee-ac48a63cb054	9ae4eede-7c8f-4807-a276-f27dd718d258	TOPUP	50000		2024-10-10 21:31:17.749466+00
3b5fe92b-dc71-42f2-9a90-b7fa47a81841	9ae4eede-7c8f-4807-a276-f27dd718d258	PAYMENT	-50000	Pembayaran Voucher Pulsa	2024-10-10 21:31:27.725131+00
7e19c10d-61ed-4e5b-956c-2ea61a3a1a5a	9ae4eede-7c8f-4807-a276-f27dd718d258	PAYMENT	-50000	Pembayaran Voucher Pulsa	2024-10-10 21:31:32.4516+00
0e22810d-3df6-434f-9cfd-ccfd4ade9693	9ae4eede-7c8f-4807-a276-f27dd718d258	PAYMENT	-50000	Pembayaran Voucher Pulsa	2024-10-10 21:31:33.573984+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, first_name, last_name, phone_number, pin, address, created_at, updated_at) FROM stdin;
9ae4eede-7c8f-4807-a276-f27dd718d258	Tom	Araya	0811255501	$2a$10$cHKxOoIbEdlQBj7hk4gxdO1HAIIJS6OErJPQY0yG.zEWjmC.pjt2O	Jl. Diponegoro No. 215	2024-10-10 19:04:06.398756+00	2024-10-10 20:18:45.145127+00
9dfa9b39-5c5a-4bde-b541-c5ebec9431da	Muhammad	Arya Dyas	0895613367705	$2a$10$xtvTe8bBxgwWBQ0cBNuHM.ZJmyIrHnKBvbC6yAdGuw02l9TWeSYSe	Jl. Kebon Sirih No. 1	2024-10-10 21:08:03.885469+00	2024-10-10 21:08:03.885469+00
\.


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: transactions_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX transactions_user_id_idx ON public.transactions USING btree (user_id);


--
-- Name: transactions transactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

