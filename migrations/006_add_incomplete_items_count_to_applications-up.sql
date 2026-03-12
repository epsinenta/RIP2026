ALTER TABLE department_applications
    ADD COLUMN IF NOT EXISTS incomplete_items_count INTEGER NOT NULL DEFAULT 0;

UPDATE department_applications da
SET incomplete_items_count = (
    SELECT COUNT(*)::INTEGER
    FROM department_application_departments dad
    WHERE dad.department_application_id = da.department_application_id
      AND (COALESCE(TRIM(dad.role), '') = '' OR dad.salary IS NULL)
);
