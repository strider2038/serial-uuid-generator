CREATE TABLE public.uuid_sequence
(
  sequence VARCHAR (32) PRIMARY KEY,
  start_value BIGINT NOT NULL DEFAULT 0
);
