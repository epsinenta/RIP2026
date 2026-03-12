ALTER TABLE department_applications
    ADD COLUMN IF NOT EXISTS main_department_id INTEGER REFERENCES departments(department_id);

UPDATE department_applications da
SET main_department_id = (
    SELECT dad.department_id
    FROM department_application_departments dad
    WHERE dad.department_application_id = da.department_application_id
      AND dad.is_main = TRUE
    LIMIT 1
)
WHERE main_department_id IS NULL;

ALTER TABLE department_application_departments
    DROP COLUMN IF EXISTS is_main;
