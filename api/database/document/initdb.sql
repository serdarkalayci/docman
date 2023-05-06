USE docman;
CREATE TABLE public.spaces (
  id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  passive BOOL NOT NULL DEFAULT false,
  CONSTRAINT spaces_pkey PRIMARY KEY (id ASC),
  UNIQUE INDEX spaces_index_2 (name ASC)
);

CREATE TABLE public.documents (
  id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
  parent_id UUID NULL,
  name VARCHAR(255) NULL,
  content STRING NULL,
  space_id UUID NOT NULL,
  passive BOOL NOT NULL DEFAULT false,
  CONSTRAINT documents_pkey PRIMARY KEY (id ASC),
  CONSTRAINT documents_hierarchy FOREIGN KEY (parent_id) REFERENCES public.documents(id) ON DELETE SET NULL,
  CONSTRAINT documents_space FOREIGN KEY (space_id) REFERENCES public.spaces(id),
  UNIQUE INDEX documents_index_2 (space_id ASC, name ASC)
);

CREATE TABLE public.document_history (
  id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now():::TIMESTAMP,
  document_id UUID NOT NULL,
  content STRING NOT NULL,
  CONSTRAINT document_history_pkey PRIMARY KEY (id ASC),
  CONSTRAINT document_history_relation_1 FOREIGN KEY (document_id) REFERENCES public.documents(id)
);

INSERT INTO spaces (id, name) VALUES ('ec46101a-4779-423a-ac6b-b3e98d7f6990', 'space1');
INSERT INTO spaces (id, name) VALUES ('cd745f1c-6fdd-4ffc-9045-07f2ab2a4918', 'space2');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('34668788-1c63-4025-8787-5ff769458a51', NULL, 's1-doc1', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('d80ab59c-aced-4dee-aab1-22c55c3174d3', NULL, 's1-doc2', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('7acf232a-a4ad-41d4-9827-f0fd0e77a8fa', '34668788-1c63-4025-8787-5ff769458a51', 's1-doc1-1', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('a26498dd-1fed-44cf-839c-f0b807756709', '34668788-1c63-4025-8787-5ff769458a51', 's1-doc1-2', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('fb980aa4-f9d1-4a91-8097-e827175a0e1b', '7acf232a-a4ad-41d4-9827-f0fd0e77a8fa', 's1-doc1-1-1', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('4053c691-31dc-4ce6-af70-16e05274e9fa', 'd80ab59c-aced-4dee-aab1-22c55c3174d3', 's1-doc2-1', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('46e6e99c-3654-437d-8a2f-4a203ae09d0a', 'd80ab59c-aced-4dee-aab1-22c55c3174d3', 's1-doc2-2', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('1ffc1f08-8d5b-4419-a874-52c0458f423d', '46e6e99c-3654-437d-8a2f-4a203ae09d0a', 's1-doc2-2-2', 'ec46101a-4779-423a-ac6b-b3e98d7f6990', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('e3eb4455-ec0c-45c5-b826-a1a6c9e6dc47', NULL, 's2-doc1', 'cd745f1c-6fdd-4ffc-9045-07f2ab2a4918', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('5766836a-6e92-46ee-b296-8ee98fe09054', NULL, 's2-doc2', 'cd745f1c-6fdd-4ffc-9045-07f2ab2a4918', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');
INSERT INTO documents (id, parent_id, name, space_id, content) VALUES ('864b8c58-52e3-4b79-ad8e-43ecffc0cd37', 'e3eb4455-ec0c-45c5-b826-a1a6c9e6dc47', 's2-doc1-1', 'cd745f1c-6fdd-4ffc-9045-07f2ab2a4918', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.');

INSERT INTO document_history (id, document_id, content) VALUES ('cc2ffb2a-4df3-4f79-93b9-c2cd01010d83', '34668788-1c63-4025-8787-5ff769458a51', 'Old content of s1-doc1');
INSERT INTO document_history (id, document_id, content) VALUES ('4d052b18-c102-4943-b830-ce4824f7301d', 'd80ab59c-aced-4dee-aab1-22c55c3174d3', 'Old content of s1-doc2');