ALTER TABLE departments ADD COLUMN IF NOT EXISTS short_description_en VARCHAR(500);

UPDATE departments SET short_description_en = 'Corporate IT systems, infrastructure and internal software development.'
WHERE department_id = 1 AND (short_description_en IS NULL OR short_description_en = '');

UPDATE departments SET short_description_en = 'Accounting, tax compliance and financial reporting.'
WHERE department_id = 2 AND (short_description_en IS NULL OR short_description_en = '');

UPDATE departments SET short_description_en = 'Human resources, recruitment and staff development.'
WHERE department_id = 3 AND (short_description_en IS NULL OR short_description_en = '');

UPDATE departments SET short_description_en = 'Legal support for business operations and contracts.'
WHERE department_id = 4 AND (short_description_en IS NULL OR short_description_en = '');

UPDATE departments SET short_description_en = 'Procurement and supply of materials for the company.'
WHERE department_id = 5 AND (short_description_en IS NULL OR short_description_en = '');

UPDATE departments SET short_description_en = 'Office facilities management and internal logistics.'
WHERE department_id = 6 AND (short_description_en IS NULL OR short_description_en = '');
