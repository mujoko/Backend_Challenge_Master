CREATE TABLE stock
(
 id serial NOT NULL,
 name character varying NOT NULL,
 price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
 availability integer,
 is_active boolean, 
 CONSTRAINT pk_buku PRIMARY KEY (id )
)
WITH (
 OIDS=FALSE
);
ALTER TABLE stock
 OWNER TO postgres;
