ALTER TABLE department_application_departments
    ADD COLUMN IF NOT EXISTS main_department_id INTEGER REFERENCES departments(department_id);

ALTER TABLE department_application_departments
    ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0;
