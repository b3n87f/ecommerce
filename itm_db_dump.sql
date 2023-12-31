PGDMP      2            
    {            itm    16.0    16.0 6    )           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            *           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            +           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            ,           1262    16397    itm    DATABASE     w   CREATE DATABASE itm WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Turkish_Turkey.1254';
    DROP DATABASE itm;
                postgres    false            �            1259    16414 	   customers    TABLE       CREATE TABLE public.customers (
    first_name text NOT NULL,
    last_name text NOT NULL,
    phone text NOT NULL,
    address text NOT NULL,
    email text NOT NULL,
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
    DROP TABLE public.customers;
       public         heap    postgres    false            �            1259    16433    customers_id_seq    SEQUENCE     y   CREATE SEQUENCE public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 '   DROP SEQUENCE public.customers_id_seq;
       public          postgres    false    216            -           0    0    customers_id_seq    SEQUENCE OWNED BY     E   ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;
          public          postgres    false    218            �            1259    25031    loggers    TABLE     O   CREATE TABLE public.loggers (
    log text NOT NULL,
    id bigint NOT NULL
);
    DROP TABLE public.loggers;
       public         heap    postgres    false            �            1259    25030    loggers_id_seq    SEQUENCE     w   CREATE SEQUENCE public.loggers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 %   DROP SEQUENCE public.loggers_id_seq;
       public          postgres    false    226            .           0    0    loggers_id_seq    SEQUENCE OWNED BY     A   ALTER SEQUENCE public.loggers_id_seq OWNED BY public.loggers.id;
          public          postgres    false    225            �            1259    16538    order_details    TABLE     N  CREATE TABLE public.order_details (
    order_id text NOT NULL,
    product_id bigint NOT NULL,
    product_base_price bigint NOT NULL,
    product_discount bigint NOT NULL,
    product_pay_price bigint NOT NULL,
    product_vat bigint NOT NULL,
    product_name text NOT NULL,
    quantity bigint NOT NULL,
    id bigint NOT NULL
);
 !   DROP TABLE public.order_details;
       public         heap    postgres    false            �            1259    16559    order_details_id_seq    SEQUENCE     }   CREATE SEQUENCE public.order_details_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.order_details_id_seq;
       public          postgres    false    222            /           0    0    order_details_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.order_details_id_seq OWNED BY public.order_details.id;
          public          postgres    false    224            �            1259    16510    orders    TABLE     m  CREATE TABLE public.orders (
    customer bigint NOT NULL,
    inserted_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_date timestamp without time zone,
    status bigint NOT NULL,
    order_id text NOT NULL,
    transaction_method text NOT NULL,
    total_price bigint NOT NULL,
    id bigint NOT NULL,
    products text NOT NULL
);
    DROP TABLE public.orders;
       public         heap    postgres    false            �            1259    16548    orders_id_seq    SEQUENCE     v   CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 $   DROP SEQUENCE public.orders_id_seq;
       public          postgres    false    221            0           0    0    orders_id_seq    SEQUENCE OWNED BY     ?   ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;
          public          postgres    false    223            �            1259    16398    products    TABLE     F  CREATE TABLE public.products (
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_valid boolean NOT NULL,
    id bigint NOT NULL,
    base_price bigint NOT NULL,
    discount bigint DEFAULT 0 NOT NULL,
    pay_price bigint NOT NULL,
    vat bigint,
    category text
);
    DROP TABLE public.products;
       public         heap    postgres    false            �            1259    16456    products_id_seq    SEQUENCE     x   CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.products_id_seq;
       public          postgres    false    215            1           0    0    products_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;
          public          postgres    false    219            �            1259    16421    tokens    TABLE       CREATE TABLE public.tokens (
    token text NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    validity_date timestamp without time zone NOT NULL,
    id bigint NOT NULL,
    is_admin bigint NOT NULL
);
    DROP TABLE public.tokens;
       public         heap    postgres    false            �            1259    16466    tokens_id_seq    SEQUENCE     v   CREATE SEQUENCE public.tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 $   DROP SEQUENCE public.tokens_id_seq;
       public          postgres    false    217            2           0    0    tokens_id_seq    SEQUENCE OWNED BY     ?   ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;
          public          postgres    false    220            l           2604    16434    customers id    DEFAULT     l   ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);
 ;   ALTER TABLE public.customers ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    218    216            s           2604    25034 
   loggers id    DEFAULT     h   ALTER TABLE ONLY public.loggers ALTER COLUMN id SET DEFAULT nextval('public.loggers_id_seq'::regclass);
 9   ALTER TABLE public.loggers ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    225    226    226            r           2604    16560    order_details id    DEFAULT     t   ALTER TABLE ONLY public.order_details ALTER COLUMN id SET DEFAULT nextval('public.order_details_id_seq'::regclass);
 ?   ALTER TABLE public.order_details ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    224    222            q           2604    16549 	   orders id    DEFAULT     f   ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);
 8   ALTER TABLE public.orders ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    223    221            j           2604    16457    products id    DEFAULT     j   ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);
 :   ALTER TABLE public.products ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    219    215            o           2604    16467 	   tokens id    DEFAULT     f   ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);
 8   ALTER TABLE public.tokens ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    220    217                      0    16414 	   customers 
   TABLE DATA           a   COPY public.customers (first_name, last_name, phone, address, email, id, created_at) FROM stdin;
    public          postgres    false    216   i<       &          0    25031    loggers 
   TABLE DATA           *   COPY public.loggers (log, id) FROM stdin;
    public          postgres    false    226   �=       "          0    16538    order_details 
   TABLE DATA           �   COPY public.order_details (order_id, product_id, product_base_price, product_discount, product_pay_price, product_vat, product_name, quantity, id) FROM stdin;
    public          postgres    false    222   �>       !          0    16510    orders 
   TABLE DATA           �   COPY public.orders (customer, inserted_date, updated_date, status, order_id, transaction_method, total_price, id, products) FROM stdin;
    public          postgres    false    221   K?                 0    16398    products 
   TABLE DATA           r   COPY public.products (name, created_at, is_valid, id, base_price, discount, pay_price, vat, category) FROM stdin;
    public          postgres    false    215   �?                 0    16421    tokens 
   TABLE DATA           Y   COPY public.tokens (token, user_id, created_at, validity_date, id, is_admin) FROM stdin;
    public          postgres    false    217   �@       3           0    0    customers_id_seq    SEQUENCE SET     ?   SELECT pg_catalog.setval('public.customers_id_seq', 22, true);
          public          postgres    false    218            4           0    0    loggers_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('public.loggers_id_seq', 2, true);
          public          postgres    false    225            5           0    0    order_details_id_seq    SEQUENCE SET     H   SELECT pg_catalog.setval('public.order_details_id_seq', 4627831, true);
          public          postgres    false    224            6           0    0    orders_id_seq    SEQUENCE SET     A   SELECT pg_catalog.setval('public.orders_id_seq', 1118220, true);
          public          postgres    false    223            7           0    0    products_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.products_id_seq', 11, true);
          public          postgres    false    219            8           0    0    tokens_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('public.tokens_id_seq', 1, false);
          public          postgres    false    220            w           2606    16436    customers customers_pkey 
   CONSTRAINT     V   ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);
 B   ALTER TABLE ONLY public.customers DROP CONSTRAINT customers_pkey;
       public            postgres    false    216            y           2606    16432    customers email_unique 
   CONSTRAINT     R   ALTER TABLE ONLY public.customers
    ADD CONSTRAINT email_unique UNIQUE (email);
 @   ALTER TABLE ONLY public.customers DROP CONSTRAINT email_unique;
       public            postgres    false    216                       2606    25081    orders idx_orders_order_id 
   CONSTRAINT     Y   ALTER TABLE ONLY public.orders
    ADD CONSTRAINT idx_orders_order_id UNIQUE (order_id);
 D   ALTER TABLE ONLY public.orders DROP CONSTRAINT idx_orders_order_id;
       public            postgres    false    221            �           2606    25038    loggers loggers_pkey 
   CONSTRAINT     R   ALTER TABLE ONLY public.loggers
    ADD CONSTRAINT loggers_pkey PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.loggers DROP CONSTRAINT loggers_pkey;
       public            postgres    false    226            �           2606    16537 
   orders oid 
   CONSTRAINT     I   ALTER TABLE ONLY public.orders
    ADD CONSTRAINT oid UNIQUE (order_id);
 4   ALTER TABLE ONLY public.orders DROP CONSTRAINT oid;
       public            postgres    false    221            �           2606    16562     order_details order_details_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT order_details_pkey PRIMARY KEY (id);
 J   ALTER TABLE ONLY public.order_details DROP CONSTRAINT order_details_pkey;
       public            postgres    false    222            �           2606    16551    orders orders_pkey 
   CONSTRAINT     P   ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);
 <   ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_pkey;
       public            postgres    false    221            {           2606    16491    customers phone_unique 
   CONSTRAINT     R   ALTER TABLE ONLY public.customers
    ADD CONSTRAINT phone_unique UNIQUE (phone);
 @   ALTER TABLE ONLY public.customers DROP CONSTRAINT phone_unique;
       public            postgres    false    216            u           2606    16459    products products_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.products DROP CONSTRAINT products_pkey;
       public            postgres    false    215            }           2606    16469    tokens tokens_pkey 
   CONSTRAINT     P   ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);
 <   ALTER TABLE ONLY public.tokens DROP CONSTRAINT tokens_pkey;
       public            postgres    false    217            �           1259    16683    orders_customerid    INDEX     H   CREATE INDEX orders_customerid ON public.orders USING btree (customer);
 %   DROP INDEX public.orders_customerid;
       public            postgres    false    221            �           1259    16682    orders_orderid_index    INDEX     R   CREATE UNIQUE INDEX orders_orderid_index ON public.orders USING btree (order_id);
 (   DROP INDEX public.orders_orderid_index;
       public            postgres    false    221            �           2606    16516    orders customer    FK CONSTRAINT     }   ALTER TABLE ONLY public.orders
    ADD CONSTRAINT customer FOREIGN KEY (customer) REFERENCES public.customers(id) NOT VALID;
 9   ALTER TABLE ONLY public.orders DROP CONSTRAINT customer;
       public          postgres    false    4727    221    216            �           2606    16543    order_details oid    FK CONSTRAINT     x   ALTER TABLE ONLY public.order_details
    ADD CONSTRAINT oid FOREIGN KEY (order_id) REFERENCES public.orders(order_id);
 ;   ALTER TABLE ONLY public.order_details DROP CONSTRAINT oid;
       public          postgres    false    221    222    4737               z  x����N�0�g�)����dJ�(CUĊ�b�FqR)M+�/�[Б�.LlI�;jQ��_����?6	xl&S;`!!�R����SU)8�4�����A����?o�c{
(�T�����@1e�C����#&�W�y3��t	�V��̭�X-�wm�5��a�_�#�?������U���%k�ݓ�4�<���C*�u�C����.u5o]$�)ÇǷe��m�V9�ynA �k��zn��M-�s��HDl�9��*3&��0d".U�1����I,�B��)?>�Y�G!&h$(��{#���݆ �o}*<J[��/��� dٛ����%E�b7@���8�E�wf�^�߭��+��!���<��|��-      &   �   x�M��
�0���S��4i�֛T�D<x��-�����Z���0;�#�f��{
�DW�Y����ſ����x*��	���p�iX)0�����m���Z�
��,[��e��9x��˙J�/������ڰ��O��ش��כּ�b���e�ʔ�υo��G�      "   �   x�m�1�0���9EN�l�N� � ,�4	ъ��Cڎ�����,w#䖔�����uFD`K���u}��9=�����YF��T�A<EAq(y��/��?�%�"��\=)�d	}�[����ʡP��B?fmB      !   e   x�M�;�0�z}
�43k������+]�/P����L��d6��8��@�ܶ�h�=��3��P2�b������(�h$G	ƧN�/�-U�Zغ�_�         �   x���=n�0��>��"�%z�
� ]:8h� 	�=~+C��Dć����z�����r�q:|]��eG�L��H��T�����4$���m:M��y�/�����$ٳX���40��>�MR��X.UF�Zҙ2�O�����B��ΕG2��[/k�ٽ�#a�N��>!mIRBG�L��$!��]Pz�uS7̖�(Wgs�f�"C�c㟗���h~�� �3զ��i���a���{         b   x���9
�0@�z���0K����L!b��D��*��a\�|�s�� $j��RC>0���;mI��?���瞶*N�$@�$���~%�\��-"^\d@      